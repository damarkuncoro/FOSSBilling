import { useState, useEffect, useCallback } from 'react';
import { adminFormBuilderService } from '@/services/admin_form_builder.service';
import type { CustomForm, FormField } from '../types/formBuilder';

export function useFormBuilder() {
  const [forms, setForms] = useState<CustomForm[]>([]);
  const [selectedForm, setSelectedForm] = useState<CustomForm | null>(null);
  const [loading, setLoading] = useState(true);
  const [isNewFormModal, setIsNewFormModal] = useState(false);
  const [isFieldModal, setIsFieldModal] = useState(false);
  const [editingField, setEditingField] = useState<FormField | null>(null);

  const fetchForms = useCallback(async () => {
    setLoading(true);
    try {
      const data = await adminFormBuilderService.listForms();
      setForms(data || []);
      if (data && data.length > 0 && !selectedForm) {
        setSelectedForm(data[0]);
      }
    } catch (err) {
      console.error('Failed to fetch forms:', err);
    } finally {
      setLoading(false);
    }
  }, [selectedForm]);

  useEffect(() => {
    fetchForms();
  }, [fetchForms]);

  const createForm = async (name: string, type: string) => {
    try {
      const newForm = await adminFormBuilderService.createForm(name, type);
      setForms((prev) => [newForm, ...prev]);
      setSelectedForm(newForm);
      setIsNewFormModal(false);
    } catch (err) {
      console.error('Failed to create form:', err);
    }
  };

  const deleteForm = async (id: number) => {
    if (!confirm('Are you sure you want to remove this form template?')) return;
    try {
      await adminFormBuilderService.deleteForm(id);
      setForms((prev) => prev.filter((f) => f.id !== id));
      if (selectedForm?.id === id) {
        setSelectedForm(null);
      }
    } catch (err) {
      console.error('Failed to delete form:', err);
    }
  };

  const saveField = async (field: Partial<FormField>) => {
    if (!selectedForm) return;
    try {
      if (field.id) {
        await adminFormBuilderService.updateField(field.id, field);
      } else {
        await adminFormBuilderService.addField(selectedForm.id, field);
      }
      // Refresh selected form data
      const updatedForms = await adminFormBuilderService.listForms();
      setForms(updatedForms);
      const found = updatedForms.find(f => f.id === selectedForm.id);
      if (found) setSelectedForm(found);

      setIsFieldModal(false);
      setEditingField(null);
    } catch (err) {
      console.error('Failed to save field:', err);
    }
  };

  const removeField = async (fieldId: number) => {
    if (!confirm('Remove this field?')) return;
    try {
      await adminFormBuilderService.deleteField(fieldId);
      if (selectedForm) {
        const updatedFields = selectedForm.fields.filter((f) => f.id !== fieldId);
        setSelectedForm({ ...selectedForm, fields: updatedFields });
        setForms((prev) => prev.map(f => f.id === selectedForm.id ? { ...f, fields: updatedFields } : f));
      }
    } catch (err) {
      console.error('Failed to delete field:', err);
    }
  };

  return {
    forms,
    selectedForm,
    setSelectedForm,
    loading,
    isNewFormModal,
    setIsNewFormModal,
    isFieldModal,
    setIsFieldModal,
    editingField,
    setEditingField,
    createForm,
    deleteForm,
    saveField,
    removeField,
  };
}
