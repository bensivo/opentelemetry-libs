package main

import (
	"context"
	"time"
    "errors"

	"github.com/bensivo/opentelemetry-libs/packages/opentelemetry-tracing-go/pkg/tracing"
	"go.opentelemetry.io/otel/codes"
)

func main() {
    tracing.Initialize(tracing.InitializeOptions{
        ServiceName: "opentelemetry-tracing-go-example",
        ServiceVersion: "1.0.0",
        DeploymentEnvironment: "local",

        // OTLP Settings, in this case, we're pushing to NewRelic
        OtlpEndpoint: "https://otlp.nr-data.net:4318/v1/traces",
        OtlpHeaders: map[string]string{
            // Add your NewRelic License Key here
            // "api-key": "REPLACEME",
        },
    })
    defer tracing.Shutdown(context.Background())

    span1 := tracing.StartSpan(context.Background(), "span-1")
    span2 := tracing.StartSpan(context.Background(), "span-2")
    span3 := tracing.StartSpan(context.Background(), "span-3")

    time.Sleep(2 * time.Second)

    err := errors.New("This is an error")
    span3.SetStatus(codes.Error, err.Error())
    span3.RecordError(err)

    span1.End()
    span2.End()
    span3.End()

}
