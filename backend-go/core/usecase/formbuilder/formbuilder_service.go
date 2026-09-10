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
	ErrInvalidFieldType     = errors.New("invalid form field type")
	ErrFieldNameNumeric     = errors.New("field name cannot start with a number")
	ErrFieldNameEmpty       = errors.New("field name cannot be empty")
	ErrDuplicateFieldName   = errors.New("duplicate field name")
	ErrRequiredFieldMissing = errors.New("required field missing")
	ErrInvalidURLField      = errors.New("invalid URL format")
	validFieldTypes         = map[string]bool{"text": true, "url": true, "select": true, "radio": true, "checkbox": true, "textarea": true}
)

type CreateFormDTO struct { Name string `json:"name"`; Style domain.FormStyle `json:"style"` }
type UpdateFormDTO struct { Name string `json:"name"`; Style domain.FormStyle `json:"style"` }
type FormFieldDTO struct { Name, Label, Description, Type, DefaultValue, Prefix, Suffix string; HideLabel, Required, Hidden, Readonly bool; Options map[string]interface{}; TextSize int }

type FormbuilderService struct { repo domain.FormbuilderRepository }

func NewFormbuilderService(r domain.FormbuilderRepository) *FormbuilderService { return &FormbuilderService{r} }

func (s *FormbuilderService) Slugify(t string) (string, error) {
	t = strings.TrimSpace(t); if t == "" { return "", ErrFieldNameEmpty }
	var b strings.Builder
	for _, r := range t { if unicode.IsLetter(r) || unicode.IsDigit(r) { b.WriteRune(unicode.ToLower(r)) } else { b.WriteRune('_') } }
	res := regexp.MustCompile(`_+`).ReplaceAllString(b.String(), "_"); res = strings.Trim(res, "_")
	if res == "" { return "", ErrFieldNameEmpty }
	if unicode.IsDigit(rune(res[0])) { return "", ErrFieldNameNumeric }
	return res, nil
}

func (s *FormbuilderService) GetForm(ctx context.Context, id int64) (*domain.Form, error) { return s.repo.GetFormByID(ctx, id) }
func (s *FormbuilderService) ListForms(ctx context.Context, l, o int) ([]*domain.Form, int, error) {
	if l <= 0 { l = 20 }; return s.repo.ListForms(ctx, l, o)
}

func (s *FormbuilderService) CreateForm(ctx context.Context, dto CreateFormDTO) (*domain.Form, error) {
	if strings.TrimSpace(dto.Name) == "" { return nil, appErrors.ErrInvalidInput }
	f := &domain.Form{Name: strings.TrimSpace(dto.Name), Style: domain.FormStyle{Type: strings.ToLower(dto.Style.Type), ShowTitle: dto.Style.ShowTitle}}
	if f.Style.Type == "" { f.Style.Type = "horizontal" }
	return f, s.repo.CreateForm(ctx, f)
}

func (s *FormbuilderService) UpdateForm(ctx context.Context, id int64, dto UpdateFormDTO) (*domain.Form, error) {
	f, err := s.repo.GetFormByID(ctx, id); if err != nil { return nil, err }
	if dto.Name != "" { f.Name = dto.Name }
	if dto.Style.Type != "" { f.Style.Type = strings.ToLower(dto.Style.Type) }
	f.Style.ShowTitle = dto.Style.ShowTitle
	return f, s.repo.UpdateForm(ctx, f)
}

func (s *FormbuilderService) AddField(ctx context.Context, fID int64, dto FormFieldDTO) (*domain.FormField, error) {
	tp := strings.ToLower(strings.TrimSpace(dto.Type)); if !validFieldTypes[tp] { return nil, ErrInvalidFieldType }
	cnt, _ := s.repo.CountFields(ctx, fID); num := cnt + 1
	nm := dto.Name; if nm == "" { nm = fmt.Sprintf("new_%s_%d", tp, num) }
	sl, err := s.Slugify(nm); if err != nil { return nil, err }
	ex, _ := s.repo.GetFieldsByFormID(ctx, fID); for _, e := range ex { if e.Name == sl { return nil, ErrDuplicateFieldName } }
	lb := dto.Label; if lb == "" { lb = strings.Title(tp) + fmt.Sprintf(" %d", num) }
	opt := dto.Options; if (tp == "select" || tp == "radio" || tp == "checkbox") && len(opt) == 0 { opt = map[string]interface{}{"Option 1": "1", "Option 2": "2"} }
	fd := &domain.FormField{FormID: fID, Name: sl, Label: lb, HideLabel: dto.HideLabel, Description: dto.Description, Type: tp, DefaultValue: dto.DefaultValue, Required: dto.Required, Hidden: dto.Hidden, Readonly: dto.Readonly, Options: opt, Prefix: dto.Prefix, Suffix: dto.Suffix, TextSize: dto.TextSize}
	return fd, s.repo.AddField(ctx, fd)
}

func (s *FormbuilderService) UpdateField(ctx context.Context, id int64, dto FormFieldDTO) (*domain.FormField, error) {
	f, err := s.repo.GetFieldByID(ctx, id); if err != nil { return nil, err }
	if dto.Type != "" { if !validFieldTypes[strings.ToLower(dto.Type)] { return nil, ErrInvalidFieldType }; f.Type = strings.ToLower(dto.Type) }
	if dto.Name != "" {
		sl, _ := s.Slugify(dto.Name); if sl != f.Name {
			ex, _ := s.repo.GetFieldsByFormID(ctx, f.FormID); for _, e := range ex { if e.ID != id && e.Name == sl { return nil, ErrDuplicateFieldName } }
			f.Name = sl
		}
	}
	if dto.Label != "" { f.Label = dto.Label }; f.HideLabel, f.Description, f.DefaultValue, f.Required, f.Hidden, f.Readonly, f.Prefix, f.Suffix = dto.HideLabel, dto.Description, dto.DefaultValue, dto.Required, dto.Hidden, dto.Readonly, dto.Prefix, dto.Suffix
	if dto.Options != nil { f.Options = dto.Options }; if dto.TextSize > 0 { f.TextSize = dto.TextSize }
	return f, s.repo.UpdateField(ctx, f)
}

func (s *FormbuilderService) DeleteForm(ctx context.Context, id int64) error { return s.repo.DeleteForm(ctx, id) }
func (s *FormbuilderService) DeleteField(ctx context.Context, id int64) error { return s.repo.DeleteField(ctx, id) }

func (s *FormbuilderService) ValidateFormSubmission(ctx context.Context, fID int64, sub map[string]interface{}) (map[string]interface{}, error) {
	fs, err := s.repo.GetFieldsByFormID(ctx, fID); if err != nil { return nil, err }
	cl := make(map[string]interface{})
	for _, f := range fs {
		v, ex := sub[f.Name]; sv := fmt.Sprintf("%v", v); if !ex || v == nil { sv = "" }
		if f.Required && (sv == "") { return nil, fmt.Errorf("%w: %s", ErrRequiredFieldMissing, f.Label) }
		if !ex || v == nil { if f.DefaultValue != "" { cl[f.Name] = f.DefaultValue }; continue }
		if f.Type == "url" && sv != "" { if u, err := url.ParseRequestURI(sv); err != nil || (u.Scheme != "http" && u.Scheme != "https") { return nil, ErrInvalidURLField } }
		if f.Type == "checkbox" { if b, ok := v.(bool); ok { cl[f.Name] = b } else { cl[f.Name] = (sv == "1" || sv == "true" || sv == "on") }; continue }
		cl[f.Name] = v
	}
	return cl, nil
}
