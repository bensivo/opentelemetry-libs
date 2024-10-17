package logging

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var logger log.Logger
var provider sdklog.LoggerProvider

type InitializeOptions struct {
	ServiceName           string            // Name of the service. Use lowercase letters and hyphens. (e.g. "my-service")
	ServiceVersion        string            // Version of the service. Use semantic versioning. (e.g. "1.0.0")
	DeploymentEnvironment string            // Environment. Use short lowercase names where possible (e.g. "dev", "test", "uat", "prod")
	Exporter              string            // "otlp" or "console"
	OtlpEndpoint          string            // If exporter = 'otlp' GRPC endpoint and port. (e.g "localhost:4317")
	OtlpHeaders           map[string]string // if exporter = 'otlp' additional headers to attach to the request, usually for authentication
}

func Initialize(opts InitializeOptions) error {
	// Epxorter definese where the logs should be sent to
	var exporter sdklog.Exporter
	var err error
	if opts.Exporter == "otlp" {
		exporter, err = otlploggrpc.New(context.Background(),
			otlploggrpc.WithEndpoint(opts.OtlpEndpoint),
			otlploggrpc.WithHeaders(opts.OtlpHeaders),
		)
		if err != nil {
			fmt.Println("Failed to create otlp exporter:", err)
			return err
		}
	} else if opts.Exporter == "console" {
		exporter, err = stdoutlog.New(
		// stdoutlog.WithPrettyPrint(),
		)
		if err != nil {
			fmt.Println("Failed to create otlp exporter:", err)
			return err
		}
	} else {
		return errors.New("Invalid exporter. Use 'otlp' or 'console'")
	}

	// Resource defines key-value attributes attached to this tracer
	// these will be sent along with the traces themselves to the OTLP receiver
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(opts.ServiceName),
			semconv.ServiceVersion(opts.ServiceVersion),
			semconv.DeploymentEnvironment(opts.DeploymentEnvironment),
		),
	)
	if err != nil {
		fmt.Println("Failed to create resource:", err)
		return err
	}

	batchProcessor := sdklog.NewBatchProcessor(exporter,
		sdklog.WithExportInterval(5*time.Second),
	)

	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(batchProcessor),
		sdklog.WithResource(resource),
	)
	logger = provider.Logger("github.com/bensivo/opentelemetry-libs/packages/opentelemetry-logging-go/pkg/logging")
	return nil
}

func Info(msg string) {
	record := log.Record{
		// timestamp:    time.Now(),
		// severity:     log.SeverityInfo,
		// severityText: "INFO",
		// body:         msg,
	}
	record.SetTimestamp(time.Now())
	record.SetSeverity(log.SeverityInfo)
	record.SetSeverityText(log.SeverityInfo.String())
	record.SetBody(log.StringValue(msg))

	// TODO: support adding attributes

	logger.Emit(context.Background(), record)
}
