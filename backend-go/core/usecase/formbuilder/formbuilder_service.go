package formbuilder

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

var (
	ErrInvalidFieldType     = errors.New("invalid form field type: must be text, url, select, radio, checkbox, or textarea")
	ErrFieldNameNumeric     = errors.New("field name cannot start with a number")
	ErrFieldNameEmpty       = errors.New("field name cannot be empty")
	ErrDuplicateFieldName   = errors.New("a field with this name already exists in the form")
	ErrRequiredFieldMissing = errors.New("required form field is missing")
	ErrInvalidURLField      = errors.New("invalid URL format provided")
)

var validFieldTypes = map[string]bool{
	"text":     true,
	"url":      true,
	"select":   true,
	"radio":    true,
	"checkbox": true,
	"textarea": true,
}

type CreateFormDTO struct {
	Name  string           `json:"name"`
	Style domain.FormStyle `json:"style"`
}

type UpdateFormDTO struct {
	Name  string           `json:"name"`
	Style domain.FormStyle `json:"style"`
}

type FormFieldDTO struct {
	Name         string                 `json:"name"`
	Label        string                 `json:"label"`
	HideLabel    bool                   `json:"hide_label"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"`
	DefaultValue string                 `json:"default_value"`
	Required     bool                   `json:"required"`
	Hidden       bool                   `json:"hidden"`
	Readonly     bool                   `json:"readonly"`
	Options      map[string]interface{} `json:"options"`
	Prefix       string                 `json:"prefix"`
	Suffix       string                 `json:"suffix"`
	TextSize     int                    `json:"text_size"`
}

type FormbuilderService struct {
	repo domain.FormbuilderRepository
}

func NewFormbuilderService(repo domain.FormbuilderRepository) *FormbuilderService {
	return &FormbuilderService{repo: repo}
}

func (s *FormbuilderService) Slugify(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", ErrFieldNameEmpty
	}

	// Replace non-alphanumeric with underscore
	var b strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune('_')
		}
	}

	res := regexp.MustCompile(`_+`).ReplaceAllString(b.String(), "_")
	res = strings.Trim(res, "_")

	if res == "" {
		return "", ErrFieldNameEmpty
	}

	if unicode.IsDigit(rune(res[0])) {
		return "", ErrFieldNameNumeric
	}

	return res, nil
}

func (s *FormbuilderService) GetForm(ctx context.Context, id int64) (*domain.Form, error) {
	return s.repo.GetFormByID(ctx, id)
}

func (s *FormbuilderService) ListForms(ctx context.Context, limit, offset int) ([]*domain.Form, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListForms(ctx, limit, offset)
}

func (s *FormbuilderService) CreateForm(ctx context.Context, dto CreateFormDTO) (*domain.Form, error) {
	if strings.TrimSpace(dto.Name) == "" {
		return nil, fmt.Errorf("%w: form name is required", appErrors.ErrInvalidInput)
	}

	styleType := strings.ToLower(strings.TrimSpace(dto.Style.Type))
	if styleType == "" {
		styleType = "horizontal"
	}

	form := &domain.Form{
		Name: strings.TrimSpace(dto.Name),
		Style: domain.FormStyle{
			Type:      styleType,
			ShowTitle: dto.Style.ShowTitle,
		},
	}

	if err := s.repo.CreateForm(ctx, form); err != nil {
		return nil, err
	}
	return form, nil
}

func (s *FormbuilderService) UpdateForm(ctx context.Context, id int64, dto UpdateFormDTO) (*domain.Form, error) {
	form, err := s.repo.GetFormByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(dto.Name) != "" {
		form.Name = strings.TrimSpace(dto.Name)
	}
	if dto.Style.Type != "" {
		form.Style.Type = strings.ToLower(strings.TrimSpace(dto.Style.Type))
	}
	form.Style.ShowTitle = dto.Style.ShowTitle

	if err := s.repo.UpdateForm(ctx, form); err != nil {
		return nil, err
	}
	return form, nil
}

func (s *FormbuilderService) DeleteForm(ctx context.Context, id int64) error {
	return s.repo.DeleteForm(ctx, id)
}

func (s *FormbuilderService) AddField(ctx context.Context, formID int64, dto FormFieldDTO) (*domain.FormField, error) {
	_, err := s.repo.GetFormByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	fieldType := strings.ToLower(strings.TrimSpace(dto.Type))
	if !validFieldTypes[fieldType] {
		return nil, ErrInvalidFieldType
	}

	fieldCount, _ := s.repo.CountFields(ctx, formID)
	fieldNumber := fieldCount + 1

	name := dto.Name
	if strings.TrimSpace(name) == "" {
		name = fmt.Sprintf("new_%s_%d", fieldType, fieldNumber)
	}

	slugName, err := s.Slugify(name)
	if err != nil {
		return nil, err
	}

	// Check duplicate field name in form
	existingFields, _ := s.repo.GetFieldsByFormID(ctx, formID)
	for _, ef := range existingFields {
		if ef.Name == slugName {
			return nil, ErrDuplicateFieldName
		}
	}

	label := dto.Label
	if strings.TrimSpace(label) == "" {
		label = fmt.Sprintf("%s %d", strings.Title(fieldType), fieldNumber)
	}

	options := dto.Options
	if (fieldType == "select" || fieldType == "radio" || fieldType == "checkbox") && len(options) == 0 {
		options = map[string]interface{}{
			"Option 1": "1",
			"Option 2": "2",
		}
	}

	field := &domain.FormField{
		FormID:       formID,
		Name:         slugName,
		Label:        label,
		HideLabel:    dto.HideLabel,
		Description:  dto.Description,
		Type:         fieldType,
		DefaultValue: dto.DefaultValue,
		Required:     dto.Required,
		Hidden:       dto.Hidden,
		Readonly:     dto.Readonly,
		Options:      options,
		Prefix:       dto.Prefix,
		Suffix:       dto.Suffix,
		TextSize:     dto.TextSize,
	}

	if err := s.repo.AddField(ctx, field); err != nil {
		return nil, err
	}
	return field, nil
}

func (s *FormbuilderService) UpdateField(ctx context.Context, fieldID int64, dto FormFieldDTO) (*domain.FormField, error) {
	field, err := s.repo.GetFieldByID(ctx, fieldID)
	if err != nil {
		return nil, err
	}

	if dto.Type != "" {
		fieldType := strings.ToLower(strings.TrimSpace(dto.Type))
		if !validFieldTypes[fieldType] {
			return nil, ErrInvalidFieldType
		}
		field.Type = fieldType
	}

	if strings.TrimSpace(dto.Name) != "" {
		slugName, err := s.Slugify(dto.Name)
		if err != nil {
			return nil, err
		}

		if slugName != field.Name {
			// Check if name is taken by other field in form
			existingFields, _ := s.repo.GetFieldsByFormID(ctx, field.FormID)
			for _, ef := range existingFields {
				if ef.ID != fieldID && ef.Name == slugName {
					return nil, ErrDuplicateFieldName
				}
			}
			field.Name = slugName
		}
	}

	if strings.TrimSpace(dto.Label) != "" {
		field.Label = strings.TrimSpace(dto.Label)
	}
	field.HideLabel = dto.HideLabel
	field.Description = dto.Description
	field.DefaultValue = dto.DefaultValue
	field.Required = dto.Required
	field.Hidden = dto.Hidden
	field.Readonly = dto.Readonly
	if dto.Options != nil {
		field.Options = dto.Options
	}
	field.Prefix = dto.Prefix
	field.Suffix = dto.Suffix
	if dto.TextSize > 0 {
		field.TextSize = dto.TextSize
	}

	if err := s.repo.UpdateField(ctx, field); err != nil {
		return nil, err
	}
	return field, nil
}

func (s *FormbuilderService) DeleteField(ctx context.Context, fieldID int64) error {
	return s.repo.DeleteField(ctx, fieldID)
}

// ValidateFormSubmission validates client submitted values against form fields
func (s *FormbuilderService) ValidateFormSubmission(ctx context.Context, formID int64, submitted map[string]interface{}) (map[string]interface{}, error) {
	fields, err := s.repo.GetFieldsByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}

	cleaned := make(map[string]interface{})

	for _, f := range fields {
		val, exists := submitted[f.Name]
		strVal := ""
		if exists && val != nil {
			strVal = fmt.Sprintf("%v", val)
		}

		// Check required
		if f.Required && (!exists || strings.TrimSpace(strVal) == "") {
			return nil, fmt.Errorf("%w: '%s' (%s)", ErrRequiredFieldMissing, f.Label, f.Name)
		}

		if !exists || val == nil {
			if f.DefaultValue != "" {
				cleaned[f.Name] = f.DefaultValue
			}
			continue
		}

		// Type-specific validations
		if f.Type == "url" && strings.TrimSpace(strVal) != "" {
			u, err := url.ParseRequestURI(strVal)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
				return nil, fmt.Errorf("%w for field '%s'", ErrInvalidURLField, f.Name)
			}
		}

		if f.Type == "checkbox" {
			if b, ok := val.(bool); ok {
				cleaned[f.Name] = b
			} else if strVal == "1" || strVal == "true" || strVal == "on" {
				cleaned[f.Name] = true
			} else {
				cleaned[f.Name] = false
			}
			continue
		}

		cleaned[f.Name] = val
	}

	return cleaned, nil
}
