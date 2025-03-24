package errors

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"runtime/debug"
	"strings"
)

type ErrorCode struct {
	Code  string
	Title string
}

type ErrCodeErr interface {
	error
	GetId() string
	GetErrorCode() ErrorCode
	GetDetail() string
	GetStackTrace() string
}

type MultiCauseErr interface {
	ErrCodeErr
	GetCauses() []ErrCodeErr
}

type ErrCodeError struct {
	Id         string
	Name       string
	ErrorCode  ErrorCode
	Detail     string
	StackTrace []byte
	Cause      error
}

type MultiCauseError struct {
	*ErrCodeError
	Causes []ErrCodeErr
}

func (e *ErrCodeError) Error() string {
	name := e.Name
	if name == "" {
		name = "ErrCodeError"
	}
	prefix := fmt.Sprintf("%s [%s][%s] ", name, e.ErrorCode.Code, e.Id)
	if e.Detail != "" {
		return prefix + e.Detail
	} else {
		return prefix + e.ErrorCode.Title
	}
}

func (e *ErrCodeError) Unwrap() error { return e.Cause }

func (e *ErrCodeError) GetId() string           { return e.Id }
func (e *ErrCodeError) GetErrorCode() ErrorCode { return e.ErrorCode }
func (e *ErrCodeError) GetDetail() string       { return e.Detail }

func (e *ErrCodeError) GetStackTrace() string {
	s := e.Error() + "\n" + string(e.StackTrace)
	prefix := " "
	switch i := e.Cause.(type) {
	case ErrCodeErr:
		s += "Caused by: "
		stackTrace := i.GetStackTrace()
		nlCount := strings.Count(stackTrace, "\n")
		s += strings.Replace(stackTrace, "\n", "\n"+prefix, nlCount-1)
	case error:
		s += "Caused by: "
		s += prefix + i.Error()
	}
	return s
}

func (e *MultiCauseError) Error() string {
	prefix := "MultiCauseError [" + e.ErrorCode.Code + "][" + e.Id + "] "
	if e.Detail != "" {
		return prefix + e.Detail
	} else {
		return prefix + e.ErrorCode.Title
	}
}

func (e *MultiCauseError) GetStackTrace() string {
	s := e.Error() + "\n"
	prefix := " "
	for i, cause := range e.Causes {
		s += fmt.Sprintf("Caused by (%d/%d): ", i+1, len(e.Causes))
		switch i := cause.(type) {
		case ErrCodeErr:
			stackTrace := i.GetStackTrace()
			nlCount := strings.Count(stackTrace, "\n")
			s += strings.Replace(stackTrace, "\n", "\n"+prefix, nlCount-1)
		default:
			s += prefix + i.Error()
		}
	}
	return s
}
func (e *MultiCauseError) GetCauses() []ErrCodeErr {
	return e.Causes
}

func NewError(code ErrorCode, detail string, cause error) *ErrCodeError {
	return &ErrCodeError{
		Id:         uuid.New().String(),
		ErrorCode:  code,
		Detail:     detail,
		StackTrace: debug.Stack(),
		Cause:      cause,
	}
}

// New
// Instantiates new ErrCodeErr of type U with auto setting embedded *ErrCodeError and returns reference to it
// template - struct of type U with pre-set fields
// cause - an optional parameter (allowed only 1 or no error)
func New[U ErrCodeErr](template U, code ErrorCode, detail string, cause ...error) *U {
	templateRef := &template
	s := reflect.ValueOf(templateRef).Elem()
	if s.Kind() != reflect.Struct {
		panic("template must be a Struct")
	}
	f := s.FieldByName("ErrCodeError")
	if !f.IsValid() && !f.CanSet() && f.Kind() != reflect.Pointer {
		panic("template struct must embed *ErrCodeError struct")
	}
	var cErr error
	if len(cause) == 1 {
		cErr = cause[0]
	}
	errCode := &ErrCodeError{
		Id:         uuid.New().String(),
		Name:       reflect.TypeOf(templateRef).Elem().Name(),
		ErrorCode:  code,
		Detail:     detail,
		StackTrace: debug.Stack(),
		Cause:      cErr,
	}
	f.Set(reflect.ValueOf(errCode))
	return templateRef
}

// NewDefaultMultiCauseError todo use varargs for causes
func NewDefaultMultiCauseError(causes []ErrCodeErr) *MultiCauseError {
	return NewMultiCauseError(ErrorCode{"NC-COMMON-2100", "multi-cause error"},
		"multiple independent errors have happened", causes)
}

// NewMultiCauseError todo use varargs for causes
func NewMultiCauseError(code ErrorCode, detail string, causes []ErrCodeErr) *MultiCauseError {
	return &MultiCauseError{
		ErrCodeError: &ErrCodeError{
			Id:         uuid.New().String(),
			ErrorCode:  code,
			Detail:     detail,
			StackTrace: debug.Stack(),
			Cause:      nil,
		},
		Causes: causes,
	}
}
