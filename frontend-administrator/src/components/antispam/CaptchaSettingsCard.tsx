import React from 'react';
import { ShieldCheck, CheckCircle } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { SecuritySettings } from '@/types/api';

interface CaptchaSettingsCardProps {
  settings: SecuritySettings;
  onChange: (settings: SecuritySettings) => void;
  onSubmit: (e: React.FormEvent) => void;
  saving: boolean;
  message: string | null;
}

export const CaptchaSettingsCard: React.FC<CaptchaSettingsCardProps> = ({
  settings,
  onChange,
  onSubmit,
  saving,
  message,
}) => {
  return (
    <Card className="border-border/60 shadow-sm">
      <CardHeader>
        <div className="flex items-center gap-2">
          <ShieldCheck className="h-5 w-5 text-primary" />
          <CardTitle className="text-base font-semibold">Captcha & Bot Defense</CardTitle>
        </div>
        <CardDescription>
          Protect registration and login forms against automated spam attacks
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
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Captcha Protection Provider</label>
            <select
              value={settings.recaptcha_provider}
              onChange={(e: any) => onChange({ ...settings, recaptcha_provider: e.target.value })}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus:outline-none focus:ring-1 focus:ring-ring"
            >
              <option value="cloudflare_turnstile">Cloudflare Turnstile (Recommended - Zero Friction)</option>
              <option value="google_recaptcha">Google reCAPTCHA v2 / v3 Invisible</option>
            </select>
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Public Site Key</label>
            <Input
              required
              value={settings.site_key}
              onChange={(e) => onChange({ ...settings, site_key: e.target.value })}
              placeholder="0x4AAAAAA..."
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Max Login Attempts</label>
              <Input
                type="number"
                value={settings.max_login_attempts}
                onChange={(e) => onChange({ ...settings, max_login_attempts: parseInt(e.target.value) || 5 })}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Lockout Duration (Mins)</label>
              <Input
                type="number"
                value={settings.lockout_time_minutes}
                onChange={(e) => onChange({ ...settings, lockout_time_minutes: parseInt(e.target.value) || 15 })}
              />
            </div>
          </div>

          <Button type="submit" disabled={saving} className="font-semibold shadow-sm">
            {saving ? 'Saving...' : 'Save Bot Protection'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
};
