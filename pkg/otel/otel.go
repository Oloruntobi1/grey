package otel

import (
	"context"
	"log/slog"
	"runtime"
	"strconv"

	"github.com/Oloruntobi1/grey/appconstants"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// OtelHandler implements slog.Handler
// It adds;
// (a) TraceIds & spanIds to logs.
// (b) Logs(as events) to the active span.
type OtelHandler struct{ H slog.Handler }

func (s OtelHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true /* support all logging levels*/
}

func (s OtelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return OtelHandler{H: s.H.WithAttrs(attrs)}
}

func (s OtelHandler) WithGroup(name string) slog.Handler {
	return OtelHandler{H: s.H.WithGroup(name)}
}

func (s OtelHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	if ctx == nil {
		return s.H.Handle(ctx, r)
	}

	span := oteltrace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return s.H.Handle(ctx, r)
	}

	{ // (a) adds TraceIds & spanIds to logs.
		//
		// TODO: (komuw) add stackTraces maybe.
		//
		sCtx := span.SpanContext()
		attrs := make([]slog.Attr, 0)
		if sCtx.HasTraceID() {
			attrs = append(attrs,
				slog.Attr{Key: "traceId", Value: slog.StringValue(sCtx.TraceID().String())},
			)
		}
		if sCtx.HasSpanID() {
			attrs = append(attrs,
				slog.Attr{Key: "spanId", Value: slog.StringValue(sCtx.SpanID().String())},
			)
		}
		if len(attrs) > 0 {
			r.AddAttrs(attrs...)
		}
	}

	{ // (b) adds logs to the active span as events.

		// code from: https://github.com/uptrace/opentelemetry-go-extra/tree/main/otellogrus
		// which is BSD 2-Clause license.

		attrs := make([]attribute.KeyValue, 0)

		logSeverityKey := attribute.Key("log.severity")
		logMessageKey := attribute.Key("log.message")
		attrs = append(attrs, logSeverityKey.String(r.Level.String()))
		attrs = append(attrs, logMessageKey.String(r.Message))

		callerKey := attribute.Key("caller")

		frames := runtime.CallersFrames([]uintptr{r.PC})

		for {
			frame, more := frames.Next()

			attrs = append(attrs, callerKey.StringSlice([]string{
				frame.Function,
				frame.File,
				strconv.Itoa(frame.Line),
			}))
			if !more {
				break
			}

		}

		// Ensuring non-zero time and non-empty keys for attributes
		if !r.Time.IsZero() {
			attrs = append(attrs, attribute.Key("time").String(r.Time.Format(appconstants.YearMonthDayHourMinuteSecond)))
		}

		// Only add non-empty keys for attributes
		r.Attrs(func(a slog.Attr) bool {
			if a.Key != "" {
				attrs = append(attrs,
					attribute.KeyValue{
						Key:   attribute.Key(a.Key),
						Value: attribute.StringValue(a.Value.String()),
					},
				)
			}
			return true
		})

		span.AddEvent("log", oteltrace.WithAttributes(attrs...))
		if r.Level >= slog.LevelError {
			span.SetStatus(codes.Error, r.Message)
		}
	}

	return s.H.Handle(ctx, r)
}
