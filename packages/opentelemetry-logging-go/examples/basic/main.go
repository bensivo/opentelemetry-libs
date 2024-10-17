package main

import (
	"time"

	"github.com/bensivo/opentelemetry-libs/packages/opentelemetry-logging-go/pkg/logging"
)

func main() {

	logging.Initialize(logging.InitializeOptions{
		ServiceName:           "opentelemetry-logging-go-example",
		ServiceVersion:        "1.0.0",
		DeploymentEnvironment: "local",
		Exporter:              "console",
		// Exporter:     "otlp",
		// OtlpEndpoint: "otlp.nr-data.net:4317",
		// OtlpHeaders: map[string]string{
		//     "api-key": "replaceme", // Add your NewRelic License Key here
		// },
	})

	logging.Info("This is an info log")

	time.Sleep(10 * time.Second)
}
