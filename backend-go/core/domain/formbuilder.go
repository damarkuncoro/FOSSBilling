package domain

import (
	"context"
	"time"
)

type FormStyle struct {
	Type      string `json:"type"`       // "horizontal", "vertical", "default"
	ShowTitle bool   `json:"show_title"` // whether to render title in customer checkout
}

type FormField struct {
	ID           int64                  `json:"id"`
	FormID       int64                  `json:"form_id"`
	Name         string                 `json:"name"`
	Label        string                 `json:"label"`
	HideLabel    bool                   `json:"hide_label"`
	Description  string                 `json:"description,omitempty"`
	Type         string                 `json:"type"` // "text", "url", "select", "radio", "checkbox", "textarea"
	DefaultValue string                 `json:"default_value,omitempty"`
	Required     bool                   `json:"required"`
	Hidden       bool                   `json:"hidden"`
	Readonly     bool                   `json:"readonly"`
	Options      map[string]interface{} `json:"options,omitempty"` // For select, radio, checkbox, textarea dimensions
	Prefix       string                 `json:"prefix,omitempty"`
	Suffix       string                 `json:"suffix,omitempty"`
	TextSize     int                    `json:"text_size,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type Form struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	Style     FormStyle    `json:"style"`
	Fields    []*FormField `json:"fields,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type FormbuilderRepository interface {
	GetFormByID(ctx context.Context, id int64) (*Form, error)
	ListForms(ctx context.Context, limit, offset int) ([]*Form, int, error)
	CreateForm(ctx context.Context, form *Form) error
	UpdateForm(ctx context.Context, form *Form) error
	DeleteForm(ctx context.Context, id int64) error

	GetFieldByID(ctx context.Context, fieldID int64) (*FormField, error)
	GetFieldsByFormID(ctx context.Context, formID int64) ([]*FormField, error)
	AddField(ctx context.Context, field *FormField) error
	UpdateField(ctx context.Context, field *FormField) error
	DeleteField(ctx context.Context, fieldID int64) error
	CountFields(ctx context.Context, formID int64) (int, error)
}
