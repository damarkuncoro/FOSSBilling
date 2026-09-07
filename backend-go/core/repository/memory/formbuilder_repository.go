package memory

import (
	"context"
	"sync"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
)

type MockFormbuilderRepository struct {
	mu          sync.RWMutex
	forms       map[int64]*domain.Form
	fields      map[int64]*domain.FormField
	nextFormID  int64
	nextFieldID int64
}

func NewMockFormbuilderRepository() *MockFormbuilderRepository {
	return &MockFormbuilderRepository{
		forms:       make(map[int64]*domain.Form),
		fields:      make(map[int64]*domain.FormField),
		nextFormID:  1,
		nextFieldID: 1,
	}
}

func (r *MockFormbuilderRepository) GetFormByID(ctx context.Context, id int64) (*domain.Form, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	form, exists := r.forms[id]
	if !exists {
		return nil, appErrors.ErrNotFound
	}

	// Clone and attach fields
	copyForm := *form
	copyForm.Fields = make([]*domain.FormField, 0)
	for _, f := range r.fields {
		if f.FormID == id {
			fieldCopy := *f
			copyForm.Fields = append(copyForm.Fields, &fieldCopy)
		}
	}
	return &copyForm, nil
}

func (r *MockFormbuilderRepository) ListForms(ctx context.Context, limit, offset int) ([]*domain.Form, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*domain.Form
	for _, f := range r.forms {
		copyForm := *f
		copyForm.Fields = make([]*domain.FormField, 0)
		for _, fld := range r.fields {
			if fld.FormID == f.ID {
				fieldCopy := *fld
				copyForm.Fields = append(copyForm.Fields, &fieldCopy)
			}
		}
		list = append(list, &copyForm)
	}

	total := len(list)
	if offset >= total {
		return []*domain.Form{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return list[offset:end], total, nil
}

func (r *MockFormbuilderRepository) CreateForm(ctx context.Context, form *domain.Form) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	form.ID = r.nextFormID
	r.nextFormID++
	if form.CreatedAt.IsZero() {
		form.CreatedAt = time.Now()
	}
	form.UpdatedAt = time.Now()
	r.forms[form.ID] = form
	return nil
}

func (r *MockFormbuilderRepository) UpdateForm(ctx context.Context, form *domain.Form) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.forms[form.ID]; !exists {
		return appErrors.ErrNotFound
	}
	form.UpdatedAt = time.Now()
	r.forms[form.ID] = form
	return nil
}

func (r *MockFormbuilderRepository) DeleteForm(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.forms[id]; !exists {
		return appErrors.ErrNotFound
	}
	delete(r.forms, id)
	for fID, f := range r.fields {
		if f.FormID == id {
			delete(r.fields, fID)
		}
	}
	return nil
}

func (r *MockFormbuilderRepository) GetFieldByID(ctx context.Context, fieldID int64) (*domain.FormField, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	f, exists := r.fields[fieldID]
	if !exists {
		return nil, appErrors.ErrNotFound
	}
	copyField := *f
	return &copyField, nil
}

func (r *MockFormbuilderRepository) GetFieldsByFormID(ctx context.Context, formID int64) ([]*domain.FormField, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*domain.FormField
	for _, f := range r.fields {
		if f.FormID == formID {
			copyField := *f
			list = append(list, &copyField)
		}
	}
	return list, nil
}

func (r *MockFormbuilderRepository) AddField(ctx context.Context, field *domain.FormField) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.forms[field.FormID]; !exists {
		return appErrors.ErrNotFound
	}

	field.ID = r.nextFieldID
	r.nextFieldID++
	if field.CreatedAt.IsZero() {
		field.CreatedAt = time.Now()
	}
	field.UpdatedAt = time.Now()
	r.fields[field.ID] = field
	return nil
}

func (r *MockFormbuilderRepository) UpdateField(ctx context.Context, field *domain.FormField) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.fields[field.ID]; !exists {
		return appErrors.ErrNotFound
	}
	field.UpdatedAt = time.Now()
	r.fields[field.ID] = field
	return nil
}

func (r *MockFormbuilderRepository) DeleteField(ctx context.Context, fieldID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.fields[fieldID]; !exists {
		return appErrors.ErrNotFound
	}
	delete(r.fields, fieldID)
	return nil
}

func (r *MockFormbuilderRepository) CountFields(ctx context.Context, formID int64) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, f := range r.fields {
		if f.FormID == formID {
			count++
		}
	}
	return count, nil
}
