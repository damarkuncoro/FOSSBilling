import React, { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { Input } from '@/components/ui/button'; // Wait, Input is in separate file
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';

interface ProductConfigFormProps {
  formId: number;
  onConfigChange: (config: Record<string, any>) => void;
  initialConfig?: Record<string, any>;
}

export const ProductConfigForm: React.FC<ProductConfigFormProps> = ({
  formId,
  onConfigChange,
  initialConfig = {},
}) => {
  const [form, setForm] = useState<any>(null);
  const [values, setValues] = useState<Record<string, any>>(initialConfig);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getForm(formId)
      .then(setForm)
      .finally(() => setLoading(false));
  }, [formId]);

  const handleChange = (name: string, value: any) => {
    const newValues = { ...values, [name]: value };
    setValues(newValues);
    onConfigChange(newValues);
  };

  if (loading) return <div className="text-[10px] text-muted-foreground animate-pulse">Loading configuration fields...</div>;
  if (!form || !form.fields || form.fields.length === 0) return null;

  return (
    <div className="mt-3 p-4 bg-muted/30 rounded-xl border border-border/40 space-y-4">
      <div className="flex items-center gap-2 mb-1">
        <Badge variant="secondary" className="text-[9px] uppercase font-bold">Configuration Required</Badge>
        <span className="text-xs font-semibold text-foreground/80">{form.name}</span>
      </div>

      <div className={`grid gap-4 ${form.style?.type === 'horizontal' ? 'grid-cols-1 md:grid-cols-2' : 'grid-cols-1'}`}>
        {form.fields.map((field: any) => (
          <div key={field.id} className="space-y-1.5">
            {!field.hide_label && (
              <label className="text-[11px] font-bold text-muted-foreground uppercase tracking-tight">
                {field.label} {field.required && <span className="text-destructive">*</span>}
              </label>
            )}

            {field.type === 'text' && (
              <input
                type="text"
                className="w-full h-9 rounded-lg border border-input bg-background px-3 py-1 text-sm shadow-sm focus:ring-1 focus:ring-primary outline-none"
                placeholder={field.description}
                value={values[field.name] || ''}
                onChange={(e) => handleChange(field.name, e.target.value)}
                required={field.required}
              />
            )}

            {field.type === 'textarea' && (
              <Textarea
                className="min-h-[80px]"
                placeholder={field.description}
                value={values[field.name] || ''}
                onChange={(e) => handleChange(field.name, e.target.value)}
                required={field.required}
              />
            )}

            {field.type === 'select' && (
              <select
                className="w-full h-9 rounded-lg border border-input bg-background px-3 py-1 text-sm shadow-sm focus:ring-1 focus:ring-primary outline-none"
                value={values[field.name] || ''}
                onChange={(e) => handleChange(field.name, e.target.value)}
                required={field.required}
              >
                <option value="">Select an option...</option>
                {field.options && Object.entries(field.options).map(([label, val]: [string, any]) => (
                  <option key={label} value={val}>{label}</option>
                ))}
              </select>
            )}

            {field.type === 'checkbox' && (
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id={`field-${field.id}`}
                  checked={!!values[field.name]}
                  onChange={(e) => handleChange(field.name, e.target.checked)}
                />
                <label htmlFor={`field-${field.id}`} className="text-xs text-muted-foreground cursor-pointer">
                  {field.description || 'Enable this option'}
                </label>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};
