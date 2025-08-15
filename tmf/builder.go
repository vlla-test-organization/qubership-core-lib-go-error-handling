package tmf

import (
	"strconv"

	"github.com/vlla-test-organization/qubership-core-lib-go-error-handling/v3/errors"
)

type ResponseBuilder struct {
	id             string
	code           string
	reason         string
	message        string
	referenceError *string
	status         *string
	source         any
	meta           *map[string]any
	errors         *[]Error
	schemaType     string
	schemaLocation *string
}

type ErrorBuilder struct {
	id             string
	code           string
	reason         string
	message        *string
	referenceError *string
	status         *string
	source         any
	meta           *map[string]any
}

func NewResponseBuilder(err errors.ErrCodeErr) *ResponseBuilder {
	var causes *[]Error
	if multiErrT, ok := err.(errors.MultiCauseErr); ok {
		var multiErrors []Error
		for _, cause := range multiErrT.GetCauses() {
			var causeErr *Error
			if convertableErr, convertable := cause.(interface{ ToError() *Error }); convertable {
				causeErr = convertableErr.ToError()
			} else {
				causeErr = NewErrorBuilder(cause).Build()
			}
			multiErrors = append(multiErrors, *causeErr)
		}
		causes = &multiErrors
	} else {
		causes = nil
	}
	return &ResponseBuilder{
		id:         err.GetId(),
		code:       err.GetErrorCode().Code,
		reason:     err.GetErrorCode().Title,
		message:    err.GetDetail(),
		errors:     causes,
		schemaType: TypeV1_0,
	}
}

func (b *ResponseBuilder) Id(id string) *ResponseBuilder {
	b.id = id
	return b
}

func (b *ResponseBuilder) Code(code string) *ResponseBuilder {
	b.code = code
	return b
}

func (b *ResponseBuilder) Reason(reason string) *ResponseBuilder {
	b.reason = reason
	return b
}

func (b *ResponseBuilder) Message(message string) *ResponseBuilder {
	b.message = message
	return b
}

func (b *ResponseBuilder) ReferenceError(referenceError string) *ResponseBuilder {
	b.referenceError = &referenceError
	return b
}

func (b *ResponseBuilder) Status(status int) *ResponseBuilder {
	s := strconv.Itoa(status)
	b.status = &s
	return b
}

func (b *ResponseBuilder) Source(source any) *ResponseBuilder {
	b.source = source
	return b
}

func (b *ResponseBuilder) Meta(meta map[string]any) *ResponseBuilder {
	b.meta = &meta
	return b
}

func (b *ResponseBuilder) Errors(errors ...Error) *ResponseBuilder {
	b.errors = &errors
	return b
}

func (b *ResponseBuilder) Type(schemaType string) *ResponseBuilder {
	b.schemaType = schemaType
	return b
}

func (b *ResponseBuilder) SchemaLocation(schemaLocation string) *ResponseBuilder {
	b.schemaLocation = &schemaLocation
	return b
}

func (b *ResponseBuilder) Build() *Response {
	return &Response{
		Id:             b.id,
		Code:           b.code,
		Reason:         b.reason,
		Message:        b.message,
		ReferenceError: b.referenceError,
		Status:         b.status,
		Source:         b.source,
		Meta:           b.meta,
		Errors:         b.errors,
		Type:           b.schemaType,
		SchemaLocation: b.schemaLocation,
	}
}

func NewErrorBuilder(err errors.ErrCodeErr) *ErrorBuilder {
	detail := err.GetDetail()
	return &ErrorBuilder{
		id:      err.GetId(),
		code:    err.GetErrorCode().Code,
		reason:  err.GetErrorCode().Title,
		message: &detail,
	}
}

func (b *ErrorBuilder) Id(id string) *ErrorBuilder {
	b.id = id
	return b
}

func (b *ErrorBuilder) Code(code string) *ErrorBuilder {
	b.code = code
	return b
}

func (b *ErrorBuilder) Reason(reason string) *ErrorBuilder {
	b.reason = reason
	return b
}

func (b *ErrorBuilder) Message(message string) *ErrorBuilder {
	b.message = &message
	return b
}

func (b *ErrorBuilder) ReferenceError(referenceError string) *ErrorBuilder {
	b.referenceError = &referenceError
	return b
}

func (b *ErrorBuilder) Status(status int) *ErrorBuilder {
	s := strconv.Itoa(status)
	b.status = &s
	return b
}

func (b *ErrorBuilder) Source(source any) *ErrorBuilder {
	b.source = source
	return b
}

func (b *ErrorBuilder) Meta(meta map[string]any) *ErrorBuilder {
	b.meta = &meta
	return b
}

func (b *ErrorBuilder) Build() *Error {
	return &Error{
		Id:             b.id,
		Code:           b.code,
		Reason:         b.reason,
		Message:        b.message,
		ReferenceError: b.referenceError,
		Status:         b.status,
		Source:         b.source,
		Meta:           b.meta,
	}
}
