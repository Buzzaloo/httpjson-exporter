// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package httpjsonexporter // import "github.com/ExaForce/httpjsonexporter"

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type httpJSONLogsExporter struct {
	config         *Config
	client         *http.Client
	clientSettings *confighttp.ClientConfig
	logger         *zap.Logger
	settings       component.TelemetrySettings
}

func newLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	config *Config,
) (*httpJSONLogsExporter, error) {
	return &httpJSONLogsExporter{
		config:         config,
		clientSettings: &config.ClientConfig,
		client:         nil,
		logger:         set.Logger,
		settings:       set.TelemetrySettings,
	}, nil
}

func (e *httpJSONLogsExporter) start(ctx context.Context, host component.Host) error {
	client, err := e.clientSettings.ToClient(ctx, host, e.settings)
	if err != nil {
		return err
	}
	e.client = client
	return nil
}

func (e *httpJSONLogsExporter) shutdown(ctx context.Context) error {
	return nil
}

func (e *httpJSONLogsExporter) pushLogs(ctx context.Context, ld plog.Logs) error {
	// Convert logs to JSON
	jsonLogs, err := e.logsToJSON(ld)
	if err != nil {
		return fmt.Errorf("failed to convert logs to JSON: %w", err)
	}

	if len(jsonLogs) == 0 {
		e.logger.Debug("No logs to export")
		return nil
	}

	// Prepare the HTTP request body
	body, err := e.prepareRequestBody(jsonLogs)
	if err != nil {
		return fmt.Errorf("failed to prepare request body: %w", err)
	}

	// Send the HTTP request
	if err := e.sendRequest(ctx, body); err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}

	e.logger.Debug("Successfully exported logs", zap.Int("count", len(jsonLogs)))
	return nil
}

// logsToJSON converts plog.Logs to a slice of JSON byte slices (one per log record)
func (e *httpJSONLogsExporter) logsToJSON(ld plog.Logs) ([][]byte, error) {
	var jsonLogs [][]byte

	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resourceLogs := ld.ResourceLogs().At(i)
		
		for j := 0; j < resourceLogs.ScopeLogs().Len(); j++ {
			scopeLogs := resourceLogs.ScopeLogs().At(j)
			
			for k := 0; k < scopeLogs.LogRecords().Len(); k++ {
				logRecord := scopeLogs.LogRecords().At(k)
				
				// Convert log record to a map
				logMap := make(map[string]interface{})
				
				// Add timestamp
				if logRecord.Timestamp() != 0 {
					logMap["timestamp"] = logRecord.Timestamp().AsTime().Format("2006-01-02T15:04:05.000000000Z07:00")
				}
				
				// Add severity
				if logRecord.SeverityNumber() != plog.SeverityNumberUnspecified {
					logMap["severity_number"] = int(logRecord.SeverityNumber())
					logMap["severity_text"] = logRecord.SeverityText()
				}
				
				// Add body
				body := logRecord.Body()
				if body.Type() != 0 { // 0 is empty type
					logMap["body"] = body.AsString()
				}
				
				// Add attributes
				logRecord.Attributes().Range(func(k string, v pcommon.Value) bool {
					logMap[k] = v.AsRaw()
					return true
				})
				
				// Add resource attributes with a prefix to avoid collisions
				resourceLogs.Resource().Attributes().Range(func(k string, v pcommon.Value) bool {
					logMap["resource."+k] = v.AsRaw()
					return true
				})
				
				// Marshal to JSON
				jsonBytes, err := json.Marshal(logMap)
				if err != nil {
					e.logger.Error("Failed to marshal log record to JSON", zap.Error(err))
					continue
				}
				
				jsonLogs = append(jsonLogs, jsonBytes)
			}
		}
	}

	return jsonLogs, nil
}

// prepareRequestBody creates the request body with optional compression
func (e *httpJSONLogsExporter) prepareRequestBody(jsonLogs [][]byte) (io.Reader, error) {
	// Join all JSON logs with newlines (NDJSON format)
	var buf bytes.Buffer
	for i, jsonLog := range jsonLogs {
		buf.Write(jsonLog)
		if i < len(jsonLogs)-1 {
			buf.WriteByte('\n')
		}
	}

	// Apply compression if configured
	if e.config.Compression == "gzip" {
		var gzipBuf bytes.Buffer
		gzipWriter := gzip.NewWriter(&gzipBuf)
		
		if _, err := gzipWriter.Write(buf.Bytes()); err != nil {
			return nil, err
		}
		
		if err := gzipWriter.Close(); err != nil {
			return nil, err
		}
		
		return &gzipBuf, nil
	}

	return &buf, nil
}

// sendRequest sends the HTTP POST request
func (e *httpJSONLogsExporter) sendRequest(ctx context.Context, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.config.Endpoint, body)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-ndjson")
	req.Header.Set("Authorization", "Bearer "+string(e.config.BearerToken))
	
	if e.config.Compression == "gzip" {
		req.Header.Set("Content-Encoding", "gzip")
	}

	// Add any additional headers from config
	for key, value := range e.config.Headers {
		req.Header.Set(key, string(value))
	}

	// Send the request
	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
