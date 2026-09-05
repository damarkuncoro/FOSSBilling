import React, { useState, useEffect } from 'react';
import { X, Plus, Trash2 } from 'lucide-react';
import type { FormField, FormFieldType } from '../../types/formBuilder';

interface FieldEditorDialogProps {
  isOpen: boolean;
  field: FormField | null;
  onClose: () => void;
  onSave: (field: FormField) => void;
}

export const FieldEditorDialog: React.FC<FieldEditorDialogProps> = ({
  isOpen,
  field,
  onClose,
  onSave,
}) => {
  const [name, setName] = useState('');
  const [label, setLabel] = useState('');
  const [type, setType] = useState<FormFieldType>('text');
  const [required, setRequired] = useState(false);
  const [placeholder, setPlaceholder] = useState('');
  const [description, setDescription] = useState('');
  const [options, setOptions] = useState<string[]>([]);
  const [newOption, setNewOption] = useState('');

  useEffect(() => {
    if (field) {
      setName(field.name);
      setLabel(field.label);
      setType(field.type);
      setRequired(field.required);
      setPlaceholder(field.placeholder || '');
      setDescription(field.description || '');
      setOptions(field.options || []);
    } else {
      setName('');
      setLabel('');
      setType('text');
      setRequired(false);
      setPlaceholder('');
      setDescription('');
      setOptions([]);
    }
  }, [field, isOpen]);

  if (!isOpen) return null;

  const handleAddOption = () => {
    if (newOption.trim()) {
      setOptions([...options, newOption.trim()]);
      setNewOption('');
    }
  };

  const handleRemoveOption = (index: number) => {
    setOptions(options.filter((_, i) => i !== index));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave({
      id: field?.id || `f_${Date.now()}`,
      name: name.toLowerCase().replace(/[^a-z0-9_]/g, '_'),
      label,
      type,
      required,
      placeholder,
      description,
      options: ['dropdown', 'radio', 'checkbox'].includes(type) ? options : undefined,
    });
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl animate-in fade-in zoom-in duration-200">
        <div className="flex items-center justify-between pb-4 border-b border-gray-100">
          <h3 className="text-gray-900 font-bold text-base">{field ? 'Edit Field' : 'Add Custom Field'}</h3>
          <button onClick={onClose} className="p-1 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100">
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Field Label</label>
            <input
              type="text"
              required
              placeholder="e.g. Server Hostname"
              value={label}
              onChange={(e) => setLabel(e.target.value)}
              className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Variable Name</label>
              <input
                type="text"
                required
                placeholder="hostname"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-3 py-2 font-mono border border-gray-200 rounded-lg text-xs outline-none focus:border-indigo-500"
              />
            </div>
            <div>
              <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Field Type</label>
              <select
                value={type}
                onChange={(e) => setType(e.target.value as FormFieldType)}
                className="w-full px-3 py-2 border border-gray-200 rounded-lg text-xs outline-none focus:border-indigo-500 bg-white"
              >
                <option value="text">Text Input</option>
                <option value="textarea">Textarea</option>
                <option value="number">Number</option>
                <option value="dropdown">Dropdown Select</option>
                <option value="radio">Radio Options</option>
                <option value="checkbox">Checkbox Toggle</option>
              </select>
            </div>
          </div>

          {['dropdown', 'radio'].includes(type) && (
            <div className="p-3 bg-gray-50 rounded-xl border border-gray-100 space-y-2">
              <label className="block text-xs font-semibold text-gray-700 uppercase">Choice Options</label>
              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder="Add option..."
                  value={newOption}
                  onChange={(e) => setNewOption(e.target.value)}
                  className="flex-1 px-3 py-1.5 border border-gray-200 rounded-lg text-xs bg-white outline-none"
                />
                <button
                  type="button"
                  onClick={handleAddOption}
                  className="px-3 py-1.5 bg-indigo-600 text-white rounded-lg text-xs font-medium flex items-center gap-1"
                >
                  <Plus className="w-3.5 h-3.5" /> Add
                </button>
              </div>
              <div className="flex flex-wrap gap-1.5 mt-2">
                {options.map((opt, i) => (
                  <span key={i} className="inline-flex items-center gap-1 px-2.5 py-1 bg-white border border-gray-200 rounded-lg text-xs text-gray-700">
                    {opt}
                    <button type="button" onClick={() => handleRemoveOption(i)} className="text-gray-400 hover:text-rose-600">
                      <Trash2 className="w-3 h-3" />
                    </button>
                  </span>
                ))}
              </div>
            </div>
          )}

          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Help / Description</label>
            <input
              type="text"
              placeholder="e.g. Instructions for the client"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-xs outline-none focus:border-indigo-500"
            />
          </div>

          <div className="flex items-center gap-2 pt-2">
            <input
              type="checkbox"
              id="req_toggle"
              checked={required}
              onChange={(e) => setRequired(e.target.checked)}
              className="w-4 h-4 text-indigo-600 rounded border-gray-300 focus:ring-indigo-500"
            />
            <label htmlFor="req_toggle" className="text-xs font-medium text-gray-700">
              Required field (client must fill this out at checkout)
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-4 border-t border-gray-100">
            <button type="button" onClick={onClose} className="px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 rounded-lg">
              Cancel
            </button>
            <button type="submit" className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-sm">
              Save Field
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
