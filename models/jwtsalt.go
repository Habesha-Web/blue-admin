package models

// JWTSalt Database model info
// @Description App type information
type JWTSalt struct {
	ID    uint   `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	SaltA string `gorm:"not null; unique;" json:"salt_a,omitempty"`
	SaltB string `gorm:"not null; unique;" json:"salt_b,omitempty"`
}
