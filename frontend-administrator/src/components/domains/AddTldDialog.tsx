import React from 'react';
import { Plus } from 'lucide-react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

interface AddTldDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: any;
  onChange: (form: any) => void;
  onSubmit: (e: React.FormEvent) => void;
  saving: boolean;
}

export const AddTldDialog: React.FC<AddTldDialogProps> = ({
  open,
  onOpenChange,
  form,
  onChange,
  onSubmit,
  saving,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button size="sm" className="gap-1.5 shadow-sm">
          <Plus className="h-4 w-4" />
          Add New TLD
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Configure Domain TLD Pricing</DialogTitle>
          <DialogDescription>Define registration and renewal rates for top-level domains</DialogDescription>
        </DialogHeader>
        <form onSubmit={onSubmit} className="space-y-4 pt-2">
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">TLD Extension</label>
              <Input
                required
                placeholder=".com or .id"
                value={form.tld}
                onChange={(e) => onChange({ ...form, tld: e.target.value })}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Registrar Provider</label>
              <select
                value={form.registrar}
                onChange={(e) => onChange({ ...form, registrar: e.target.value })}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus:outline-none focus:ring-1 focus:ring-ring"
              >
                <option value="namecheap">Namecheap</option>
                <option value="enom">eNom</option>
                <option value="resellerclub">ResellerClub</option>
                <option value="custom">DigitalRegistrar (.ID)</option>
              </select>
            </div>
          </div>

          <div className="grid grid-cols-3 gap-3">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Register ($)</label>
              <Input
                type="number"
                step="0.01"
                required
                value={form.price_registration}
                onChange={(e) => onChange({ ...form, price_registration: parseFloat(e.target.value) || 0 })}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Renew ($)</label>
              <Input
                type="number"
                step="0.01"
                required
                value={form.price_renewal}
                onChange={(e) => onChange({ ...form, price_renewal: parseFloat(e.target.value) || 0 })}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Transfer ($)</label>
              <Input
                type="number"
                step="0.01"
                required
                value={form.price_transfer}
                onChange={(e) => onChange({ ...form, price_transfer: parseFloat(e.target.value) || 0 })}
              />
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Saving...' : 'Save TLD'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
