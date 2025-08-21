package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Element struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	PageID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"page_id"`
	Kind      string         `gorm:"not null;check:kind IN ('text','image','sticker','shape')" json:"kind"`
	X         float64        `gorm:"not null" json:"x"`
	Y         float64        `gorm:"not null" json:"y"`
	W         float64        `gorm:"not null" json:"w"`
	H         float64        `gorm:"not null" json:"h"`
	Rotation  float64        `gorm:"default:0" json:"rotation"`
	Z         int            `gorm:"default:0;index" json:"z"`
	Visible   bool           `gorm:"default:true" json:"visible"`
	Locked    bool           `gorm:"default:false" json:"locked"`
	Payload   datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (e *Element) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
