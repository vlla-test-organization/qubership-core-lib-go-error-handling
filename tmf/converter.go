package tmf

import (
	"strconv"

	"github.com/netcracker/qubership-core-lib-go-error-handling/v3/errors"
)

type Converter interface {
	BuildErrorCodeError(response Response) error
}

type DefaultConverter struct {
}

func (c *DefaultConverter) BuildErrorCodeError(response Response) error {
	errorCode := errors.ErrorCode{
		Code:  response.Code,
		Title: response.Reason,
	}
	var status *int
	if response.Status != nil {
		statusVal, _ := strconv.Atoi(*response.Status)
		status = &statusVal
	}
	var meta map[string]interface{}
	if response.Meta != nil {
		meta = *response.Meta
	} else {
		meta = make(map[string]interface{})
	}
	if response.Errors != nil && len(*response.Errors) > 0 {
		var causes []*errors.RemoteErrCodeError
		for _, tmfErr := range *response.Errors {
			causes = append(causes, buildErrorCodeError(tmfErr))
		}
		return errors.NewRemoteMultiCauseError(response.Id, errorCode, response.Message, meta, status,
			response.Source, causes)
	} else {
		return errors.NewRemoteErrCodeError(response.Id, errorCode, response.Message, meta, status,
			response.Source)
	}
}

func buildErrorCodeError(tmfErr Error) *errors.RemoteErrCodeError {
	errorCode := errors.ErrorCode{
		Code:  tmfErr.Code,
		Title: tmfErr.Reason,
	}
	var status *int
	if tmfErr.Status != nil {
		statusVal, _ := strconv.Atoi(*tmfErr.Status)
		status = &statusVal
	}
	var meta map[string]interface{}
	if tmfErr.Meta != nil {
		meta = *tmfErr.Meta
	} else {
		meta = make(map[string]interface{})
	}
	var message string
	if tmfErr.Message != nil {
		message = *tmfErr.Message
	} else {
		message = ""
	}
	return errors.NewRemoteErrCodeError(tmfErr.Id, errorCode, message, meta, status, tmfErr.Source)
}

// ErrToResponse
// Deprecated: Use NewResponseBuilder().Build() and NewErrorBuilder().Build() functions instead
func ErrToResponse(e errors.ErrCodeErr, status int) Response {
	var causes *[]Error
	if multiErrT, ok := e.(errors.MultiCauseErr); ok {
		var multiErrors []Error
		for _, cause := range multiErrT.GetCauses() {
			detail := cause.GetDetail()
			multiErrors = append(multiErrors, Error{
				Id:      cause.GetId(),
				Code:    cause.GetErrorCode().Code,
				Reason:  cause.GetErrorCode().Title,
				Message: &detail,
			})
		}
		causes = &multiErrors
	} else {
		causes = nil
	}
	statusAsStr := strconv.Itoa(status)
	return Response{
		Id:      e.GetId(),
		Code:    e.GetErrorCode().Code,
		Reason:  e.GetErrorCode().Title,
		Message: e.GetDetail(),
		Status:  &statusAsStr,
		Errors:  causes,
		Type:    TypeV1_0,
	}
}
