package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"certitrack/internal/config"
	"certitrack/internal/database"
	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthIntegrationTestSuite struct {
	suite.Suite
	db          *gorm.DB
	router      *gin.Engine
	authService services.AuthService
}

// Test Suite para tokens inválidos
type InvalidTokenTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *AuthIntegrationTestSuite) SetupSuite() {
	// Cargar configuración de prueba
	os.Setenv("APP_ENV", "test")
	cfg, err := config.Load()
	if err != nil {
		suite.T().Fatalf("Error loading config: %v", err)
	}

	// Configurar base de datos de prueba
	db, err := database.Connect(cfg)
	if err != nil {
		suite.T().Fatalf("Error connecting to database: %v", err)
	}
	suite.db = db

	// Migrar esquemas
	err = database.AutoMigrate(db)
	if err != nil {
		suite.T().Fatalf("Error migrating database: %v", err)
	}

	// Inicializar repositorios
	userRepo := repositories.NewUserRepositoryImpl(db)

	// Inicializar servicios
	suite.authService = services.NewAuthService(cfg, userRepo)

	// Configurar router
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	// Configurar rutas
	authHandler := handlers.NewAuthHandler(suite.authService)
	authGroup := suite.router.Group("/api/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}

	// Grupo protegido para probar autenticación
	api := suite.router.Group("/api")
	api.Use(middleware.AuthMiddleware(suite.authService))
	{
		api.GET("/me", func(c *gin.Context) {
			user, exists := c.Get("user")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusOK, user)
		})
	}
}

func (suite *AuthIntegrationTestSuite) TearDownTest() {
	// Limpiar base de datos después de cada prueba
	suite.db.Exec("TRUNCATE TABLE users CASCADE")
}

func (suite *AuthIntegrationTestSuite) TearDownSuite() {
	// Cerrar conexión a la base de datos
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// Test Suite para el registro de usuarios
type RegisterTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *RegisterTestSuite) SetupTest() {
	suite.router = setupTestRouter()
}

func (suite *RegisterTestSuite) TestRegisterSuccess() {
	payload := map[string]interface{}{
		"email":     "test@example.com",
		"password":  "password123",
		"firstName": "Test",
		"lastName":  "User",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

// Test Suite para el login
type LoginTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *LoginTestSuite) SetupTest() {
	suite.router = setupTestRouter()
	// Aquí podrías crear un usuario de prueba si es necesario
}

func (suite *LoginTestSuite) TestLoginSuccess() {
	// Primero registramos un usuario
	registerTestUser(suite.router)

	payload := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// Test Suite para rutas protegidas
type ProtectedRouteTestSuite struct {
	suite.Suite
	router      *gin.Engine
	accessToken string
}

func (suite *ProtectedRouteTestSuite) SetupTest() {
	suite.router = setupTestRouter()
	// Obtenemos un token válido para las pruebas
	suite.accessToken = getTestToken(suite.router)
}

func (suite *ProtectedRouteTestSuite) TestAccessProtectedRoute() {
	req, _ := http.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", suite.accessToken))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// Test Suite para refresh token
type RefreshTokenTestSuite struct {
	suite.Suite
	router       *gin.Engine
	refreshToken string
}

func (suite *RefreshTokenTestSuite) SetupTest() {
	suite.router = setupTestRouter()
	// Obtenemos un refresh token válido para las pruebas
	_, suite.refreshToken = getTestTokens(suite.router)
}

func (suite *RefreshTokenTestSuite) TestRefreshToken() {
	payload := map[string]interface{}{
		"refreshToken": suite.refreshToken,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// Funciones de ayuda
func setupTestRouter() *gin.Engine {
	// Configuración de la base de datos
	os.Setenv("APP_ENV", "test")
	cfg, _ := config.Load()
	db, _ := database.Connect(cfg)
	_ = database.AutoMigrate(db)

	// Configuración de servicios
	userRepo := repositories.NewUserRepositoryImpl(db)
	authService := services.NewAuthService(cfg, userRepo)

	// Configuración del router
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Rutas públicas
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", handlers.NewAuthHandler(authService).Register)
		authGroup.POST("/login", handlers.NewAuthHandler(authService).Login)
		authGroup.POST("/refresh", handlers.NewAuthHandler(authService).RefreshToken)
	}

	// Rutas protegidas
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authService))
	{
		api.GET("/me", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(http.StatusOK, user)
		})
	}

	return r
}

func registerTestUser(router *gin.Engine) {
	payload := map[string]interface{}{
		"email":     "test@example.com",
		"password":  "password123",
		"firstName": "Test",
		"lastName":  "User",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func getTestToken(router *gin.Engine) string {
	// Primero registramos un usuario si no existe
	registerTestUser(router)

	// Luego hacemos login
	loginPayload := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	return response["data"].(map[string]interface{})["accessToken"].(string)
}

func getTestTokens(router *gin.Engine) (string, string) {
	// Primero registramos un usuario si no existe
	registerTestUser(router)

	// Luego hacemos login
	loginPayload := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})

	return data["accessToken"].(string), data["refreshToken"].(string)
}

// Funciones para ejecutar los tests
func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func TestProtectedRouteSuite(t *testing.T) {
	suite.Run(t, new(ProtectedRouteTestSuite))
}

func TestRefreshTokenSuite(t *testing.T) {
	suite.Run(t, new(RefreshTokenTestSuite))
}

func TestInvalidTokenSuite(t *testing.T) {
	suite.Run(t, new(InvalidTokenTestSuite))
}

func TestAuthIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}
