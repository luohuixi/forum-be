package tracer

import (
	"context"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"

	"io"
	"time"
)

const TRACER_ENV = "TRACER"

type TracerValue struct {
	SamplerParam    float64
	FlusherInterval time.Duration
	LogSpans        bool
}

func NewTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
	tv := TraceValueFromEnv()

	cfg := jaegercfg.Configuration{
		ServiceName: serviceName, // tracer name
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: tv.SamplerParam,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            tv.LogSpans,
			BufferFlushInterval: tv.FlusherInterval,
		},
	}
	sender := transport.NewHTTPTransport(addr)
	reporter := jaeger.NewRemoteReporter(sender) // create Jaeger reporter
	// Initialize Opentracing tracer with Jaeger Reporter
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)
	return tracer, closer, err
}

func TraceValueFromEnv() *TracerValue {
	env := os.Getenv(TRACER_ENV)

	if env == "DEV" {
		return &TracerValue{
			SamplerParam:    1,
			FlusherInterval: 1 * time.Second,
			LogSpans:        true,
		}
	}

	return &TracerValue{
		SamplerParam:    0.05,
		FlusherInterval: 5 * time.Second,
		LogSpans:        false,
	}
}

func GetTraceId(ctx context.Context) string {
	traceID := "unknown"
	if span := opentracing.SpanFromContext(ctx); span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			traceID = sc.TraceID().String()
		}
	}

	return traceID
}
