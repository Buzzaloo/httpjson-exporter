// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package httpjsonexporter // import "github.com/ExaForce/httpjsonexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

var (
	// Type is the component type
	Type = component.MustNewType("httpjson")
)

const (
	// stabilityLevel is the stability level of the exporter
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a factory for HTTP JSON exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		Type,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, stability),
	)
}

func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	c := cfg.(*Config)

	// Create the HTTP JSON logs exporter
	logsExporter, err := newLogsExporter(ctx, set, c)
	if err != nil {
		return nil, err
	}

	// Wrap with standard exporter helper for retry, queue, and observability
	return exporterhelper.NewLogs(
		ctx,
		set,
		cfg,
		logsExporter.pushLogs,
		exporterhelper.WithStart(logsExporter.start),
		exporterhelper.WithShutdown(logsExporter.shutdown),
		exporterhelper.WithTimeout(exporterhelper.TimeoutConfig{Timeout: c.Timeout}),
		exporterhelper.WithRetry(c.BackOffConfig),
		exporterhelper.WithQueue(c.QueueSettings),
	)
}
