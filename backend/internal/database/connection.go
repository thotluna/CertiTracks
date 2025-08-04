package database

import (
	"fmt"
	"log"
	"os"

	"certitrack/internal/config"
	"certitrack/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseDSN()

	var gormLogger logger.Interface
	if os.Getenv("GO_ENV") != "test" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	if os.Getenv("GO_ENV") != "test" {
		log.Println("✅ Database connection established")
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	if os.Getenv("GO_ENV") != "test" {
		log.Println("🔄 Running database migrations...")
	}

	err := db.AutoMigrate(
		&models.User{},
		&models.PasswordResetToken{},
		// Add other models here as they are created
		// &models.Person{},
		// &models.Equipment{},
		// &models.Certification{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if os.Getenv("GO_ENV") != "test" {
		log.Println("✅ Database migrations completed")
	}
	return nil
}

func CreateDefaultAdmin(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check for admin users: %w", err)
	}

	if count > 0 {
		if os.Getenv("GO_ENV") != "test" {
			log.Println("ℹ️  Admin user already exists, skipping creation")
		}
		return nil
	}

	admin := models.User{
		Email:     "admin@certitrack.local",
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password: "password"
		FirstName: "System",
		LastName:  "Administrator",
		Role:      "admin",
		IsActive:  true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	if os.Getenv("GO_ENV") != "test" {
		log.Println("✅ Default admin user created (admin@certitrack.local / password)")
	}
	return nil
}
