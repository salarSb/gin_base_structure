package service_errors

type ServiceError struct {
	EndUserMessage   string `json:"endUserMessage"`
	TechnicalMessage string `json:"technicalMessage"`
	Err              error
}

func (e *ServiceError) Error() string {
	return e.EndUserMessage
}
