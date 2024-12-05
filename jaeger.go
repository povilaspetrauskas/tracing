package tracing

import (
	"io"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// InitJaeger initializes and returns a Jaeger tracer.
func InitJaeger(serviceName string) (opentracing.Tracer, io.Closer) {
	jaegerHost := os.Getenv("JAEGER_AGENT_HOST")
	if jaegerHost == "" {
		jaegerHost = "localhost" // Default to localhost if not set
	}

	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1, // Sample all traces
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: jaegerHost + ":6831", // Jaeger agent address
			LogSpans:           true,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Fatalf("Failed to initialize Jaeger: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
