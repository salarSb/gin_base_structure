package service_errors

const (
	// OtpExists => Otp
	OtpExists = "otp exists"
	// OtpUsed => Otp
	OtpUsed = "otp used"
	// OtpNotValid => Otp
	OtpNotValid = "otp not valid"

	// UnexpectedError => Token
	UnexpectedError = "unexpected error"
	// ClaimsNotFound => Token
	ClaimsNotFound = "claims not found"
	// TokenRequired => Token
	TokenRequired = "token required"
	// TokenExpired => Token
	TokenExpired = "token expired"
	// TokenInvalid => Token
	TokenInvalid = "token invalid"

	// EmailExists => User
	EmailExists = "email exists"
	// UsernameExists => User
	UsernameExists = "username exists"
	// PermissionDenied => User
	PermissionDenied = "permission denied"

	// RecordNotFound => DB
	RecordNotFound = "record not found"
)
