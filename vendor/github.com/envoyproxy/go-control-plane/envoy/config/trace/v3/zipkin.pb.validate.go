// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/config/trace/v3/zipkin.proto

package envoy_config_trace_v3

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
)

// Validate checks the field values on ZipkinConfig with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ZipkinConfig) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetCollectorCluster()) < 1 {
		return ZipkinConfigValidationError{
			field:  "CollectorCluster",
			reason: "value length must be at least 1 runes",
		}
	}

	if utf8.RuneCountInString(m.GetCollectorEndpoint()) < 1 {
		return ZipkinConfigValidationError{
			field:  "CollectorEndpoint",
			reason: "value length must be at least 1 runes",
		}
	}

	// no validation rules for TraceId_128Bit

	if v, ok := interface{}(m.GetSharedSpanContext()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ZipkinConfigValidationError{
				field:  "SharedSpanContext",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for CollectorEndpointVersion

	// no validation rules for CollectorHostname

	return nil
}

// ZipkinConfigValidationError is the validation error returned by
// ZipkinConfig.Validate if the designated constraints aren't met.
type ZipkinConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ZipkinConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ZipkinConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ZipkinConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ZipkinConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ZipkinConfigValidationError) ErrorName() string { return "ZipkinConfigValidationError" }

// Error satisfies the builtin error interface
func (e ZipkinConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sZipkinConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ZipkinConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ZipkinConfigValidationError{}
