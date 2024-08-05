package iceberg

import (
	"encoding/json"
	"strings"
)

const NamespaceDelimiter = "\x1F"

// PrimitiveTypes
// See: https://iceberg.apache.org/spec/#primitive-types
var (
	Bool          = ToType("bool")
	Int           = ToType("int")
	Long          = ToType("long")
	Float         = ToType("float")
	Double        = ToType("double")
	Date          = ToType("date")
	Time          = ToType("time")
	Timestamp     = ToType("timestamp")
	TimestampTZ   = ToType("timestamptz")
	TimestampNS   = ToType("timestamp_ns")
	TimestampTZNS = ToType("timestamptz_ns")
	String        = ToType("string")
	UUID          = ToType("uuid")
	Binary        = ToType("binary")
	// Decimal(P, S) via ToType("decimal[10, 2]")
	// Fixed L via ToType("fixed[16]")
)

func NamespaceString(ns Namespaces) string {
	return strings.Join(ns, NamespaceDelimiter)
}

func Ptr[T any](v T) *T {
	return &v
}

func ToType(v any) Type {
	data, _ := json.Marshal(v)
	return Type{union: data}
}
