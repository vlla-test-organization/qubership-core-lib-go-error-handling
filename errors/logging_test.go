package errors

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestErrCodeErrorWithLogTemplate(t *testing.T) {
	assertions := require.New(t)
	err := NewError(ErrorCode{Code: "TEST-1001", Title: "Test 1001"}, "test 1001 detail", nil)
	logStr := ToLogFormat(err)
	assertions.Equal(fmt.Sprintf("[error_code=%s] [error_id=%s] %s", err.GetErrorCode().Code, err.GetId(), err.GetStackTrace()), logStr)
}
