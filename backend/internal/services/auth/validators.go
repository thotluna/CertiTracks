package auth

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	lowerCaseRegex   = regexp.MustCompile(`[a-z]`)
	upperCaseRegex   = regexp.MustCompile(`[A-Z]`)
	digitRegex       = regexp.MustCompile(`[0-9]`)
	specialCharRegex = regexp.MustCompile(`[!@#$%^&*]`)
)

type AuthValidators struct{}

func NewAuthValidators() *AuthValidators {
	return &AuthValidators{}
}

func (v *AuthValidators) Register(validate *validator.Validate) error {
	return validate.RegisterValidation("strong_password", validateStrongPassword)
}

// validateStrongPassword checks if a password meets the following criteria:
// - At least 8 characters long
// - Contains at least one lowercase letter
// - Contains at least one uppercase letter
// - Contains at least one digit
// - Contains at least one special character (!@#$%^&*)
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	return lowerCaseRegex.MatchString(password) &&
		upperCaseRegex.MatchString(password) &&
		digitRegex.MatchString(password) &&
		specialCharRegex.MatchString(password)
}
