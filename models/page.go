package models

// Page Database model info
// @Description App type information
type Page struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Roles       []Role `gorm:"many2many:page_roles; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;" json:"roles,omitempty"`
}

// PagePost model info
// @Description PagePost type information
type PagePost struct {
	Name string `gorm:"not null; unique;" json:"name,omitempty"`

	Description string `gorm:"not null; unique;" json:"description,omitempty"`
}

// PageGet model info
// @Description PageGet type information
type PageGet struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
	Roles       []Role `gorm:"many2many:page_roles; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;" json:"roles,omitempty"`
}

// PagePut model info
// @Description PagePut type information
type PagePut struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
}

// PagePatch model info
// @Description PagePatch type information
type PagePatch struct {
	Name        string `gorm:"not null; unique;" json:"name,omitempty"`
	Active      bool   `gorm:"default:true; constraint:not null;" json:"active"`
	Description string `gorm:"not null; unique;" json:"description,omitempty"`
}
