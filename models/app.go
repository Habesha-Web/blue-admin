package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// App Database model info
// @Description App type information
type App struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	UUID        string `gorm:"constraint:not null; unique; type:string;" json:"uuid"`
	Active      bool   `gorm:"constraint:not null;" json:"active"`
	Description string `gorm:"not null;" json:"description,omitempty"`
	Roles       []Role `gorm:"association_foreignkey:AppID constraint:OnUpdate:SET NULL OnDelete:SET NULL" json:"roles,omitempty"`
}

func (app *App) BeforeCreate(tx *gorm.DB) (err error) {
	gen, _ := uuid.NewV7()
	id := gen.String()
	app.UUID = id
	return
}

// AppPost model info
// @Description AppPost type information
type AppPost struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Description string `gorm:"not null;" json:"description,omitempty"`
	Active      bool   `gorm:"constraint:not null;" json:"active"`
}

// AppGet model info
// @Description AppGet type information
type AppGet struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	UUID        string `gorm:"constraint:not null; unique; type:string;" json:"uuid"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
	Description string `gorm:"not null;" json:"description,omitempty"`
	Roles       []Role `gorm:"association_foreignkey:AppID constraint:OnUpdate:SET NULL OnDelete:SET NULL" json:"roles,omitempty"`
}

// AppPut model info
// @Description AppPut type information
type AppPut struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name string `gorm:"not null; unique;" json:"name,omitempty"`

	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
	Description string `gorm:"not null;" json:"description,omitempty"`
}

// AppPatch model info
// @Description AppPatch type information
type AppPatch struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
	Description string `gorm:"not null;" json:"description,omitempty"`
}
