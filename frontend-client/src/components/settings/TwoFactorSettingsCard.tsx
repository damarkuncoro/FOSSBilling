import React, { useState } from 'react';
import { ShieldCheck, ShieldAlert, CheckCircle2, Lock } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { clientPortalApi } from '@/lib/api/client_portal';
import { useAuth } from '@/lib/auth';

export const TwoFactorSettingsCard: React.FC = () => {
  const { user, refreshUser } = useAuth();
  const [setupData, setSetupData] = useState<{ secret: string; qr_url: string } | null>(null);
  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(false);

  const handleStartSetup = async () => {
    setLoading(true);
    try {
      const res = await clientPortalApi.setupTwoFactor();
      setSetupData(res);
    } catch (err: any) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleEnable = async () => {
    setLoading(true);
    try {
      await clientPortalApi.enableTwoFactor(code);
      alert('2FA Enabled successfully!');
      setSetupData(null);
      refreshUser();
    } catch (err: any) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDisable = async () => {
    if (!confirm('Are you sure you want to disable 2FA? This will reduce your account security.')) return;
    setLoading(true);
    try {
      await clientPortalApi.disableTwoFactor();
      alert('2FA Disabled.');
      refreshUser();
    } catch (err: any) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="border-border/60 shadow-sm">
      <CardHeader>
        <div className="flex items-center gap-2 mb-1">
          <Lock className="h-5 w-5 text-indigo-500" />
          <CardTitle className="text-base font-semibold">Two-Factor Authentication (2FA)</CardTitle>
        </div>
        <CardDescription className="text-xs">
          Secure your account using TOTP (Google Authenticator, Authy, etc).
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {user?.two_factor_enabled ? (
          <div className="p-4 rounded-xl bg-emerald-50 border border-emerald-100 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <ShieldCheck className="h-6 w-6 text-emerald-600" />
              <div>
                <p className="text-sm font-bold text-emerald-900">2FA is Enabled</p>
                <p className="text-[11px] text-emerald-700">Your account is protected by an additional layer of security.</p>
              </div>
            </div>
            <Button variant="ghost" size="sm" onClick={handleDisable} disabled={loading} className="text-rose-600 hover:text-rose-700 hover:bg-rose-50 text-xs font-bold">
              Disable
            </Button>
          </div>
        ) : setupData ? (
          <div className="space-y-4 animate-in fade-in slide-in-from-top-2">
            <div className="p-4 bg-muted/30 rounded-xl border border-dashed flex flex-col items-center gap-4 text-center">
               <p className="text-xs font-medium">Scan this QR code with your authenticator app:</p>
               <div className="bg-white p-3 rounded-lg border shadow-sm">
                  {/* In a real app we'd use a QR component here. For now we show the URL as fallback */}
                  <img
                    src={`https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=${encodeURIComponent(setupData.qr_url)}`}
                    alt="2FA QR Code"
                    className="w-40 h-40"
                  />
               </div>
               <div className="space-y-1">
                  <p className="text-[10px] text-muted-foreground uppercase font-bold tracking-wider">Secret Key</p>
                  <code className="text-sm font-bold text-indigo-600">{setupData.secret}</code>
               </div>
            </div>

            <div className="space-y-2">
              <p className="text-xs font-semibold">Verify Code to Enable</p>
              <div className="flex gap-2">
                <Input
                  placeholder="6-digit code"
                  value={code}
                  onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                  className="font-mono text-center tracking-[0.5em] h-10"
                />
                <Button onClick={handleEnable} disabled={loading || code.length !== 6}>
                  Verify & Enable
                </Button>
              </div>
            </div>
          </div>
        ) : (
          <div className="p-4 rounded-xl bg-amber-50 border border-amber-100 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <ShieldAlert className="h-6 w-6 text-amber-600" />
              <div>
                <p className="text-sm font-bold text-amber-900">2FA is Not Active</p>
                <p className="text-[11px] text-amber-700">Add an extra layer of security to prevent unauthorized access.</p>
              </div>
            </div>
            <Button size="sm" onClick={handleStartSetup} disabled={loading} className="text-xs font-bold shrink-0">
              Setup 2FA
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
};
