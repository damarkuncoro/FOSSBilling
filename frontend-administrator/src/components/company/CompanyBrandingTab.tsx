import React from 'react';
import { CompanySettings } from '@/lib/api';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

interface CompanyBrandingTabProps {
  form: CompanySettings;
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
}

export const CompanyBrandingTab: React.FC<CompanyBrandingTabProps> = ({ form, onChange }) => {
  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-1.5">
          <label className="text-xs font-semibold">Light Theme Logo URL</label>
          <Input
            name="logo_url"
            placeholder="https://example.com/assets/logo.png"
            value={form.logo_url || ''}
            onChange={onChange}
          />
          {form.logo_url && (
            <div className="mt-2 p-2 bg-slate-100 rounded-lg border w-fit">
              <img src={form.logo_url} alt="Logo Preview" className="h-7 object-contain" />
            </div>
          )}
        </div>

        <div className="space-y-1.5">
          <label className="text-xs font-semibold">Dark Theme Logo URL</label>
          <Input
            name="logo_dark_url"
            placeholder="https://example.com/assets/logo-white.png"
            value={form.logo_dark_url || ''}
            onChange={onChange}
          />
          {form.logo_dark_url && (
            <div className="mt-2 p-2 bg-slate-900 rounded-lg border w-fit">
              <img src={form.logo_dark_url} alt="Dark Logo Preview" className="h-7 object-contain" />
            </div>
          )}
        </div>
      </div>

      <div className="space-y-1.5">
        <label className="text-xs font-semibold">Favicon Icon URL</label>
        <Input
          name="favicon_url"
          placeholder="https://example.com/favicon.ico"
          value={form.favicon_url || ''}
          onChange={onChange}
        />
      </div>

      <div className="space-y-1.5">
        <label className="text-xs font-semibold">Default Email Signature (Plaintext / Markdown)</label>
        <Textarea
          name="email_signature"
          rows={4}
          placeholder="Best regards,&#10;Customer Support Team&#10;My Company Ltd."
          value={form.email_signature || ''}
          onChange={onChange}
        />
      </div>
    </div>
  );
};
