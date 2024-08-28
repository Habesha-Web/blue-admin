package models

import (
	"database/sql"
)

// Feature Database model info
// @Description App type information
type Feature struct {
	ID          uint          `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string        `gorm:"not null; unique;" json:"name,omitempty"`
	Description string        `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool          `gorm:"default:true; constraint:not null;" json:"active"`
	RoleID      sql.NullInt64 `gorm:"foreignkey:RoleID OnDelete:SET NULL" json:"role,omitempty" swaggertype:"number"`
	Endpoints   []Endpoint    `gorm:"association_foreignkey:FeatureID constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"endpoints,omitempty"`
}

// FeaturePost model info
// @Description FeaturePost type information
type FeaturePost struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
}

// FeatureGet model info
// @Description FeatureGet type information
type FeatureGet struct {
	ID          uint       `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string     `gorm:"not null; unique;" json:"name,omitempty"`
	Description string     `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool       `gorm:"default:true; constraint:not null;" json:"active"`
	Endpoints   []Endpoint `gorm:"association_foreignkey:FeatureID constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"endpoints,omitempty"`
}

// FeaturePut model info
// @Description FeaturePut type information
type FeaturePut struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
}

// FeaturePatch model info
// @Description FeaturePatch type information
type FeaturePatch struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
}
