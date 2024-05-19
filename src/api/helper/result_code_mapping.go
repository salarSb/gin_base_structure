package helper

import (
	"base_structure/src/pkg/service_errors"
)

type ResultCode int

const (
	Success         ResultCode = 0
	BadRequestError ResultCode = 40001
	ValidationError ResultCode = 42201
	AuthError       ResultCode = 40101
	ForbiddenError  ResultCode = 40301
	NotFoundError   ResultCode = 40401
	ConflictError   ResultCode = 40901
	LimiterError    ResultCode = 42901
	OtpLimiterError ResultCode = 42902
	CustomRecovery  ResultCode = 50001
	InternalError   ResultCode = 50002
)

var ResultCodeMapping = map[string]ResultCode{
	service_errors.RecordNotFound: NotFoundError,
}

func TranslateErrorToResultCode(err error) ResultCode {
	value, ok := ResultCodeMapping[err.Error()]
	if !ok {
		return InternalError
	}
	return value
}
