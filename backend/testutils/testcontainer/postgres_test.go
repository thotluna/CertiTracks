package testcontainer

import (
	"context"
	"testing"

	"certitrack/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSetupPostgres(t *testing.T) {
	// Configurar el contexto
	ctx := context.Background()

	t.Run("should initialize postgres container with test configuration", func(t *testing.T) {
		// Configuración de prueba para la base de datos
		testCfg := &config.Config{
			Database: config.DatabaseConfig{
				Name:     "test_db",
				User:     "testuser",
				Password: "testpassword",
				SSLMode:  "disable",
			},
		}

		// Configurar el contenedor de PostgreSQL
		pgContainer, err := SetupPostgres(ctx, testCfg)
		require.NoError(t, err, "No se pudo configurar el contenedor de PostgreSQL")
		defer pgContainer.Teardown(ctx)

		// Verificar que se asignó un puerto dinámico
		require.NotEmpty(t, pgContainer.Config.Database.Port, "Debe asignarse un puerto dinámico")
		require.NotEqual(t, "0", pgContainer.Config.Database.Port, "El puerto no puede ser 0 después de la asignación")

		// Verificar que la configuración se haya cargado correctamente
		require.NotEmpty(t, pgContainer.Config.Database.User, "El usuario de la base de datos no debería estar vacío")
		require.NotEmpty(t, pgContainer.Config.Database.Password, "La contraseña de la base de datos no debería estar vacía")
		require.NotEmpty(t, pgContainer.Config.Database.Name, "El nombre de la base de datos no debería estar vacío")
		require.Equal(t, "disable", pgContainer.Config.Database.SSLMode, "El modo SSL debería estar deshabilitado para pruebas")

		// Verificar que la base de datos se haya creado correctamente
		db := pgContainer.DB
		require.NotNil(t, db, "La conexión a la base de datos no debería ser nula")

		// Realizar una consulta simple para verificar la conexión
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err, "No se pudo ejecutar la consulta de prueba")
		require.Equal(t, 1, result, "El resultado de la consulta debería ser 1")

		// Verificar que las migraciones se aplicaron correctamente
		var tableExists bool
		err = db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists).Error
		require.NoError(t, err, "No se pudo verificar la existencia de la tabla 'users'")
		require.True(t, tableExists, "La tabla 'users' debería existir después de las migraciones")
	})

	t.Run("should use custom port when specified", func(t *testing.T) {
		// Usar un puerto que probablemente esté disponible
		customPort := "5434"
		t.Setenv("POSTGRES_TEST_PORT", customPort)

		// Configuración de prueba para la base de datos
		testCfg := &config.Config{
			Database: config.DatabaseConfig{
				Name:     "test_db_port",
				User:     "testuser",
				Password: "testpassword",
				SSLMode:  "disable",
			},
		}
		
		pgContainer, err := SetupPostgres(ctx, testCfg)
		require.NoError(t, err, "No se pudo configurar el contenedor de PostgreSQL")
		defer pgContainer.Teardown(ctx)

		// Verificar que el puerto sea el especificado
		require.Equal(t, customPort, pgContainer.Config.Database.Port, "Debería usar el puerto especificado en POSTGRES_TEST_PORT")
		
		// Verificar que la base de datos es accesible en el puerto personalizado
		db := pgContainer.DB
		require.NotNil(t, db, "La conexión a la base de datos no debería ser nula")

		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err, "No se pudo ejecutar la consulta de prueba en el puerto personalizado")
		require.Equal(t, 1, result, "El resultado de la consulta debería ser 1")
	})
}
