import { useState } from 'react';
import type { CustomForm, FormField } from '../types/formBuilder';

const initialForms: CustomForm[] = [
  {
    id: 1,
    name: 'VPS & Dedicated Server Setup',
    description: 'Requires hostname, root password, and operating system distribution.',
    product_ids: [1, 2],
    fields: [
      {
        id: 'f1',
        name: 'hostname',
        label: 'Server Hostname',
        type: 'text',
        required: true,
        placeholder: 'server1.mydomain.com',
        description: 'Fully qualified domain name for your server.',
      },
      {
        id: 'f2',
        name: 'root_password',
        label: 'Root Password',
        type: 'text',
        required: true,
        placeholder: 'Choose strong root password',
      },
      {
        id: 'f3',
        name: 'os_distribution',
        label: 'Operating System',
        type: 'dropdown',
        required: true,
        options: ['Ubuntu 24.04 LTS', 'Debian 12 Bookworm', 'AlmaLinux 9', 'Rocky Linux 9', 'Windows Server 2022'],
        default_value: 'Ubuntu 24.04 LTS',
      },
      {
        id: 'f4',
        name: 'backup_addon',
        label: 'Enable Automated Daily Backups',
        type: 'checkbox',
        required: false,
      },
    ],
    created_at: '2026-02-10T08:00:00Z',
    updated_at: '2026-08-15T10:00:00Z',
  },
  {
    id: 2,
    name: 'Managed Game Server Config',
    description: 'Game server slot options and tickrate configuration.',
    product_ids: [3],
    fields: [
      {
        id: 'g1',
        name: 'server_name',
        label: 'Public Server Name',
        type: 'text',
        required: true,
        placeholder: 'My Epic Minecraft Realm',
      },
      {
        id: 'g2',
        name: 'server_region',
        label: 'Location / Region',
        type: 'dropdown',
        required: true,
        options: ['Singapore (SG)', 'Tokyo (JP)', 'Frankfurt (DE)', 'US East (VA)'],
        default_value: 'Singapore (SG)',
      },
    ],
    created_at: '2026-04-01T12:00:00Z',
    updated_at: '2026-07-20T16:00:00Z',
  },
];

export function useFormBuilder() {
  const [forms, setForms] = useState<CustomForm[]>(initialForms);
  const [selectedForm, setSelectedForm] = useState<CustomForm | null>(forms[0]);
  const [isNewFormModal, setIsNewFormModal] = useState(false);
  const [isFieldModal, setIsFieldModal] = useState(false);
  const [editingField, setEditingField] = useState<FormField | null>(null);

  const createForm = (name: string, description: string) => {
    const newForm: CustomForm = {
      id: Date.now(),
      name,
      description,
      product_ids: [],
      fields: [],
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    setForms((prev) => [...prev, newForm]);
    setSelectedForm(newForm);
    setIsNewFormModal(false);
  };

  const deleteForm = (id: number) => {
    setForms((prev) => prev.filter((f) => f.id !== id));
    if (selectedForm?.id === id) {
      setSelectedForm(forms.find((f) => f.id !== id) || null);
    }
  };

  const saveField = (field: FormField) => {
    if (!selectedForm) return;
    const exists = selectedForm.fields.some((f) => f.id === field.id);
    const updatedFields = exists
      ? selectedForm.fields.map((f) => (f.id === field.id ? field : f))
      : [...selectedForm.fields, { ...field, id: field.id || `f_${Date.now()}` }];

    const updatedForm = { ...selectedForm, fields: updatedFields, updated_at: new Date().toISOString() };
    setSelectedForm(updatedForm);
    setForms((prev) => prev.map((f) => (f.id === updatedForm.id ? updatedForm : f)));
    setIsFieldModal(false);
    setEditingField(null);
  };

  const removeField = (fieldId: string) => {
    if (!selectedForm) return;
    const updatedFields = selectedForm.fields.filter((f) => f.id !== fieldId);
    const updatedForm = { ...selectedForm, fields: updatedFields, updated_at: new Date().toISOString() };
    setSelectedForm(updatedForm);
    setForms((prev) => prev.map((f) => (f.id === updatedForm.id ? updatedForm : f)));
  };

  return {
    forms,
    selectedForm,
    setSelectedForm,
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
