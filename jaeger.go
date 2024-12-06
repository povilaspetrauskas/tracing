package tracing

import (
	"log"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// InitJaeger initializes and returns a Jaeger tracer.
func InitJaeger(serviceName string) (opentracing.Tracer, func()) {
	jaegerHost := os.Getenv("JAEGER_AGENT_HOST")
	if jaegerHost == "" {
		jaegerHost = "jaeger" // The name of the Jaeger container in Docker Compose
	}

	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1, // Sample all traces
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: jaegerHost + ":6831", // Jaeger agent address
		},
	}

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Jaeger: %v", err)
	}
	return tracer, func() {
		_ = closer.Close()
	}
}
