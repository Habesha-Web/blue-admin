package models

import (
	"database/sql"
)

// Role Database model info
// @Description App type information
type Role struct {
	ID          uint          `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string        `gorm:"not null; unique;" json:"name,omitempty"`
	Description string        `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool          `gorm:"default:true; constraint:not null;" json:"active"`
	Users       []User        `gorm:"many2many:user_roles; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;" json:"users,omitempty"`
	Features    []Feature     `gorm:"foreignkey:RoleID; constraint:OnUpdate:CASCADE; OnDelete:SET NULL;" json:"features,omitempty"`
	Pages       []Page        `gorm:"many2many:page_roles; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"pages,omitempty"`
	AppID       sql.NullInt64 `gorm:"foreignkey:AppID OnDelete:SET NULL" json:"app,omitempty" swaggertype:"number"`
}

// RolePost model info
// @Description RolePost type information
type RolePost struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
}

// RoleGet model info
// @Description RoleGet type information
type RoleGet struct {
	ID          uint      `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string    `gorm:"not null; unique;" json:"name,omitempty"`
	Description string    `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool      `gorm:"default:true; constraint:not null;" json:"active"`
	Users       []User    `gorm:"many2many:user_roles; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;" json:"users,omitempty"`
	Features    []Feature `gorm:"foreignkey:RoleID; constraint:OnUpdate:CASCADE; OnDelete:SET NULL;" json:"features,omitempty"`
	Pages       []Page    `gorm:"many2many:page_roles; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"pages,omitempty"`
}

// RolePut model info
// @Description RolePut type information
type RolePut struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
}

// RolePatch model info
// @Description RolePatch type information
type RolePatch struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
}
