package constants

import (
	"base_structure/src/config"
	"base_structure/src/pkg/logging"
	"github.com/joho/godotenv"
	"os"
)

var logger = logging.NewLogger(config.GetConfig())
var (
	AdminRoleName          string
	DefaultRoleName        string
	DefaultUserName        string
	AdminFirstName         string
	AdminLastName          string
	AdminMobileNumber      string
	AdminPassword          string
	AdminEmail             string
	RedisOtpDefaultKey     string
	AuthorizationHeaderKey string
	UserIdKey              string
	FirstNameKey           string
	LastNameKey            string
	UsernameKey            string
	EmailKey               string
	MobileNumberKey        string
	RolesKey               string
	ExpireTimeKey          string
)

func InitConstants() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(logging.Internal, logging.StartUp, "error in reading from .env", nil)
		return
	}
	AdminRoleName = os.Getenv("ADMIN_ROLE_NAME")
	DefaultRoleName = os.Getenv("DEFAULT_ROLE_NAME")
	DefaultUserName = os.Getenv("DEFAULT_USER_NAME")
	AdminFirstName = os.Getenv("ADMIN_FIRST_NAME")
	AdminLastName = os.Getenv("ADMIN_LAST_NAME")
	AdminMobileNumber = os.Getenv("ADMIN_MOBILE_NUMBER")
	AdminPassword = os.Getenv("ADMIN_PASSWORD")
	AdminEmail = os.Getenv("ADMIN_EMAIL")
	RedisOtpDefaultKey = os.Getenv("REDIS_OTP_DEFAULT_KEY")
	AuthorizationHeaderKey = os.Getenv("AUTHORIZATION_HEADER_KEY")
	UserIdKey = os.Getenv("USER_ID_KEY")
	FirstNameKey = os.Getenv("FIRST_NAME_KEY")
	LastNameKey = os.Getenv("LAST_NAME_KEY")
	UsernameKey = os.Getenv("USERNAME_KEY")
	EmailKey = os.Getenv("EMAIL_KEY")
	MobileNumberKey = os.Getenv("MOBILE_NUMBER_KEY")
	RolesKey = os.Getenv("ROLES_KEY")
	ExpireTimeKey = os.Getenv("EXPIRE_TIME_KEY")
}
