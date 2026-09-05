import React from 'react';
import { Bot, ShieldAlert } from 'lucide-react';
import { SecuritySettings } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

interface SecurityPolicyTabProps {
  securitySettings: SecuritySettings;
  setSecuritySettings: React.Dispatch<React.SetStateAction<SecuritySettings>>;
  blacklistText: string;
  setBlacklistText: (val: string) => void;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
  saveSuccess: boolean;
}

export const SecurityPolicyTab: React.FC<SecurityPolicyTabProps> = ({
  securitySettings,
  setSecuritySettings,
  blacklistText,
  setBlacklistText,
  onSave,
  saving,
  saveSuccess,
}) => {
  return (
    <form onSubmit={onSave} className="space-y-6">
      {saveSuccess && (
        <div className="p-3 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 text-xs font-medium">
          Security settings updated successfully!
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className="border-border/60 shadow-sm">
          <CardHeader>
            <CardTitle className="text-base font-semibold flex items-center gap-2">
              <Bot className="h-4 w-4 text-primary" />
              Anti-Spam & Bot Protection
            </CardTitle>
            <CardDescription>Defend client registration and login forms against automated abuse</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between p-3 rounded-lg border bg-muted/30">
              <div>
                <p className="text-xs font-semibold">Enable CAPTCHA Challenge</p>
                <p className="text-[11px] text-muted-foreground">Verify visitors on register, login, & tickets</p>
              </div>
              <input
                type="checkbox"
                checked={securitySettings.recaptcha_enabled}
                onChange={(e) =>
                  setSecuritySettings((prev) => ({ ...prev, recaptcha_enabled: e.target.checked }))
                }
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-semibold">CAPTCHA Provider</label>
              <select
                className="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                value={securitySettings.recaptcha_provider}
                onChange={(e) =>
                  setSecuritySettings((prev) => ({ ...prev, recaptcha_provider: e.target.value as any }))
                }
              >
                <option value="cloudflare_turnstile">Cloudflare Turnstile (Recommended - Invisible)</option>
                <option value="google_recaptcha">Google reCAPTCHA v2 / v3</option>
              </select>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Site Key</label>
              <Input
                placeholder="0x4AAAAAA..."
                value={securitySettings.site_key}
                onChange={(e) => setSecuritySettings((prev) => ({ ...prev, site_key: e.target.value }))}
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Secret Key</label>
              <Input
                type="password"
                placeholder="0x4AAAAAA..."
                value={securitySettings.secret_key || ''}
                onChange={(e) => setSecuritySettings((prev) => ({ ...prev, secret_key: e.target.value }))}
              />
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/60 shadow-sm">
          <CardHeader>
            <CardTitle className="text-base font-semibold flex items-center gap-2">
              <ShieldAlert className="h-4 w-4 text-primary" />
              Brute-Force & Lockout Policy
            </CardTitle>
            <CardDescription>Automatically throttle and block repeat failed login attempts</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">Max Login Attempts</label>
                <Input
                  type="number"
                  value={securitySettings.max_login_attempts}
                  onChange={(e) =>
                    setSecuritySettings((prev) => ({
                      ...prev,
                      max_login_attempts: parseInt(e.target.value) || 5,
                    }))
                  }
                />
                <p className="text-[10px] text-muted-foreground">Before IP is temporarily banned</p>
              </div>
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">Lockout Duration (Mins)</label>
                <Input
                  type="number"
                  value={securitySettings.lockout_time_minutes}
                  onChange={(e) =>
                    setSecuritySettings((prev) => ({
                      ...prev,
                      lockout_time_minutes: parseInt(e.target.value) || 15,
                    }))
                  }
                />
                <p className="text-[10px] text-muted-foreground">Cool-off period</p>
              </div>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-semibold">IP Address Blacklist (1 per line)</label>
              <Textarea
                rows={4}
                placeholder="192.0.2.1&#10;198.51.100.44"
                className="font-mono text-xs"
                value={blacklistText}
                onChange={(e) => setBlacklistText(e.target.value)}
              />
              <p className="text-[10px] text-muted-foreground">Denied from accessing all client & admin portals</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="flex justify-end">
        <Button type="submit" disabled={saving} className="gap-2">
          {saving ? 'Saving Security Policy...' : 'Save Security Policy'}
        </Button>
      </div>
    </form>
  );
};
