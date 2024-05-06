package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

func main() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracer := otel.Tracer("otel-client")

	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "hello_test")
	defer span.End()
	fmt.Println(span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())

	propagator := otel.GetTextMapPropagator()
	// 将跟踪上下文信息注入到请求中
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())

	fmt.Println(resp)

}
