import React, { useState, useEffect } from 'react';
import { User, ShieldCheck, ShieldAlert, Key, Save, Lock } from 'lucide-react';
import { useAuth } from '@/lib/auth';
import { request } from '@/lib/api/client';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export const AdminProfile: React.FC = () => {
  const { user, refreshUser } = useAuth();
  const [setupData, setSetupData] = useState<{ secret: string; qr_url: string } | null>(null);
  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(false);

  const handleStartSetup = async () => {
    setLoading(true);
    try {
      const res = await request<any>('/admin/auth/2fa/setup', { method: 'POST' });
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
      await request('/admin/auth/2fa/enable', {
        method: 'POST',
        body: JSON.stringify({ code })
      });
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
    if (!confirm('Disable 2FA? This is not recommended for admin accounts.')) return;
    setLoading(true);
    try {
      await request('/admin/auth/2fa/disable', { method: 'POST' });
      alert('2FA Disabled.');
      refreshUser();
    } catch (err: any) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <h1 className="text-2xl font-bold tracking-tight">My Administrator Profile</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
        <Card>
          <CardHeader>
             <CardTitle className="text-base flex items-center gap-2">
                <User className="h-4 w-4" /> Personal Information
             </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
             <div className="grid gap-2">
                <label className="text-xs font-bold text-muted-foreground uppercase">Display Name</label>
                <Input value={user?.name || ''} readOnly className="bg-muted/30" />
             </div>
             <div className="grid gap-2">
                <label className="text-xs font-bold text-muted-foreground uppercase">Email Address</label>
                <Input value={user?.email || ''} readOnly className="bg-muted/30" />
             </div>
             <div className="grid gap-2">
                <label className="text-xs font-bold text-muted-foreground uppercase">Role</label>
                <div className="px-3 py-2 rounded-lg bg-indigo-50 text-indigo-700 text-sm font-bold border border-indigo-100 w-fit">
                   {user?.role?.toUpperCase()}
                </div>
             </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
             <CardTitle className="text-base flex items-center gap-2">
                <Lock className="h-4 w-4" /> Two-Factor Authentication
             </CardTitle>
             <CardDescription className="text-xs">Enhance your administrative account security.</CardDescription>
          </CardHeader>
          <CardContent>
             {user?.two_factor_enabled ? (
               <div className="p-4 rounded-xl bg-emerald-50 border border-emerald-100 flex items-center justify-between">
                 <div className="flex items-center gap-3">
                   <ShieldCheck className="h-6 w-6 text-emerald-600" />
                   <p className="text-sm font-bold text-emerald-900">2FA is Active</p>
                 </div>
                 <Button variant="outline" size="sm" onClick={handleDisable} className="text-rose-600 hover:text-rose-700 border-rose-200">
                   Disable
                 </Button>
               </div>
             ) : setupData ? (
               <div className="space-y-4">
                  <div className="p-4 bg-muted/30 rounded-xl border border-dashed flex flex-col items-center gap-4 text-center">
                    <img
                      src={`https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${encodeURIComponent(setupData.qr_url)}`}
                      alt="QR"
                    />
                    <code className="text-xs font-bold text-indigo-600 bg-white px-2 py-1 rounded border">{setupData.secret}</code>
                  </div>
                  <div className="flex gap-2">
                    <Input
                      placeholder="Enter 6-digit code"
                      value={code}
                      onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                    />
                    <Button onClick={handleEnable} disabled={code.length !== 6 || loading}>Enable</Button>
                  </div>
               </div>
             ) : (
               <div className="p-4 rounded-xl bg-amber-50 border border-amber-100 flex items-center justify-between">
                 <div className="flex items-center gap-3">
                   <ShieldAlert className="h-6 w-6 text-amber-600" />
                   <p className="text-sm font-bold text-amber-900">2FA Not Enabled</p>
                 </div>
                 <Button onClick={handleStartSetup} size="sm">Setup 2FA</Button>
               </div>
             )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
