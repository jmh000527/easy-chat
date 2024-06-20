package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"testing"
)

func Test_Jaeger(t *testing.T) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: fmt.Sprintf("http://%s/api/traces", "192.168.199.138:14268"),
		},
	}

	jaegerTracer, err := cfg.InitGlobalTracer("client test", jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := jaegerTracer.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// 任务的执行
	tracer := opentracing.GlobalTracer()
	// 任务节点定义Span
	parentSpan := tracer.StartSpan("A")
	defer parentSpan.Finish()

	B(tracer, parentSpan)
}

func B(tracer opentracing.Tracer, parentSpan opentracing.Span) {
	childSpan := tracer.StartSpan("B", opentracing.ChildOf(parentSpan.Context()))
	defer childSpan.Finish()
}
