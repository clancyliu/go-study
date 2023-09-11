package jaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	jaegerlog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"log"
	"time"
)

func NewJaegerTracer(serviceName string) (opentracing.Tracer, io.Closer) {
	sender, _ := jaeger.NewUDPTransport("", 0)
	tracer, closer := jaeger.NewTracer(serviceName,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender))

	return tracer, closer
}

func Hello() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831", // 替换host
		},
	}
	closer, err := cfg.InitGlobalTracer(
		"serviceName",
		jaegercfg.Logger(jaeger.StdLogger),
	)
	if err != nil {
		log.Println("Could not initialize jaeger tracer: %s", err.Error())
		return
	}

	var ctx = context.TODO()
	span1, ctx := opentracing.StartSpanFromContext(ctx, "span_1")

	span1.LogFields(jaegerlog.String("event", "string-format"))

	time.Sleep(time.Second / 2)
	span11, _ := opentracing.StartSpanFromContext(ctx, "span_1-1")
	time.Sleep(time.Second / 2)
	span11.Finish()
	span1.Finish()
	defer closer.Close()
}

var Tracer opentracing.Tracer

func NewTracer(servicename string, addr string, udp bool) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, //固定采样
			Param: 1,                       //1全采样，0不采样
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  addr, //"127.0.0.1:6831",
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}
	if udp {
		reporter := jaeger.NewRemoteReporter(sender)
		// Initialize tracer with a logger and a metrics factory
		return cfg.NewTracer(
			jaegercfg.Reporter(reporter),
		)
	}
	return cfg.NewTracer()
}
