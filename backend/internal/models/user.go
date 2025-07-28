package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string     `gorm:"uniqueIndex;not null" json:"email"`
	Password  string     `gorm:"not null" json:"-"` // Never include password in JSON
	FirstName string     `gorm:"not null" json:"firstName"`
	LastName  string     `gorm:"not null" json:"lastName"`
	Phone     string     `json:"phone"`
	Role      string     `gorm:"not null;default:'user'" json:"role"` // 'admin' or 'user'
	IsActive  bool       `gorm:"default:true" json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	LastLogin *time.Time `json:"lastLogin"`

	// Relationships (will be uncommented when other models are created)
	// CreatedPeople         []Person        `gorm:"foreignKey:CreatedBy" json:"-"`
	// CreatedEquipment      []Equipment     `gorm:"foreignKey:CreatedBy" json:"-"`
	// CreatedCertifications []Certification `gorm:"foreignKey:CreatedBy" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}
