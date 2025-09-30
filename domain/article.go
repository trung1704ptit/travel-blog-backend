package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Article is representing the Article data struct
type Article struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title" validate:"required"`
	Slug        string          `json:"slug" validate:"required,alphanumdash"`
	Content     string          `json:"content" validate:"required"`
	Thumbnail   string          `json:"thumbnail" validate:"omitempty,url"`
	Image       string          `json:"image" validate:"omitempty,url"`
	Description string          `json:"description"`
	Keywords    JSONStringSlice `json:"keywords"`
	Categories  []Category      `json:"categories"`
	Author      Author          `json:"author"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedAt   time.Time       `json:"created_at"`
}

type ArticleCategory struct {
	ID         uuid.UUID `json:"id"`
	ArticleID  uuid.UUID `json:"article_id"`
	CategoryID uuid.UUID `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// JSONStringSlice is a custom type that handles JSON marshaling/unmarshaling for string slices
type JSONStringSlice []string

// Scan implements the sql.Scanner interface for database/sql
func (j *JSONStringSlice) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface for database/sql
func (j JSONStringSlice) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
