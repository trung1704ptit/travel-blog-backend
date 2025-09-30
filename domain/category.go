package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug" validate:"required"`
	Description string     `json:"description,omitempty"`
	Image       string     `json:"image,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Parent      *Category  `json:"parent,omitempty"`
	Children    []Category `json:"children,omitempty"`
	Level       int        `json:"level,omitempty"`
	Path        string     `json:"path,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
