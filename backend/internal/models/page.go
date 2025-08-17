package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Page struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	BoardID   uuid.UUID `gorm:"type:uuid;not null;index" json:"board_id"`
	Title     string    `gorm:"not null" json:"title"`
	Date      time.Time `gorm:"not null;index" json:"date"`
	OrderIdx  int       `gorm:"not null;index" json:"order_idx"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Elements  []Element `gorm:"foreignKey:PageID" json:"elements,omitempty"`
}

func (p *Page) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
