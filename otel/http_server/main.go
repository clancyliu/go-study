package main

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

func main() {

	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracer := otel.Tracer("otel-server")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		propagator := otel.GetTextMapPropagator()
		propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))
		ctx, span := tracer.Start(ctx, "http_request")
		defer span.End()

		fmt.Println(span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())

		w.Write([]byte("Hello World"))
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}

}
