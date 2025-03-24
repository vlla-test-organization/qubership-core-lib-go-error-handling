package tmf

import (
	"testing"

	"github.com/google/uuid"
	"github.com/netcracker/qubership-core-lib-go-error-handling/v3/errors"
	"github.com/stretchr/testify/require"
)

var (
	responseBuilder = ResponseBuilder{}
	errorBuilder    = ErrorBuilder{}
	tmfConverter    = DefaultConverter{}
)

func TestBuildErrorCodeError(t *testing.T) {
	assertions := require.New(t)
	id := uuid.New().String()
	code := "TEST-0001"
	reason := "Test reason"
	message := "Test detail"
	meta := map[string]any{
		"test-key": "test-value",
	}
	status := 404
	source := "/path"

	response := responseBuilder.Id(id).Code(code).Reason(reason).Message(message).Meta(meta).Status(status).Source(source).Build()
	converter := DefaultConverter{}
	remoteErr := converter.BuildErrorCodeError(*response)
	assertions.NotNil(remoteErr)
	remoteErrCodeError, ok := remoteErr.(*errors.RemoteErrCodeError)
	assertions.True(ok)
	assertions.Equal(id, remoteErrCodeError.Id)
	assertions.Equal(code, remoteErrCodeError.GetErrorCode().Code)
	assertions.Equal(reason, remoteErrCodeError.GetErrorCode().Title)
	assertions.Equal(message, remoteErrCodeError.GetDetail())
	assertions.Equal(meta, remoteErrCodeError.Meta)
	assertions.Equal(status, *remoteErrCodeError.Status)
	assertions.Equal(source, remoteErrCodeError.Source)
}

func TestBuildErrorCodeErrorMultiCauseErr(t *testing.T) {
	assertions := require.New(t)
	id := uuid.New().String()
	code := "TEST-0001"
	reason := "Test reason"
	message := "Test detail"
	meta := map[string]any{
		"test-key": "test-value",
	}
	status := 404
	source := "/path"

	err1Id := uuid.New().String()
	err1Code := "TEST-0001"
	err1Reason := "Test reason"
	err1Message := "Test detail"
	err1Status := 404

	err2Id := uuid.New().String()
	err2Code := "TEST-0001"
	err2Reason := "Test reason"
	err2Message := "Test detail"
	err2Status := 405

	response := responseBuilder.Id(id).Code(code).Reason(reason).Message(message).Meta(meta).Status(status).Source(source).
		Errors(
			*errorBuilder.Id(err1Id).Code(err1Code).Reason(err1Reason).Message(err1Message).Status(err1Status).Build(),
			*errorBuilder.Id(err2Id).Code(err2Code).Reason(err2Reason).Message(err2Message).Status(err2Status).Build(),
		).
		Build()
	remoteErr := tmfConverter.BuildErrorCodeError(*response)
	assertions.NotNil(remoteErr)
	remoteMultiCauseErr, ok := remoteErr.(*errors.RemoteMultiCauseError)
	assertions.True(ok)
	assertions.Equal(id, remoteMultiCauseErr.Id)
	assertions.Equal(code, remoteMultiCauseErr.GetErrorCode().Code)
	assertions.Equal(reason, remoteMultiCauseErr.GetErrorCode().Title)
	assertions.Equal(message, remoteMultiCauseErr.GetDetail())
	assertions.Equal(meta, remoteMultiCauseErr.Meta)
	assertions.Equal(status, *remoteMultiCauseErr.Status)
	assertions.Equal(source, remoteMultiCauseErr.Source)

	assertions.Equal(2, len(remoteMultiCauseErr.Causes))
	causeErr1 := remoteMultiCauseErr.Causes[0]
	assertions.NotNil(causeErr1)
	assertions.Equal(err1Id, causeErr1.Id)
	assertions.Equal(err1Code, causeErr1.GetErrorCode().Code)
	assertions.Equal(err1Reason, causeErr1.GetErrorCode().Title)
	assertions.Equal(err1Message, causeErr1.GetDetail())
	assertions.Equal(err1Status, *causeErr1.Status)

	causeErr2 := remoteMultiCauseErr.Causes[1]
	assertions.NotNil(causeErr2)
	assertions.Equal(err2Id, causeErr2.Id)
	assertions.Equal(err2Code, causeErr2.GetErrorCode().Code)
	assertions.Equal(err2Reason, causeErr2.GetErrorCode().Title)
	assertions.Equal(err2Message, causeErr2.GetDetail())
	assertions.Equal(err2Status, *causeErr2.Status)
}
