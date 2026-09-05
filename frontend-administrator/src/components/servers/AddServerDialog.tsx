import React from 'react';
import { Plus } from 'lucide-react';
import { ServerItem } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

interface AddServerDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: Partial<ServerItem> & { api_token?: string };
  setForm: React.Dispatch<React.SetStateAction<Partial<ServerItem> & { api_token?: string }>>;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const AddServerDialog: React.FC<AddServerDialogProps> = ({
  open,
  onOpenChange,
  form,
  setForm,
  onSave,
  saving,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button size="sm" className="gap-1.5 shadow-sm">
          <Plus className="h-4 w-4" />
          Add Server
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-xl">
        <DialogHeader>
          <DialogTitle>Connect New Server Node</DialogTitle>
          <DialogDescription>
            Enter server credentials and control panel API parameters for automatic provisioning.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Server Name</label>
              <Input
                required
                placeholder="e.g. US-East Node 1"
                value={form.name || ''}
                onChange={(e) => setForm((prev) => ({ ...prev, name: e.target.value }))}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Control Panel Manager</label>
              <select
                className="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                value={form.manager}
                onChange={(e) => setForm((prev) => ({ ...prev, manager: e.target.value as any }))}
              >
                <option value="cpanel">cPanel / WHM</option>
                <option value="hestiacp">HestiaCP</option>
                <option value="cwp">Control Web Panel (CWP)</option>
                <option value="directadmin">DirectAdmin</option>
                <option value="plesk">Plesk</option>
                <option value="custom">Custom API Driver</option>
              </select>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Server Hostname</label>
              <Input
                required
                placeholder="e.g. srv1.yourdomain.com"
                value={form.hostname || ''}
                onChange={(e) => setForm((prev) => ({ ...prev, hostname: e.target.value }))}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Primary IP Address</label>
              <Input
                required
                placeholder="e.g. 198.51.100.1"
                value={form.ip || ''}
                onChange={(e) => setForm((prev) => ({ ...prev, ip: e.target.value }))}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Primary Nameserver (NS1)</label>
              <Input
                placeholder="ns1.yourdomain.com"
                value={form.nameserver_1 || ''}
                onChange={(e) => setForm((prev) => ({ ...prev, nameserver_1: e.target.value }))}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Secondary Nameserver (NS2)</label>
              <Input
                placeholder="ns2.yourdomain.com"
                value={form.nameserver_2 || ''}
                onChange={(e) => setForm((prev) => ({ ...prev, nameserver_2: e.target.value }))}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">API Token / Access Key</label>
              <Input
                type="password"
                placeholder="WHM Token or Hestia Access Key"
                value={form.api_token || ''}
                onChange={(e) => setForm((prev) => ({ ...prev, api_token: e.target.value }))}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Max Account Capacity</label>
              <Input
                type="number"
                value={form.max_accounts ?? 100}
                onChange={(e) => setForm((prev) => ({ ...prev, max_accounts: parseInt(e.target.value) || 100 }))}
              />
            </div>
          </div>

          <div className="flex items-center gap-2 pt-2">
            <input
              type="checkbox"
              id="is_default_server"
              checked={form.is_default ?? false}
              onChange={(e) => setForm((prev) => ({ ...prev, is_default: e.target.checked }))}
              className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
            />
            <label htmlFor="is_default_server" className="text-xs font-medium cursor-pointer">
              Set as Default Provisioning Node for new hosting orders
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-3 border-t">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Connecting...' : 'Save & Register Server'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
