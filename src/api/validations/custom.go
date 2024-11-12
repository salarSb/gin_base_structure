package validations

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Property string `json:"property"`
	Tag      string `json:"tag"`
	Value    string `json:"value"`
	Message  string `json:"message"`
}

func GetValidationErrors(err error) *[]ValidationError {
	var validationErrors []ValidationError
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, err := range err.(validator.ValidationErrors) {
			var el ValidationError
			el.Property = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			el.Message = getCustomMessage(err)
			validationErrors = append(validationErrors, el)
		}
		return &validationErrors
	}
	return nil
}

func getCustomMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "len":
		return fmt.Sprintf("%s must be %s characters in length", field, param)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", field, param)
	case "ne":
		return fmt.Sprintf("%s must not be equal to %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "eqfield":
		return fmt.Sprintf("%s must be equal to the value of %s", field, param)
	case "nefield":
		return fmt.Sprintf("%s must not be equal to the value of %s", field, param)
	case "gtfield":
		return fmt.Sprintf("%s must be greater than the value of %s", field, param)
	case "gtefield":
		return fmt.Sprintf("%s must be greater than or equal to the value of %s", field, param)
	case "ltfield":
		return fmt.Sprintf("%s must be less than the value of %s", field, param)
	case "ltefield":
		return fmt.Sprintf("%s must be less than or equal to the value of %s", field, param)
	case "alpha":
		return fmt.Sprintf("%s can only contain alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s can only contain alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be a numeric value", field)
	case "number":
		return fmt.Sprintf("%s must be a valid number", field)
	case "hexadecimal":
		return fmt.Sprintf("%s must be a valid hexadecimal", field)
	case "hexcolor":
		return fmt.Sprintf("%s must be a valid HEX color code", field)
	case "rgb":
		return fmt.Sprintf("%s must be a valid RGB color code", field)
	case "rgba":
		return fmt.Sprintf("%s must be a valid RGBA color code", field)
	case "hsl":
		return fmt.Sprintf("%s must be a valid HSL color code", field)
	case "hsla":
		return fmt.Sprintf("%s must be a valid HSLA color code", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", field)
	case "base64":
		return fmt.Sprintf("%s must be a valid Base64 string", field)
	case "contains":
		return fmt.Sprintf("%s must contain '%s'", field, param)
	case "containsany":
		return fmt.Sprintf("%s must contain at least one of the following characters: '%s'", field, param)
	case "containsrune":
		return fmt.Sprintf("%s must contain the rune '%s'", field, param)
	case "excludes":
		return fmt.Sprintf("%s must not contain '%s'", field, param)
	case "excludesall":
		return fmt.Sprintf("%s must not contain any of the following characters: '%s'", field, param)
	case "excludesrune":
		return fmt.Sprintf("%s must not contain the rune '%s'", field, param)
	case "isbn":
		return fmt.Sprintf("%s must be a valid ISBN", field)
	case "isbn10":
		return fmt.Sprintf("%s must be a valid ISBN-10", field)
	case "isbn13":
		return fmt.Sprintf("%s must be a valid ISBN-13", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "uuid3":
		return fmt.Sprintf("%s must be a valid UUID v3", field)
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID v4", field)
	case "uuid5":
		return fmt.Sprintf("%s must be a valid UUID v5", field)
	case "ascii":
		return fmt.Sprintf("%s must contain only ASCII characters", field)
	case "printascii":
		return fmt.Sprintf("%s must contain only printable ASCII characters", field)
	case "multibyte":
		return fmt.Sprintf("%s must contain multibyte characters", field)
	case "datauri":
		return fmt.Sprintf("%s must be a valid Data URI", field)
	case "latitude":
		return fmt.Sprintf("%s must be a valid latitude coordinate", field)
	case "longitude":
		return fmt.Sprintf("%s must be a valid longitude coordinate", field)
	case "ssn":
		return fmt.Sprintf("%s must be a valid SSN", field)
	case "ip":
		return fmt.Sprintf("%s must be a valid IP address", field)
	case "ipv4":
		return fmt.Sprintf("%s must be a valid IPv4 address", field)
	case "ipv6":
		return fmt.Sprintf("%s must be a valid IPv6 address", field)
	case "cidr":
		return fmt.Sprintf("%s must be a valid CIDR notation IP address", field)
	case "cidrv4":
		return fmt.Sprintf("%s must be a valid CIDR notation IPv4 address", field)
	case "cidrv6":
		return fmt.Sprintf("%s must be a valid CIDR notation IPv6 address", field)
	case "tcp4_addr":
		return fmt.Sprintf("%s must be a valid TCPv4 address", field)
	case "tcp6_addr":
		return fmt.Sprintf("%s must be a valid TCPv6 address", field)
	case "tcp_addr":
		return fmt.Sprintf("%s must be a valid TCP address", field)
	case "udp4_addr":
		return fmt.Sprintf("%s must be a valid UDPv4 address", field)
	case "udp6_addr":
		return fmt.Sprintf("%s must be a valid UDPv6 address", field)
	case "udp_addr":
		return fmt.Sprintf("%s must be a valid UDP address", field)
	case "ip_addr":
		return fmt.Sprintf("%s must be a resolvable IP address", field)
	case "unix_addr":
		return fmt.Sprintf("%s must be a resolvable Unix address", field)
	case "mac":
		return fmt.Sprintf("%s must be a valid MAC address", field)
	case "hostname":
		return fmt.Sprintf("%s must be a valid hostname", field)
	case "fqdn":
		return fmt.Sprintf("%s must be a valid fully qualified domain name", field)
	case "unique":
		return fmt.Sprintf("%s must contain unique values", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, param)
	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime in the format '%s'", field, param)
	case "dir":
		return fmt.Sprintf("%s must be a valid directory", field)
	case "file":
		return fmt.Sprintf("%s must be a valid file", field)
	case "base64url":
		return fmt.Sprintf("%s must be a valid Base64 URL-encoded string", field)
	case "btc_addr":
		return fmt.Sprintf("%s must be a valid Bitcoin address", field)
	case "btc_addr_bech32":
		return fmt.Sprintf("%s must be a valid Bech32 Bitcoin address", field)
	case "eth_addr":
		return fmt.Sprintf("%s must be a valid Ethereum address", field)
	case "hostname_port":
		return fmt.Sprintf("%s must be a valid hostname with port", field)
	case "hostname_rfc1123":
		return fmt.Sprintf("%s must be a valid hostname according to RFC 1123", field)
	case "postcode_iso3166_alpha2":
		return fmt.Sprintf("%s must be a valid postcode for ISO 3166-1 alpha-2 country code", field)
	case "postcode_iso3166_alpha3":
		return fmt.Sprintf("%s must be a valid postcode for ISO 3166-1 alpha-3 country code", field)
	case "ir_mobile":
		return fmt.Sprintf("%s must be in IR mobile number format", field)
	case "password":
		return fmt.Sprintf("%s is not safe enough", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
