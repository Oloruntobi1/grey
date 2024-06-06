package handlers

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// GenericMeta ...
type GenericMeta struct {
	TraceID *string `json:"trace_id,omitempty"`
}

func NewGenericMeta(ctx context.Context) *GenericMeta {
	return &GenericMeta{
		TraceID: Extract(ctx),
	}
}

func String(v string) *string {
	return &v
}

// ErrorField ...
type ErrorField struct {
	Name    *string `json:"name,omitempty"`
	Message *string `json:"message,omitempty"`
}

// ErrorBase ...
type ErrorBase struct {
	Code    *string      `json:"code,omitempty"`
	Message *string      `json:"message,omitempty"`
	Fields  []ErrorField `json:"fields,omitempty"`
}

type MessageBase struct {
	Message *string `json:"message,omitempty"`
}

// Error ...
type Error struct {
	Error *ErrorBase   `json:"error,omitempty"`
	Meta  *GenericMeta `json:"meta,omitempty"`
}

// Data ...
type Data struct {
	ID interface{} `json:"id,omitempty"`
}

// Success ...
type Success struct {
	Data interface{}  `json:"data,omitempty"`
	Meta *GenericMeta `json:"meta,omitempty"`
}

// ResponseWithID ...
func ResponseWithID(ctx context.Context, data ...interface{}) *Success {
	resp := &Success{
		Meta: NewGenericMeta(ctx),
	}
	if len(data) > 0 {
		resp.Data = &Data{
			ID: data[0],
		}
	}
	return resp
}

// ResponseWithObj ...
func ResponseWithObj(ctx context.Context, obj interface{}) *Success {
	resp := &Success{
		Data: obj,
		Meta: NewGenericMeta(ctx),
	}
	return resp
}

func ResponseWithError(ctx context.Context, err *ErrorBase) *Error {
	return &Error{
		Error: err,
		Meta:  NewGenericMeta(ctx),
	}
}

func Extract(ctx context.Context) *string {
	sc := trace.SpanFromContext(ctx).SpanContext()
	var xRayTraceID string
	if sc.TraceID().IsValid() {
		xRayTraceID = XRayTraceID(sc.TraceID())
	}

	return String(xRayTraceID)
}

const FieldTraceID = "trace-id"

const (
	traceIDVersion         = "1"
	traceIDDelimiter       = "-"
	traceIDFirstPartLength = 8
)

func XRayTraceID(traceID trace.TraceID) string {
	otTraceID := traceID.String()
	return traceIDVersion + traceIDDelimiter + otTraceID[:traceIDFirstPartLength] +
		traceIDDelimiter + otTraceID[traceIDFirstPartLength:]
}
