import React, { useEffect, useState } from 'react';
import { User, Key, Plus, Trash2, ShieldCheck, CheckCircle } from 'lucide-react';
import { api } from '@/lib/api';
import { useClientAuth } from '@/lib/auth';
import { formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export const Settings: React.FC = () => {
  const { user, refreshProfile } = useClientAuth();
  const [profileForm, setProfileForm] = useState({
    first_name: user?.first_name || '',
    last_name: user?.last_name || '',
    company: user?.company || '',
    country: user?.country || 'ID',
  });
  const [apiKeys, setApiKeys] = useState<any[]>([]);
  const [keyName, setKeyName] = useState('');
  const [savingProfile, setSavingProfile] = useState(false);
  const [generatingKey, setGeneratingKey] = useState(false);
  const [profileMessage, setProfileMessage] = useState<string | null>(null);

  const fetchKeys = async () => {
    try {
      const data = await api.getApiKeys();
      setApiKeys(data || []);
    } catch {
      // ignore
    }
  };

  useEffect(() => {
    fetchKeys();
  }, []);

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setSavingProfile(true);
    setProfileMessage(null);
    try {
      await api.updateProfile(profileForm);
      await refreshProfile();
      setProfileMessage('Your profile has been saved successfully!');
    } catch (err: any) {
      alert(`Update failed: ${err.message}`);
    } finally {
      setSavingProfile(false);
    }
  };

  const handleGenerateKey = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!keyName.trim()) return;

    setGeneratingKey(true);
    try {
      await api.generateApiKey(keyName.trim());
      setKeyName('');
      await fetchKeys();
    } catch (err: any) {
      alert(`Failed to generate key: ${err.message}`);
    } finally {
      setGeneratingKey(false);
    }
  };

  const handleRevokeKey = async (id: number) => {
    if (!confirm('Are you sure you want to revoke this API key?')) return;
    try {
      await api.revokeApiKey(id);
      await fetchKeys();
    } catch (err: any) {
      alert(`Failed to revoke key: ${err.message}`);
    }
  };

  return (
    <div className="space-y-8 animate-in fade-in-50 duration-300">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Account & Developer Settings</h1>
        <p className="text-sm text-muted-foreground">Manage your contact information, company profile, and programmatic API access keys.</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-start">
        {/* Profile Card */}
        <Card className="border-border/60 shadow-sm">
          <CardHeader>
            <div className="flex items-center gap-2">
              <User className="h-5 w-5 text-primary" />
              <CardTitle className="text-base font-semibold">Personal & Company Information</CardTitle>
            </div>
            <CardDescription>Update your contact details for automated billing and tax compliance</CardDescription>
          </CardHeader>
          <CardContent>
            {profileMessage && (
              <div className="mb-4 p-3 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 text-xs font-semibold flex items-center gap-2">
                <CheckCircle className="h-4 w-4" />
                <span>{profileMessage}</span>
              </div>
            )}

            <form onSubmit={handleUpdateProfile} className="space-y-4">
              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">First Name</label>
                  <Input
                    required
                    value={profileForm.first_name}
                    onChange={(e) => setProfileForm({ ...profileForm, first_name: e.target.value })}
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">Last Name</label>
                  <Input
                    required
                    value={profileForm.last_name}
                    onChange={(e) => setProfileForm({ ...profileForm, last_name: e.target.value })}
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
                  value={profileForm.company}
                  onChange={(e) => setProfileForm({ ...profileForm, company: e.target.value })}
                />
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-semibold">Country Code</label>
                <Input
                  value={profileForm.country}
                  onChange={(e) => setProfileForm({ ...profileForm, country: e.target.value })}
                />
              </div>

              <Button type="submit" disabled={savingProfile} className="font-semibold shadow-sm">
                {savingProfile ? 'Saving...' : 'Save Profile Changes'}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Developer API Keys Card */}
        <Card className="border-border/60 shadow-sm">
          <CardHeader>
            <div className="flex items-center gap-2">
              <Key className="h-5 w-5 text-primary" />
              <CardTitle className="text-base font-semibold">Developer API Keys</CardTitle>
            </div>
            <CardDescription>Programmatic keys for automated provisioning, CI/CD, and CLI integrations</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <form onSubmit={handleGenerateKey} className="flex gap-2">
              <Input
                required
                placeholder="Key label (e.g. Production Terraform Runner)..."
                value={keyName}
                onChange={(e) => setKeyName(e.target.value)}
                className="h-9 text-xs"
              />
              <Button type="submit" size="sm" disabled={generatingKey} className="gap-1.5 shrink-0">
                <Plus className="h-3.5 w-3.5" />
                Generate Key
              </Button>
            </form>

            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Label</TableHead>
                  <TableHead>API Key</TableHead>
                  <TableHead className="text-right">Action</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {apiKeys.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={3} className="text-center text-muted-foreground py-4 text-xs">
                      No API keys generated.
                    </TableCell>
                  </TableRow>
                ) : (
                  apiKeys.map((k) => (
                    <TableRow key={k.id}>
                      <TableCell className="font-medium text-xs">{k.name}</TableCell>
                      <TableCell className="font-mono text-xs text-primary">{k.key || `fb_key_${k.id}`}</TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-7 w-7 text-muted-foreground hover:text-destructive"
                          onClick={() => handleRevokeKey(k.id)}
                        >
                          <Trash2 className="h-3.5 w-3.5" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
