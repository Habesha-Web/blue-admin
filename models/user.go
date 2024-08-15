package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Database model info
// @Description App type information
type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Email         string    `gorm:"not null; unique;" json:"email,omitempty"`
	Password      string    `gorm:"not null;" json:"password,omitempty"`
	DateRegistred time.Time `gorm:"constraint:not null; default:current_timestamp;" json:"date_registered,omitempty"`
	Disabled      bool      `gorm:"default:true; constraint:not null;" json:"disabled"`
	UUID          string    `gorm:"constraint:not null; unique; type:string;" json:"uuid"`
	Roles         []Role    `gorm:"many2many:user_roles; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;" json:"roles,omitempty"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	gen, _ := uuid.NewV7()
	id := gen.String()
	user.UUID = id
	user.Password = hashfunc(user.Password)
	return
}

// UserPost model info
// @Description UserPost type information
type UserPost struct {
	Email    string `gorm:"not null; unique;" json:"email,omitempty"`
	Password string `gorm:"not null;" json:"password,omitempty"`
	Disabled bool   `gorm:"default:true; constraint:not null;" json:"disabled"`
}

// UserGet model info
// @Description UserGet type information
type UserGet struct {
	ID            uint      `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Email         string    `gorm:"not null; unique;" json:"email,omitempty"`
	DateRegistred time.Time `gorm:"constraint:not null; default:current_timestamp;" json:"date_registered,omitempty"`
	Disabled      bool      `gorm:"default:true; constraint:not null;" json:"disabled"`
	UUID          string    `gorm:"constraint:not null; unique; type:string;" json:"uuid"`
	Roles         []Role    `gorm:"many2many:user_roles; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;" json:"roles,omitempty"`
}

// UserPut model info
// @Description UserPut type information
type UserPut struct {
	Email         string    `gorm:"not null; unique;" json:"email,omitempty"`
	Password      string    `gorm:"not null;" json:"password,omitempty"`
	DateRegistred time.Time `gorm:"constraint:not null; default:current_timestamp;" json:"date_registered,omitempty"`
	Disabled      bool      `gorm:"default:true; constraint:not null;" json:"disabled"`
}

// UserPatch model info
// @Description UserPatch type information
type UserPatch struct {
	Email string `gorm:"not null; unique;" json:"email,omitempty"`

	DateRegistred time.Time `gorm:"constraint:not null; default:current_timestamp;" json:"date_registered,omitempty"`
	Disabled      bool      `gorm:"default:true; constraint:not null;" json:"disabled"`
}
