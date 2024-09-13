package models

import (
	"database/sql"
)

// Endpoint Database model info
// @Description App type information
type Endpoint struct {
	ID          uint          `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string        `gorm:"not null;unique" json:"name,omitempty"`
	RoutePath   string        `gorm:"not null;" json:"route_path,omitempty"`
	Method      string        `gorm:"not null;" json:"method,omitempty"`
	Description string        `gorm:"not null;" json:"description,omitempty"`
	FeatureID   sql.NullInt64 `gorm:"foreignkey:FeatureID default:NULL;,OnDelete:SET NULL;" json:"feature_id,omitempty" swaggertype:"number"`
}

// EndpointPost model info
// @Description EndpointPost type information
type EndpointPost struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	RoutePath   string `gorm:"not null; unique;" json:"route_path,omitempty"`
	Method      string `gorm:"not null; unique;" json:"method,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
}

// EndpointGet model info
// @Description EndpointGet type information
type EndpointGet struct {
	ID          uint          `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string        `gorm:"not null; unique;" json:"name,omitempty"`
	RoutePath   string        `gorm:"not null; unique;" json:"route_path,omitempty"`
	Method      string        `gorm:"not null; unique;" json:"method,omitempty"`
	Description string        `gorm:"not null; unique;" json:"description,omitempty"`
	FeatureID   sql.NullInt64 `gorm:"foreignkey:FeatureID default:NULL;,OnDelete:SET NULL;" json:"feature_id,omitempty" swaggertype:"number"`
}

// EndpointPut model info
// @Description EndpointPut type information
type EndpointPut struct {
	Name        string        `gorm:"not null; unique;" json:"name,omitempty"`
	RoutePath   string        `gorm:"not null; unique;" json:"route_path,omitempty"`
	Method      string        `gorm:"not null; unique;" json:"method,omitempty"`
	Description string        `gorm:"not null; unique;" json:"description,omitempty"`
	FeatureID   sql.NullInt64 `gorm:"foreignkey:FeatureID default:NULL;,OnDelete:SET NULL;" json:"feature_id,omitempty" swaggertype:"number"`
}

// EndpointPatch model info
// @Description EndpointPatch type information
type EndpointPatch struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	RoutePath   string `gorm:"not null; unique;" json:"route_path,omitempty"`
	Method      string `gorm:"not null; unique;" json:"method,omitempty"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
}
