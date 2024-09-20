package main

import (
	"context"
	"time"

	"github.com/bensivo/opentelemetry-libs/packages/opentelemetry-tracing-go/pkg/tracing"
)

func main() {
	err := tracing.Initialize(tracing.InitializeOptions{
		ServiceName:           "opentelemetry-tracing-go-example",
		ServiceVersion:        "1.0.0",
		DeploymentEnvironment: "local",
		Exporter: "console",
		// Exporter:     "otlp",
		// OtlpEndpoint: "otlp.nr-data.net:4317",
		// OtlpHeaders: map[string]string{
		// 	"api-key": "replaceme", // Add your NewRelic License Key here
		// },
	})

	if err != nil {
		panic(err)
	}
	defer tracing.Shutdown(context.Background())

	span := tracing.StartSpan(context.Background(), "example-span")
	time.Sleep(1 * time.Second)
	span.End()

	span = tracing.StartClientHTTPSpan(context.Background(), tracing.ClientHTTPSpanOptions{
		Method: "GET",
		URL:    "https://api.example.com/api/v1/users/1",
		Route:  "/api/v1/users/{id}",
	})
	time.Sleep(1 * time.Second)
	span.End()

	span = tracing.StartServerHTTPSpan(context.Background(), tracing.ServerHTTPSpanOptions{
		Method:  "GET",
		Route:   "/my-server-route",
		URLPath: "/my-server-route",
	})
	time.Sleep(1 * time.Second)
	span.End()

}
