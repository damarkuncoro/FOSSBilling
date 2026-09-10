import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminFormBuilderService } from '@/services/admin_form_builder.service';
import type { FormField } from '../types/formBuilder';

export function useFormBuilder() {
  const queryClient = useQueryClient();
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [isNewFormModal, setIsNewFormModal] = useState(false);
  const [isFieldModal, setIsFieldModal] = useState(false);
  const [editingField, setEditingField] = useState<FormField | null>(null);

  const { data: forms = [], isLoading: loading } = useQuery({
    queryKey: ['admin', 'forms'],
    queryFn: () => adminFormBuilderService.listForms(),
  });

  const createMutation = useMutation({
    mutationFn: ({ name, type }: { name: string; type: string }) => adminFormBuilderService.createForm(name, type),
    onSuccess: (f) => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'forms'] });
      setSelectedId(f.id);
      setIsNewFormModal(false);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminFormBuilderService.deleteForm(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'forms'] });
      setSelectedId(null);
    },
  });

  const fieldMutation = useMutation({
    mutationFn: (f: Partial<FormField>) => f.id ? adminFormBuilderService.updateField(f.id, f) : adminFormBuilderService.addField(selectedId!, f),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'forms'] });
      setIsFieldModal(false);
      setEditingField(null);
    },
  });

  const removeFieldMutation = useMutation({
    mutationFn: (id: number) => adminFormBuilderService.deleteField(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'forms'] }),
  });

  const selectedForm = forms.find(f => f.id === selectedId) || (forms.length > 0 ? forms[0] : null);

  return {
    forms,
    selectedForm,
    setSelectedForm: (f: any) => setSelectedId(f ? f.id : null),
    loading,
    isNewFormModal,
    setIsNewFormModal,
    isFieldModal,
    setIsFieldModal,
    editingField,
    setEditingField,
    createForm: (name: string, type: string) => createMutation.mutate({ name, type }),
    deleteForm: (id: number) => { if (confirm('Delete form?')) deleteMutation.mutate(id); },
    saveField: fieldMutation.mutate,
    removeField: (id: number) => { if (confirm('Remove field?')) removeFieldMutation.mutate(id); },
  };
}
