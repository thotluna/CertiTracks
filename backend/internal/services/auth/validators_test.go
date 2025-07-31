package auth_test

import (
	"certitrack/internal/services/auth"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateStrongPassword(t *testing.T) {
	v, _ := binding.Validator.Engine().(*validator.Validate)
	_ = v.RegisterValidation("strong_password", auth.ValidateStrongPassword)

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "valid password",
			password: "SecurePass123!",
			want:     true,
		},
		{
			name:     "too short",
			password: "Short1!",
			want:     false,
		},
		{
			name:     "missing uppercase",
			password: "lowercase123!",
			want:     false,
		},
		{
			name:     "missing lowercase",
			password: "UPPERCASE123!",
			want:     false,
		},
		{
			name:     "missing number",
			password: "NoNumbers!",
			want:     false,
		},
		{
			name:     "missing special char",
			password: "NoSpecialChar123",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.password, "strong_password")
			if tt.want {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
