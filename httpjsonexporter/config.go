// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package httpjsonexporter // import "github.com/ExaForce/httpjsonexporter"

import (
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for the HTTP JSON exporter.
type Config struct {
	// HTTPClientSettings contains HTTP client configuration.
	confighttp.ClientConfig `mapstructure:",squash"`

	// BearerToken is the authentication token for the HTTP endpoint.
	// This will be sent as "Authorization: Bearer <token>" header.
	BearerToken configopaque.String `mapstructure:"bearer_token"`

	// BatchSize is the number of log records to send in a single HTTP request.
	// Default is 100.
	BatchSize int `mapstructure:"batch_size"`

	// Compression specifies the compression algorithm to use.
	// Supported values: "none", "gzip"
	// Default is "none".
	Compression string `mapstructure:"compression"`

	// QueueSettings defines configuration for the exporter queue.
	QueueSettings exporterhelper.QueueConfig `mapstructure:"sending_queue"`

	// BackOffConfig defines configuration for retrying failed requests.
	BackOffConfig configretry.BackOffConfig `mapstructure:"retry_on_failure"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid.
func (cfg *Config) Validate() error {
	if cfg.Endpoint == "" {
		return errors.New("endpoint must be specified")
	}

	if cfg.BearerToken == "" {
		return errors.New("bearer_token must be specified")
	}

	if cfg.BatchSize < 1 {
		return errors.New("batch_size must be at least 1")
	}

	if cfg.Compression != "" && cfg.Compression != "none" && cfg.Compression != "gzip" {
		return errors.New("compression must be 'none' or 'gzip'")
	}

	return nil
}

// createDefaultConfig creates the default configuration for the exporter.
func createDefaultConfig() component.Config {
	return &Config{
		ClientConfig: confighttp.ClientConfig{
			Timeout: 30 * time.Second,
		},
		BatchSize:     100,
		Compression:   "none",
		QueueSettings: exporterhelper.NewDefaultQueueConfig(),
		BackOffConfig: configretry.NewDefaultBackOffConfig(),
	}
}
