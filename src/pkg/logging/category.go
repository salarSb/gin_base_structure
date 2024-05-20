package logging

type Category string
type SubCategory string
type ExtraKey string

const (
	General         Category = "General"
	IO              Category = "Io"
	Internal        Category = "Internal"
	Postgres        Category = "Postgres"
	Redis           Category = "Redis"
	Validation      Category = "Validation"
	RequestResponse Category = "RequestResponse"
	Prometheus      Category = "Prometheus"
)

const (
	// StartUp => General, Internal, Redis
	StartUp SubCategory = "StartUp"
	// ExternalService => General
	ExternalService SubCategory = "ExternalService"

	// Migration => Postgres
	Migration SubCategory = "Migration"
	// Select => Postgres
	Select SubCategory = "Select"
	// Rollback => Postgres
	Rollback SubCategory = "Rollback"
	// Update => Postgres
	Update SubCategory = "Update"
	// Delete => Postgres
	Delete SubCategory = "Delete"
	// Insert => Postgres
	Insert SubCategory = "Insert"
	// Closing => Postgres
	Closing SubCategory = "Closing"

	// Api => Internal
	Api SubCategory = "Api"
	// HashPassword => Internal
	HashPassword SubCategory = "HashPassword"
	// DefaultRoleNotFound => Internal
	DefaultRoleNotFound SubCategory = "DefaultRoleNotFound"

	// MobileValidation => Validation
	MobileValidation SubCategory = "MobileValidation"
	// PasswordValidation => Validation
	PasswordValidation SubCategory = "PasswordValidation"

	// RemoveFile => IO
	RemoveFile SubCategory = "RemoveFile"
)

const (
	AppName      ExtraKey = "AppName"
	LoggerName   ExtraKey = "Logger"
	ClientIp     ExtraKey = "ClientIp"
	HostIp       ExtraKey = "HostIp"
	Method       ExtraKey = "Method"
	StatusCode   ExtraKey = "StatusCode"
	BodySize     ExtraKey = "BodySize"
	Path         ExtraKey = "Path"
	Latency      ExtraKey = "Latency"
	RequestBody  ExtraKey = "RequestBody"
	ResponseBody ExtraKey = "ResponseBody"
	ErrorMessage ExtraKey = "ErrorMessage"
)
