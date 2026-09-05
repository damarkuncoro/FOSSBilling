import React from 'react';
import { RefreshCw } from 'lucide-react';
import { useAntiSpam } from '@/hooks/useAntiSpam';
import { Button } from '@/components/ui/button';
import { CaptchaSettingsCard } from '@/components/antispam/CaptchaSettingsCard';
import { BlacklistIpCard } from '@/components/antispam/BlacklistIpCard';

export const AntiSpam: React.FC = () => {
  const {
    settings,
    setSettings,
    newIp,
    setNewIp,
    loading,
    saving,
    saveMessage,
    handleSave,
    handleAddIp,
    handleRemoveIp,
  } = useAntiSpam();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Anti-Spam & Bot Shield Protection</h1>
          <p className="text-sm text-muted-foreground">
            Configure automated Captcha validation, login rate limiting, and malicious IP blacklist.
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-start">
        <CaptchaSettingsCard
          settings={settings}
          onChange={setSettings}
          onSubmit={handleSave}
          saving={saving}
          message={saveMessage}
        />

        <BlacklistIpCard
          blacklist={settings.ip_blacklist}
          newIp={newIp}
          onNewIpChange={setNewIp}
          onAddIp={handleAddIp}
          onRemoveIp={handleRemoveIp}
        />
      </div>
    </div>
  );
};
