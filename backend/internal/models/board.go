package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Skin        string    `gorm:"default:'default'" json:"skin"`
	EditToken   uuid.UUID `gorm:"type:uuid;unique;not null" json:"-"`
	PublicToken uuid.UUID `gorm:"type:uuid;unique;not null" json:"public_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Pages       []Page    `gorm:"foreignKey:BoardID" json:"pages,omitempty"`
}

func (b *Board) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if b.EditToken == uuid.Nil {
		b.EditToken = uuid.New()
	}
	if b.PublicToken == uuid.Nil {
		b.PublicToken = uuid.New()
	}
	return nil
}
