package main

import (
	"context"
	"time"

	"github.com/bensivo/opentelemetry-libs/packages/opentelemetry-tracing-go/pkg/tracing"
)

func main() {
    tracing.Initialize(tracing.InitializeOptions{
        ServiceName: "opentelemetry-tracing-go-example",
        ServiceVersion: "1.0.0",
        DeploymentEnvironment: "local",
    })
    defer tracing.Shutdown(context.Background())

    span1 := tracing.StartSpan(context.Background(), "span-1")
    span2 := tracing.StartSpan(context.Background(), "span-2")
    span3 := tracing.StartSpan(context.Background(), "span-3")

    time.Sleep(2 * time.Second)

    span1.End()
    span2.End()
    span3.End()
}
