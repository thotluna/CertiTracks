package validators

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestEmailValidation(t *testing.T) {
	// Obtener el validador de Gin
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("No se pudo obtener el validador de Gin")
	}

	tests := []struct {
		name     string
		email    string
		hasError bool
	}{
		{
			name:     "email válido",
			email:    "usuario@ejemplo.com",
			hasError: false,
		},
		{
			name:     "email con subdominio",
			email:    "usuario@sub.dominio.com",
			hasError: false,
		},
		{
			name:     "email con caracteres especiales",
			email:    "usuario+tag@ejemplo.com",
			hasError: false,
		},
		{
			name:     "email inválido - sin arroba",
			email:    "usuario.ejemplo.com",
			hasError: true,
		},
		{
			name:     "email inválido - sin dominio",
			email:    "usuario@",
			hasError: true,
		},
		{
			name:     "email inválido - solo espacios",
			email:    "   ",
			hasError: true,
		},
		{
			name:     "email vacío",
			email:    "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.email, "required,email")
			if tt.hasError {
				assert.Error(t, err, "Se esperaba un error para el email: %s", tt.email)
			} else {
				assert.NoError(t, err, "No se esperaba un error para el email: %s", tt.email)
			}
		})
	}
}
