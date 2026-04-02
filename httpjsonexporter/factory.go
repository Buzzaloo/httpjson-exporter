// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package httpjsonexporter // import "github.com/ExaForce/httpjsonexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// typeStr is the type of the exporter
	typeStr = "httpjson"
	// stabilityLevel is the stability level of the exporter
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a factory for HTTP JSON exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, stability),
	)
}

func createLogsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Logs, error) {
	c := cfg.(*Config)

	// Create the HTTP JSON logs exporter
	logsExporter, err := newLogsExporter(ctx, set, c)
	if err != nil {
		return nil, err
	}

	// Wrap with standard exporter helper for retry, queue, and observability
	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,
		logsExporter.pushLogs,
		exporterhelper.WithStart(logsExporter.start),
		exporterhelper.WithShutdown(logsExporter.shutdown),
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: c.Timeout}),
		exporterhelper.WithRetry(c.BackOffConfig),
		exporterhelper.WithQueue(c.QueueSettings),
	)
}
