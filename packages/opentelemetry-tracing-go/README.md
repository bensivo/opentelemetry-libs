# opentelemetry-tracing-go

My opinionated OpenTelemetry tracing library for Go.

## Installation
```shell
go get github.com/bensivo/opentelemetry-libs/packages/opentelemetry-tracing-go/pkg/tracing
```


## Usage
```go
import (
    "context"
    "errors"

    "github.com/bensivo/opentelemetry-libs/packages/opentelemetry-tracing-go/pkg/tracing"
)

func main() {
    tracing.Initialize(tracing.InitializeOptions{
        ServiceName: "opentelemetry-tracing-go-example",
        ServiceVersion: "1.0.0",
        DeploymentEnvironment: "local",
        OtlpEndpoint: "https//otlp.nr-data.net:4318/v1/traces", // OTLP Endpoint, in this case, we're pushing to NewRelic
        OtlpHeaders: map[string]string{
            "api-key": "REPLACEME", // Add your NewRelic License Key here
        },
    })
    defer tracing.Shutdown(context.Background())

    span := tracing.StartSpan(context.Background(), "span-1")
    // Do something...
    span.End()
}
```
