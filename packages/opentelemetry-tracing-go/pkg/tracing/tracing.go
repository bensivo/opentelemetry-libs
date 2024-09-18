package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var provider *sdktrace.TracerProvider

type InitializeOptions struct {
    ServiceName string
    ServiceVersion string
    DeploymentEnvironment string

    OtlpEndpoint string // GRPC endpoint and port, do not include protocol. e.g "localhost:4317"
    OtlpHeaders map[string]string // Headers, usually for authentication
}

func Initialize(opts InitializeOptions) {
    // Exporter defines where the traces should be sent to
    // In this case, we export traces using the OTLP Protocol, over gRPC
    exporter, err := otlptracegrpc.New(context.Background(), 
        otlptracegrpc.WithEndpointURL(opts.OtlpEndpoint),
        otlptracegrpc.WithHeaders(opts.OtlpHeaders),
    )
    if err != nil {
        fmt.Println("Failed to create stdout exporter:", err)
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

    // Provider is used to retrieve the tracer instance
    provider = sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter), // NOTE: Default batch timeout is 5 seconds
        sdktrace.WithResource(resource),
    )

    otel.SetTracerProvider(provider)
    tracer = provider.Tracer("github.com/bensivo/opentelemetry-libs/packages/opentelemetry-tracing-go")
}

func Shutdown(ctx context.Context) {
    err := provider.ForceFlush(ctx)
    if err != nil {
        fmt.Println("Failed to flush:", err)
    }
    provider.Shutdown(ctx)
}

// StartSpan creates a new span with the given name
//   TODO: Add support for a parent span
func StartSpan(ctx context.Context, name string) trace.Span {
    ctx, span := tracer.Start(ctx, name)
    return span
}

