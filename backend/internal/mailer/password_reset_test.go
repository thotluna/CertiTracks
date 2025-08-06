package mailer_test

import (
	"bytes"
	"html/template"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordResetTemplate(t *testing.T) {
	// Configuración de prueba
	testData := struct {
		Data struct {
			ResetURL string
			Name     string
		}
	}{
		Data: struct {
			ResetURL string
			Name     string
		}{
			ResetURL: "http://example.com/reset?token=abc123",
			Name:     "Test User",
		},
	}

	// Obtener la ruta absoluta al directorio actual
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "templates/password_reset.html")

	// Cargar y renderizar el template
	tmpl, err := template.ParseFiles(templatePath)
	require.NoError(t, err, "No se pudo cargar el template")

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, testData)

	// Aserciones
	assert.NoError(t, err, "Error al renderizar el template")
	assert.Contains(t, buf.String(), testData.Data.ResetURL, "El template debe contener la URL de reset")
	assert.Contains(t, buf.String(), testData.Data.Name, "El template debe contener el nombre del usuario")
}
