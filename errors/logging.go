package errors

import "fmt"

const ErrorLogTemplate = "[error_code=%s] [error_id=%s] %s"

func ToLogFormat(err ErrCodeErr) string {
	return fmt.Sprintf(ErrorLogTemplate, err.GetErrorCode().Code, err.GetId(), err.GetStackTrace())
}

func ToLogFormatWithoutStackTrace(err ErrCodeErr) string {
	return fmt.Sprintf(ErrorLogTemplate, err.GetErrorCode().Code, err.GetId(), err.Error())
}
