package tracing

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var provider *sdktrace.TracerProvider

type InitializeOptions struct {
	ServiceName           string            // Name of the service. Use lowercase letters and hyphens. (e.g. "my-service")
	ServiceVersion        string            // Version of the service. Use semantic versioning. (e.g. "1.0.0")
	DeploymentEnvironment string            // Environment. Use short lowercase names where possible (e.g. "dev", "test", "uat", "prod")
	Exporter              string            // "otlp" or "console"
	OtlpEndpoint          string            // If exporter = 'otlp' GRPC endpoint and port. (e.g "localhost:4317")
	OtlpHeaders           map[string]string // if exporter = 'otlp' additional headers to attach to the request, usually for authentication
}

func Initialize(opts InitializeOptions) error {
	// Exporter defines where the traces should be sent to
	// In this case, we export traces using the OTLP Protocol, over gRPC

	var exporter sdktrace.SpanExporter
	var err error
	if opts.Exporter == "otlp" {
		exporter, err = otlptracegrpc.New(context.Background(),
			otlptracegrpc.WithEndpoint(opts.OtlpEndpoint),
			otlptracegrpc.WithHeaders(opts.OtlpHeaders),
		)
		if err != nil {
			fmt.Println("Failed to create otlp exporter:", err)
			return err
		}
	} else if opts.Exporter == "console" {
		exporter, err = stdouttrace.New(
		// stdouttrace.WithPrettyPrint(),
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

	// Provider is used to retrieve the tracer instance
	provider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // NOTE: Default batch timeout is 5 seconds
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(provider)
	tracer = provider.Tracer("github.com/bensivo/opentelemetry-libs/packages/opentelemetry-tracing-go")
	return nil
}

func Shutdown(ctx context.Context) {
	err := provider.ForceFlush(ctx)
	if err != nil {
		fmt.Println("Failed to flush:", err)
	}
	provider.Shutdown(ctx)
}

// StartSpan creates a new span with the given name
//
//	TODO: Add support for a parent span
func StartSpan(ctx context.Context, name string) trace.Span {
	ctx, span := tracer.Start(ctx, name)
	return span
}

type ClientHTTPSpanOptions struct {
	Method string // HTTP method (e.g. GET, POST, PUT, DELETE)
	URL    string // Full URL of the request (e.g. https://api.example.com/api/v1/users?name=ben)
	Route  string // (if known) HTTP route being hit, using placeholders where appropriate (e.g. /api/v1/users/{id})
}

// StartClientHTTPSpan creates a new span following the semantic conventions for HTTP spans
// See: https://opentelemetry.io/docs/specs/semconv/http/http-spans/
func StartClientHTTPSpan(ctx context.Context, opts ClientHTTPSpanOptions) trace.Span {
	name := fmt.Sprintf("%s %s", strings.ToUpper(opts.Method), opts.Route)

	ctx, span := tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindClient),
	)

	span.SetAttributes(attribute.KeyValue{
		Key:   "http.request.method",
		Value: attribute.StringValue(strings.ToUpper(opts.Method)),
	})
	span.SetAttributes(attribute.KeyValue{
		Key:   "url.full",
		Value: attribute.StringValue(opts.URL),
	})
	span.SetAttributes(attribute.KeyValue{
		Key:   "http.route",
		Value: attribute.StringValue(opts.Route),
	})

	// TODO: Parse the URL and set the following attributes:
	// - url.scheme
	// - url.path
	// - url.query
	return span
}

type ServerHTTPSpanOptions struct {
	Method  string // HTTP method (e.g. GET, POST, PUT, DELETE)
	Route   string // HTTP route being hit, using placeholders where appropriate (e.g. /api/v1/users/{id})
	URLPath string // The URL path, with placeholders filled in, and with query strings (e.g. /api/v1/users/123?name=ben)
}

// StartServerHTTPSpan creates a new span following the semantic conventions for HTTP spans
// See: https://opentelemetry.io/docs/specs/semconv/http/http-spans/
func StartServerHTTPSpan(ctx context.Context, opts ServerHTTPSpanOptions) trace.Span {
	name := fmt.Sprintf("%s %s", strings.ToUpper(opts.Method), opts.Route)

	ctx, span := tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindServer),
	)

	span.SetAttributes(attribute.KeyValue{
		Key:   "http.request.method",
		Value: attribute.StringValue(strings.ToUpper(opts.Method)),
	})
	span.SetAttributes(attribute.KeyValue{
		Key:   "http.route",
		Value: attribute.StringValue(opts.Route),
	})
	span.SetAttributes(attribute.KeyValue{
		Key:   "url.path",
		Value: attribute.StringValue(opts.URLPath),
	})
	return span
}
