package tmf

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vlla-test-organization/qubership-core-lib-go-error-handling/v3/errors"
)

type CustomErr struct {
	*errors.ErrCodeError
	metaInt int
	source  string
}

func (e *CustomErr) ToError() *Error {
	return NewErrorBuilder(e).Meta(map[string]any{"int": e.metaInt}).Source(e.source).Build()
}

func TestNewResponseBuilder(t *testing.T) {
	assertions := require.New(t)
	err := errors.NewError(errors.ErrorCode{Code: "TEST-1001", Title: "Test 1001"}, "test 1001 detail", nil)

	meta := map[string]any{"test": "test"}
	response := NewResponseBuilder(err).
		Meta(meta).
		Source("/path").
		Status(404).
		Build()

	status := "404"
	expectedResponse := Response{
		Id:      err.Id,
		Code:    err.GetErrorCode().Code,
		Reason:  err.GetErrorCode().Title,
		Message: err.GetDetail(),
		Meta:    &meta,
		Status:  &status,
		Source:  "/path",
		Type:    TypeV1_0,
	}
	responseAsJson, pErr := json.Marshal(response)
	assertions.NoError(pErr)
	expectedResponseAsJson, pErr := json.Marshal(expectedResponse)
	assertions.NoError(pErr)
	expected := string(expectedResponseAsJson)
	actual := string(responseAsJson)
	assertions.Equal(expected, actual)
}

func TestNewResponseBuilderCustomErr(t *testing.T) {
	assertions := require.New(t)
	err := errors.New(CustomErr{}, errors.ErrorCode{Code: "TEST-1001", Title: "Test 1001"}, "test 1001 detail")

	meta := map[string]any{"test": "test"}
	response := NewResponseBuilder(err).
		Meta(meta).
		Source("/path").
		Status(404).
		Build()

	status := "404"
	expectedResponse := Response{
		Id:      err.Id,
		Code:    err.GetErrorCode().Code,
		Reason:  err.GetErrorCode().Title,
		Message: err.GetDetail(),
		Meta:    &meta,
		Status:  &status,
		Source:  "/path",
		Type:    TypeV1_0,
	}
	responseAsJson, pErr := json.Marshal(response)
	assertions.NoError(pErr)
	expectedResponseAsJson, pErr := json.Marshal(expectedResponse)
	assertions.NoError(pErr)
	expected := string(expectedResponseAsJson)
	actual := string(responseAsJson)
	assertions.Equal(expected, actual)
}

func TestNewResponseBuilderMultiCauseError(t *testing.T) {
	assertions := require.New(t)
	detail1 := "test 1002 1 detail"
	detail2 := "test 1002 1 detail"
	err1 := errors.NewError(errors.ErrorCode{Code: "TEST-1002", Title: "Test 1002"}, detail1, nil)
	err2 := errors.NewError(errors.ErrorCode{Code: "TEST-1002", Title: "Test 1002"}, detail2, nil)
	multiCauseError := errors.NewDefaultMultiCauseError([]errors.ErrCodeErr{err1, err2})

	meta := map[string]any{"test": "test"}
	response := NewResponseBuilder(multiCauseError).
		Meta(meta).
		Source("/path").
		Status(404).
		Build()

	status := "404"
	expectedResponse := Response{
		Id:      multiCauseError.Id,
		Code:    multiCauseError.GetErrorCode().Code,
		Reason:  multiCauseError.GetErrorCode().Title,
		Message: multiCauseError.GetDetail(),
		Meta:    &meta,
		Status:  &status,
		Source:  "/path",
		Errors: &[]Error{
			{
				Id:      err1.Id,
				Code:    err1.ErrorCode.Code,
				Reason:  err1.ErrorCode.Title,
				Message: &detail1,
			},
			{
				Id:      err2.Id,
				Code:    err2.ErrorCode.Code,
				Reason:  err2.ErrorCode.Title,
				Message: &detail2,
			}},
		Type: TypeV1_0,
	}
	responseAsJson, pErr := json.Marshal(response)
	assertions.NoError(pErr)
	expectedResponseAsJson, pErr := json.Marshal(expectedResponse)
	assertions.NoError(pErr)
	expected := string(expectedResponseAsJson)
	actual := string(responseAsJson)
	assertions.Equal(expected, actual)
}

func TestNewResponseBuilderMultiCauseErrorWithConvert(t *testing.T) {
	assertions := require.New(t)
	detail1 := "test 1002 1 detail"
	detail2 := "test 1002 2 detail"
	errCode1002 := errors.ErrorCode{Code: "TEST-1002", Title: "Test 1002"}
	err1 := errors.New(CustomErr{source: "/path/1", metaInt: 1}, errCode1002, detail1, nil)
	err2 := errors.New(CustomErr{source: "/path/2", metaInt: 2}, errCode1002, detail2, nil)
	multiCauseError := errors.NewDefaultMultiCauseError([]errors.ErrCodeErr{err1, err2})

	meta := map[string]any{"test": "test"}
	response := NewResponseBuilder(multiCauseError).
		Meta(meta).
		Source("/path").
		Status(404).
		Build()

	status := "404"
	expectedResponse := Response{
		Id:      multiCauseError.Id,
		Code:    multiCauseError.GetErrorCode().Code,
		Reason:  multiCauseError.GetErrorCode().Title,
		Message: multiCauseError.GetDetail(),
		Meta:    &meta,
		Status:  &status,
		Source:  "/path",
		Errors: &[]Error{
			{
				Id:      err1.Id,
				Code:    err1.ErrorCode.Code,
				Reason:  err1.ErrorCode.Title,
				Message: &detail1,
				Meta:    &map[string]any{"int": 1},
				Source:  "/path/1",
			},
			{
				Id:      err2.Id,
				Code:    err2.ErrorCode.Code,
				Reason:  err2.ErrorCode.Title,
				Message: &detail2,
				Meta:    &map[string]any{"int": 2},
				Source:  "/path/2",
			}},
		Type: TypeV1_0,
	}
	responseAsJson, pErr := json.Marshal(response)
	assertions.NoError(pErr)
	expectedResponseAsJson, pErr := json.Marshal(expectedResponse)
	assertions.NoError(pErr)
	expected := string(expectedResponseAsJson)
	actual := string(responseAsJson)
	assertions.Equal(expected, actual)
}
