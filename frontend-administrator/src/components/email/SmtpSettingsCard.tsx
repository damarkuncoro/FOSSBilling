import React from 'react';
import { Mail, RefreshCw, Send } from 'lucide-react';
import { MailConfig } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

interface SmtpSettingsCardProps {
  mailConfig: MailConfig;
  setMailConfig: React.Dispatch<React.SetStateAction<MailConfig>>;
  onSaveConfig: (e: React.FormEvent) => void;
  savingConfig: boolean;
  testEmail: string;
  setTestEmail: (val: string) => void;
  onSendTest: () => void;
  sendingTest: boolean;
  testStatus: { success: boolean; message: string } | null;
}

export const SmtpSettingsCard: React.FC<SmtpSettingsCardProps> = ({
  mailConfig,
  setMailConfig,
  onSaveConfig,
  savingConfig,
  testEmail,
  setTestEmail,
  onSendTest,
  sendingTest,
  testStatus,
}) => {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <Card className="lg:col-span-2 border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">SMTP Connection Configuration</CardTitle>
          <CardDescription>Configure the outbound mail transport server used for all transactional emails</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSaveConfig} className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">From Sender Name</label>
                <Input
                  required
                  placeholder="e.g. FOSSBilling Support"
                  value={mailConfig.from_name}
                  onChange={(e) => setMailConfig({ ...mailConfig, from_name: e.target.value })}
                />
              </div>
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">From Sender Email</label>
                <Input
                  type="email"
                  required
                  placeholder="noreply@yourdomain.com"
                  value={mailConfig.from_email}
                  onChange={(e) => setMailConfig({ ...mailConfig, from_email: e.target.value })}
                />
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <div className="sm:col-span-2 space-y-1.5">
                <label className="text-xs font-semibold">SMTP Hostname</label>
                <Input
                  required
                  placeholder="smtp.mailgun.org / smtp.gmail.com"
                  value={mailConfig.smtp_host}
                  onChange={(e) => setMailConfig({ ...mailConfig, smtp_host: e.target.value })}
                />
              </div>
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">SMTP Port</label>
                <Input
                  type="number"
                  required
                  value={mailConfig.smtp_port}
                  onChange={(e) => setMailConfig({ ...mailConfig, smtp_port: parseInt(e.target.value) || 587 })}
                />
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">SMTP Username</label>
                <Input
                  placeholder="postmaster@..."
                  value={mailConfig.smtp_username}
                  onChange={(e) => setMailConfig({ ...mailConfig, smtp_username: e.target.value })}
                />
              </div>
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">SMTP Password</label>
                <Input
                  type="password"
                  placeholder="••••••••••••"
                  value={mailConfig.smtp_password || ''}
                  onChange={(e) => setMailConfig({ ...mailConfig, smtp_password: e.target.value })}
                />
              </div>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Encryption Type</label>
              <select
                className="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                value={mailConfig.smtp_encryption}
                onChange={(e) => setMailConfig({ ...mailConfig, smtp_encryption: e.target.value as any })}
              >
                <option value="tls">TLS (STARTTLS - Recommended on Port 587)</option>
                <option value="ssl">SSL (Port 465)</option>
                <option value="none">None (Plaintext - Port 25)</option>
              </select>
            </div>

            <div className="flex justify-end pt-3 border-t">
              <Button type="submit" disabled={savingConfig} className="gap-2">
                {savingConfig ? 'Saving...' : 'Save Mail Settings'}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card className="border-border/60 shadow-sm h-fit">
        <CardHeader>
          <CardTitle className="text-base font-semibold flex items-center gap-2">
            <Send className="h-4 w-4 text-primary" />
            Send Test Email
          </CardTitle>
          <CardDescription>Verify your SMTP handshake and delivery instantly</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Recipient Email Address</label>
            <Input
              type="email"
              placeholder="admin@example.com"
              value={testEmail}
              onChange={(e) => setTestEmail(e.target.value)}
            />
          </div>

          <Button
            variant="outline"
            className="w-full gap-2 text-xs"
            onClick={onSendTest}
            disabled={sendingTest || !testEmail}
          >
            {sendingTest ? <RefreshCw className="h-3.5 w-3.5 animate-spin" /> : <Mail className="h-3.5 w-3.5" />}
            {sendingTest ? 'Sending Test...' : 'Send Test Message'}
          </Button>

          {testStatus && (
            <div
              className={`p-3 rounded-lg border text-xs ${
                testStatus.success
                  ? 'bg-emerald-500/10 border-emerald-500/20 text-emerald-600 dark:text-emerald-400'
                  : 'bg-destructive/10 border-destructive/20 text-destructive'
              }`}
            >
              <p className="font-semibold">{testStatus.success ? 'Success' : 'Error'}</p>
              <p className="text-[11px] mt-0.5">{testStatus.message}</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};
