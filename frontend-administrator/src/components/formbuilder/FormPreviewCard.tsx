import React from 'react';
import { Eye, Edit3, Trash2 } from 'lucide-react';
import type { FormField } from '../../types/formBuilder';

interface FormPreviewCardProps {
  fields: FormField[];
  onEditField: (field: FormField) => void;
  onRemoveField: (id: string) => void;
}

export const FormPreviewCard: React.FC<FormPreviewCardProps> = ({
  fields,
  onEditField,
  onRemoveField,
}) => {
  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-6">
      <div className="flex items-center gap-2 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4 pb-3 border-b border-gray-100">
        <Eye className="w-4 h-4 text-indigo-600" /> Client Checkout Form Preview
      </div>

      {fields.length === 0 ? (
        <div className="text-center py-8 text-gray-400 text-sm">
          No custom fields added yet. Click &quot;Add Custom Field&quot; to build this form.
        </div>
      ) : (
        <div className="space-y-4">
          {fields.map((f) => (
            <div key={f.id} className="group relative p-3.5 bg-gray-50/70 hover:bg-gray-50 border border-gray-200/80 rounded-xl transition-all">
              <div className="flex items-center justify-between mb-1.5">
                <label className="text-xs font-semibold text-gray-800 flex items-center gap-1.5">
                  {f.label}
                  {f.required && <span className="text-rose-500 font-bold">*</span>}
                  <span className="text-[10px] font-mono font-normal text-gray-400 px-1.5 py-0.5 bg-white border border-gray-200 rounded">
                    {f.name} ({f.type})
                  </span>
                </label>
                <div className="flex items-center gap-1 opacity-80 group-hover:opacity-100 transition-opacity">
                  <button onClick={() => onEditField(f)} className="p-1 text-gray-400 hover:text-indigo-600 rounded">
                    <Edit3 className="w-3.5 h-3.5" />
                  </button>
                  <button onClick={() => onRemoveField(f.id)} className="p-1 text-gray-400 hover:text-rose-600 rounded">
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>

              {f.description && <p className="text-xs text-gray-500 mb-2">{f.description}</p>}

              {f.type === 'text' && (
                <input type="text" disabled placeholder={f.placeholder || f.label} className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-xs" />
              )}
              {f.type === 'textarea' && (
                <textarea disabled placeholder={f.placeholder || f.label} rows={2} className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-xs" />
              )}
              {f.type === 'dropdown' && (
                <select disabled className="w-full px-3 py-2 bg-white border border-gray-200 rounded-lg text-xs">
                  {f.options?.map((opt, i) => <option key={i}>{opt}</option>)}
                </select>
              )}
              {f.type === 'checkbox' && (
                <div className="flex items-center gap-2">
                  <input type="checkbox" disabled className="rounded border-gray-300" />
                  <span className="text-xs text-gray-600">Check to enable</span>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
