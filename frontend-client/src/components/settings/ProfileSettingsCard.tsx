import React from 'react';
import { User, CheckCircle } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { ClientProfile } from '@/types/api';

interface ProfileSettingsCardProps {
  user: ClientProfile | null;
  form: { first_name: string; last_name: string; company: string; country: string };
  onChange: (form: { first_name: string; last_name: string; company: string; country: string }) => void;
  onSubmit: (e: React.FormEvent) => void;
  saving: boolean;
  message: string | null;
}

export const ProfileSettingsCard: React.FC<ProfileSettingsCardProps> = ({
  user,
  form,
  onChange,
  onSubmit,
  saving,
  message,
}) => {
  return (
    <Card className="border-border/60 shadow-sm">
      <CardHeader>
        <div className="flex items-center gap-2">
          <User className="h-5 w-5 text-primary" />
          <CardTitle className="text-base font-semibold">
            Personal & Company Information
          </CardTitle>
        </div>
        <CardDescription>
          Update your contact details for automated billing and tax compliance
        </CardDescription>
      </CardHeader>
      <CardContent>
        {message && (
          <div className="mb-4 p-3 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 text-xs font-semibold flex items-center gap-2">
            <CheckCircle className="h-4 w-4" />
            <span>{message}</span>
          </div>
        )}

        <form onSubmit={onSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">First Name</label>
              <Input
                required
                value={form.first_name}
                onChange={(e) => onChange({ ...form, first_name: e.target.value })}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Last Name</label>
              <Input
                required
                value={form.last_name}
                onChange={(e) => onChange({ ...form, last_name: e.target.value })}
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Email Address (Read-only)</label>
            <Input disabled value={user?.email || ''} className="bg-muted" />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Company Name (Optional)</label>
            <Input
              placeholder="e.g. PT Nusantara Solusi Digital"
              value={form.company}
              onChange={(e) => onChange({ ...form, company: e.target.value })}
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Country Code</label>
            <Input
              value={form.country}
              onChange={(e) => onChange({ ...form, country: e.target.value })}
            />
          </div>

          <Button type="submit" disabled={saving} className="font-semibold shadow-sm">
            {saving ? 'Saving...' : 'Save Profile Changes'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
};
