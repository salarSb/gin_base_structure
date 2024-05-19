package helper

import (
	"base_structure/src/pkg/service_errors"
	"net/http"
)

var StatusCodeMapping = map[string]int{
	//Otp
	service_errors.OtpExists:   http.StatusConflict,
	service_errors.OtpUsed:     http.StatusConflict,
	service_errors.OtpNotValid: http.StatusBadRequest,

	//User
	service_errors.EmailExists:      http.StatusConflict,
	service_errors.UsernameExists:   http.StatusConflict,
	service_errors.PermissionDenied: http.StatusForbidden,

	//DB
	service_errors.RecordNotFound: http.StatusNotFound,
}

func TranslateErrorToStatusCode(err error) int {
	value, ok := StatusCodeMapping[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return value
}
