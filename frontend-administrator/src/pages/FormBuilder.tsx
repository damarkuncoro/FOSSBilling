import React from 'react';
import { Layers, Plus, Trash2, FileText } from 'lucide-react';
import { useFormBuilder } from '../hooks/useFormBuilder';
import { FormPreviewCard } from '../components/formbuilder/FormPreviewCard';
import { FieldEditorDialog } from '../components/formbuilder/FieldEditorDialog';
import { CreateFormDialog } from '../components/formbuilder/CreateFormDialog';

export const FormBuilder: React.FC = () => {
  const {
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
  } = useFormBuilder();

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
            <Layers className="w-7 h-7 text-indigo-600" /> Custom Form Builder
          </h1>
          <p className="text-sm text-gray-500 mt-1">
            Build custom input fields (Hostname, OS, Dropdowns) attached to product checkouts.
          </p>
        </div>
        <button
          onClick={() => setIsNewFormModal(true)}
          className="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-sm rounded-xl shadow-sm transition-all"
        >
          <Plus className="w-4 h-4" /> New Form Template
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left column: List of forms */}
        <div className="space-y-3">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider px-1">Form Templates</h3>
          {forms.map((form) => (
            <div
              key={form.id}
              onClick={() => setSelectedForm(form)}
              className={`p-4 rounded-xl border cursor-pointer transition-all ${
                selectedForm?.id === form.id
                  ? 'bg-indigo-50/50 border-indigo-200 ring-2 ring-indigo-500/20 shadow-sm'
                  : 'bg-white border-gray-200/80 hover:border-gray-300'
              }`}
            >
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2.5">
                  <FileText className={`w-4 h-4 ${selectedForm?.id === form.id ? 'text-indigo-600' : 'text-gray-400'}`} />
                  <span className="font-semibold text-sm text-gray-900">{form.name}</span>
                </div>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    deleteForm(form.id);
                  }}
                  className="p-1 text-gray-400 hover:text-rose-600 rounded"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                </button>
              </div>
              <p className="text-xs text-gray-500 mt-2 line-clamp-2">{form.description}</p>
              <div className="mt-3 flex items-center gap-2">
                <span className="px-2 py-0.5 rounded text-[11px] font-medium bg-gray-100 text-gray-600">
                  {form.fields.length} Custom Fields
                </span>
              </div>
            </div>
          ))}
        </div>

        {/* Right column: Form Fields Editor and Preview */}
        <div className="lg:col-span-2 space-y-4">
          {selectedForm ? (
            <>
              <div className="flex items-center justify-between bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
                <div>
                  <h2 className="text-base font-bold text-gray-900">{selectedForm.name}</h2>
                  <p className="text-xs text-gray-500">{selectedForm.description}</p>
                </div>
                <button
                  onClick={() => {
                    setEditingField(null);
                    setIsFieldModal(true);
                  }}
                  className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-xs rounded-lg"
                >
                  <Plus className="w-3.5 h-3.5" /> Add Field
                </button>
              </div>

              <FormPreviewCard
                fields={selectedForm.fields}
                onEditField={(field) => {
                  setEditingField(field);
                  setIsFieldModal(true);
                }}
                onRemoveField={removeField}
              />
            </>
          ) : (
            <div className="bg-white rounded-xl border border-gray-100 p-12 text-center text-gray-400">
              Select or create a form template to start adding fields.
            </div>
          )}
        </div>
      </div>

      <CreateFormDialog
        isOpen={isNewFormModal}
        onClose={() => setIsNewFormModal(false)}
        onSubmit={createForm}
      />

      <FieldEditorDialog
        isOpen={isFieldModal}
        field={editingField}
        onClose={() => {
          setIsFieldModal(false);
          setEditingField(null);
        }}
        onSave={saveField}
      />
    </div>
  );
};
