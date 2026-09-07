package formbuilder_test

import (
	"context"
	"errors"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/formbuilder"
)

func setupFormbuilderService() (*formbuilder.FormbuilderService, *memory.MockFormbuilderRepository) {
	repo := memory.NewMockFormbuilderRepository()
	svc := formbuilder.NewFormbuilderService(repo)
	return svc, repo
}

func TestFormbuilderService_FormLifecycle(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupFormbuilderService()

	// 1. Create Form
	createDTO := formbuilder.CreateFormDTO{
		Name: "VPS Custom Options",
		Style: domain.FormStyle{
			Type:      "horizontal",
			ShowTitle: true,
		},
	}
	form, err := svc.CreateForm(ctx, createDTO)
	if err != nil {
		t.Fatalf("CreateForm failed: %v", err)
	}
	if form.ID == 0 || form.Name != "VPS Custom Options" {
		t.Errorf("Unexpected form result: %+v", form)
	}

	// 2. Get Form
	fetched, err := svc.GetForm(ctx, form.ID)
	if err != nil {
		t.Fatalf("GetForm failed: %v", err)
	}
	if fetched.Name != form.Name {
		t.Errorf("Fetched name = %s; want %s", fetched.Name, form.Name)
	}

	// 3. List Forms
	forms, total, err := svc.ListForms(ctx, 10, 0)
	if err != nil || total != 1 || len(forms) != 1 {
		t.Fatalf("ListForms failed: total=%d, len=%d, err=%v", total, len(forms), err)
	}

	// 4. Update Form
	updated, err := svc.UpdateForm(ctx, form.ID, formbuilder.UpdateFormDTO{
		Name: "VPS Configurator Pro",
		Style: domain.FormStyle{
			Type: "vertical",
		},
	})
	if err != nil {
		t.Fatalf("UpdateForm failed: %v", err)
	}
	if updated.Name != "VPS Configurator Pro" || updated.Style.Type != "vertical" {
		t.Errorf("Unexpected updated form: %+v", updated)
	}

	// 5. Delete Form
	err = svc.DeleteForm(ctx, form.ID)
	if err != nil {
		t.Fatalf("DeleteForm failed: %v", err)
	}
	_, err = svc.GetForm(ctx, form.ID)
	if err == nil {
		t.Error("Expected error retrieving deleted form, got nil")
	}
}

func TestFormbuilderService_FieldsLifecycle(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupFormbuilderService()

	form, _ := svc.CreateForm(ctx, formbuilder.CreateFormDTO{
		Name: "Dedicated Server Addons",
	})

	// 1. Add Text Field with slugification
	field1, err := svc.AddField(ctx, form.ID, formbuilder.FormFieldDTO{
		Name:     "Root Password",
		Label:    "Custom Root Password",
		Type:     "text",
		Required: true,
	})
	if err != nil {
		t.Fatalf("AddField text failed: %v", err)
	}
	if field1.Name != "root_password" {
		t.Errorf("Field name slugification failed: got %s; want root_password", field1.Name)
	}

	// 2. Add Select Field
	field2, err := svc.AddField(ctx, form.ID, formbuilder.FormFieldDTO{
		Name:  "operating_system",
		Label: "Select OS",
		Type:  "select",
		Options: map[string]interface{}{
			"Ubuntu 22.04 LTS": "ubuntu_22",
			"Debian 12":        "debian_12",
			"AlmaLinux 9":      "alma_9",
		},
		Required: true,
	})
	if err != nil {
		t.Fatalf("AddField select failed: %v", err)
	}
	if len(field2.Options) != 3 {
		t.Errorf("Options count = %d; want 3", len(field2.Options))
	}

	// 3. Duplicate Field Name in Same Form should fail
	_, err = svc.AddField(ctx, form.ID, formbuilder.FormFieldDTO{
		Name: "root_password",
		Type: "text",
	})
	if !errors.Is(err, formbuilder.ErrDuplicateFieldName) {
		t.Errorf("Expected ErrDuplicateFieldName, got: %v", err)
	}

	// 4. Numeric starting field name should fail
	_, err = svc.AddField(ctx, form.ID, formbuilder.FormFieldDTO{
		Name: "123_invalid_field",
		Type: "text",
	})
	if !errors.Is(err, formbuilder.ErrFieldNameNumeric) {
		t.Errorf("Expected ErrFieldNameNumeric, got: %v", err)
	}

	// 5. Invalid field type should fail
	_, err = svc.AddField(ctx, form.ID, formbuilder.FormFieldDTO{
		Name: "valid_name",
		Type: "unsupported_type",
	})
	if !errors.Is(err, formbuilder.ErrInvalidFieldType) {
		t.Errorf("Expected ErrInvalidFieldType, got: %v", err)
	}

	// 6. Update Field
	updatedField, err := svc.UpdateField(ctx, field1.ID, formbuilder.FormFieldDTO{
		Label:    "Initial Root Password (min 12 chars)",
		Required: true,
	})
	if err != nil {
		t.Fatalf("UpdateField failed: %v", err)
	}
	if updatedField.Label != "Initial Root Password (min 12 chars)" {
		t.Errorf("Updated label = %s; want 'Initial Root Password (min 12 chars)'", updatedField.Label)
	}

	// 7. Check Form has attached fields
	formWithFields, err := svc.GetForm(ctx, form.ID)
	if err != nil || len(formWithFields.Fields) != 2 {
		t.Fatalf("Expected form to have 2 fields, got %d", len(formWithFields.Fields))
	}

	// 8. Delete Field
	err = svc.DeleteField(ctx, field1.ID)
	if err != nil {
		t.Fatalf("DeleteField failed: %v", err)
	}
	formAfterDelete, _ := svc.GetForm(ctx, form.ID)
	if len(formAfterDelete.Fields) != 1 {
		t.Errorf("Expected 1 field remaining, got %d", len(formAfterDelete.Fields))
	}
}

func TestFormbuilderService_ValidateFormSubmission(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupFormbuilderService()

	form, _ := svc.CreateForm(ctx, formbuilder.CreateFormDTO{Name: "Domain Setup"})
	_, _ = svc.AddField(ctx, form.ID, formbuilder.FormFieldDTO{
		Name:     "target_url",
		Label:    "Target URL",
		Type:     "url",
		Required: true,
	})
	_, _ = svc.AddField(ctx, form.ID, formbuilder.FormFieldDTO{
		Name:     "enable_ssl",
		Label:    "Enable SSL",
		Type:     "checkbox",
		Required: false,
	})

	// 1. Valid Submission
	validData := map[string]interface{}{
		"target_url": "https://mywebsite.com",
		"enable_ssl": "1",
	}
	cleaned, err := svc.ValidateFormSubmission(ctx, form.ID, validData)
	if err != nil {
		t.Fatalf("Expected valid submission to pass, got: %v", err)
	}
	if cleaned["enable_ssl"] != true {
		t.Errorf("Expected enable_ssl to be boolean true, got: %v", cleaned["enable_ssl"])
	}

	// 2. Missing Required Field
	missingData := map[string]interface{}{
		"enable_ssl": true,
	}
	_, err = svc.ValidateFormSubmission(ctx, form.ID, missingData)
	if !errors.Is(err, formbuilder.ErrRequiredFieldMissing) {
		t.Errorf("Expected ErrRequiredFieldMissing, got: %v", err)
	}

	// 3. Invalid URL Field
	invalidURLData := map[string]interface{}{
		"target_url": "not-a-valid-url",
	}
	_, err = svc.ValidateFormSubmission(ctx, form.ID, invalidURLData)
	if !errors.Is(err, formbuilder.ErrInvalidURLField) {
		t.Errorf("Expected ErrInvalidURLField, got: %v", err)
	}
}
