package tracing

import (
	"context"
	"log"
	"os"
	"path"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer(getName())

type SNSMessageAttribute struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

func getName() string {
	name, err := os.Executable()
	if err != nil {
		return "fiber-server"
	}
	return path.Base(name)
}

func InitTracer(serviceName string, tempoEndpoint string) *sdktrace.TracerProvider {
	ctx := context.Background()
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(tempoEndpoint), otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	stdoutexporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithBatcher(stdoutexporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func extractMetricsAttributesFromSpan(span oteltrace.Span) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	readOnlySpan, ok := span.(trace.ReadOnlySpan)
	if !ok {
		return attrs
	}
	attrs = append(attrs, readOnlySpan.Attributes()...)
	return attrs
}

func NewSpan(ctx context.Context, spanName string) (context.Context, oteltrace.Span) {
	attrs := extractMetricsAttributesFromSpan(oteltrace.SpanFromContext(ctx))
	ctx, span := tracer.Start(ctx, spanName, oteltrace.WithAttributes(
		attrs...,
	))
	return ctx, span
}

func GenerateContext(oldCtx context.Context, spanName string) (context.Context, oteltrace.Span, error) {
	ctx, span := generateSpanAndContext(oldCtx, spanName, false)
	return ctx, span, nil
}

func generateSpanAndContext(oldCtx context.Context, spanName string, generateNewContext bool) (context.Context, oteltrace.Span) {
	ctx := oldCtx
	if generateNewContext {
		ctx = context.WithoutCancel(ctx)
	}
	ctx, span := NewSpan(ctx, spanName)
	return ctx, span
}

func StringAttribute(key, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

func IntAttribute(key string, value int) attribute.KeyValue {
	return attribute.Int(key, value)
}

func NewSpanWithTraceparent(ctx context.Context, spanName string, traceparent string) (context.Context, oteltrace.Span) {
	if traceparent != "" {
		carrier := propagation.MapCarrier{}
		carrier.Set("traceparent", traceparent)
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	}
	return NewSpan(ctx, spanName)
}
