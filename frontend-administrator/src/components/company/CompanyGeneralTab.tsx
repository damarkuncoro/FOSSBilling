import React from 'react';
import { CompanySettings } from '@/lib/api';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

interface CompanyGeneralTabProps {
  form: CompanySettings;
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
}

export const CompanyGeneralTab: React.FC<CompanyGeneralTabProps> = ({ form, onChange }) => {
  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-1.5">
          <label className="text-xs font-semibold">Company / Business Name</label>
          <Input
            required
            name="name"
            placeholder="e.g. My Hosting Company Ltd."
            value={form.name}
            onChange={onChange}
          />
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-semibold">Contact & Billing Email</label>
          <Input
            type="email"
            required
            name="email"
            placeholder="billing@example.com"
            value={form.email}
            onChange={onChange}
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-1.5">
          <label className="text-xs font-semibold">Official Phone Number</label>
          <Input
            name="phone"
            placeholder="+1 (555) 019-2834"
            value={form.phone}
            onChange={onChange}
          />
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-semibold">Tax / VAT Registration Number</label>
          <Input
            name="vat_number"
            placeholder="e.g. EU123456789 or ID-NPWP"
            value={form.vat_number || ''}
            onChange={onChange}
          />
        </div>
      </div>

      <div className="space-y-1.5">
        <label className="text-xs font-semibold">Street Address Line 1</label>
        <Input
          required
          name="address_1"
          placeholder="123 Business Avenue, Suite 400"
          value={form.address_1}
          onChange={onChange}
        />
      </div>

      <div className="space-y-1.5">
        <label className="text-xs font-semibold">Street Address Line 2 (Optional)</label>
        <Input
          name="address_2"
          placeholder="Building B, Floor 2"
          value={form.address_2 || ''}
          onChange={onChange}
        />
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div className="space-y-1.5">
          <label className="text-xs font-semibold">City</label>
          <Input
            required
            name="city"
            placeholder="San Francisco"
            value={form.city}
            onChange={onChange}
          />
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-semibold">State / Province</label>
          <Input
            required
            name="state"
            placeholder="CA"
            value={form.state}
            onChange={onChange}
          />
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-semibold">Postal Code</label>
          <Input
            required
            name="postcode"
            placeholder="94105"
            value={form.postcode}
            onChange={onChange}
          />
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-semibold">Country</label>
          <Input
            required
            name="country"
            placeholder="United States"
            value={form.country}
            onChange={onChange}
          />
        </div>
      </div>

      <div className="space-y-1.5">
        <label className="text-xs font-semibold">Terms & Conditions URL</label>
        <Input
          name="terms_url"
          placeholder="https://example.com/terms"
          value={form.terms_url || ''}
          onChange={onChange}
        />
      </div>
    </div>
  );
};
