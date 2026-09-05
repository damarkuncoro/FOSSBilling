import React from 'react';
import { Plus } from 'lucide-react';
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

interface AddCurrencyDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: {
    code: string;
    title: string;
    conversion_rate: number;
    format: string;
  };
  setForm: React.Dispatch<
    React.SetStateAction<{
      code: string;
      title: string;
      conversion_rate: number;
      format: string;
    }>
  >;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const AddCurrencyDialog: React.FC<AddCurrencyDialogProps> = ({
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
          Add Currency
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add New Currency</DialogTitle>
          <DialogDescription>Define a new international currency and its conversion rate</DialogDescription>
        </DialogHeader>
        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">ISO Currency Code</label>
            <Input
              required
              placeholder="e.g. JPY, GBP, AUD"
              value={form.code}
              onChange={(e) => setForm((prev) => ({ ...prev, code: e.target.value }))}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Title / Description</label>
            <Input
              required
              placeholder="e.g. Japanese Yen"
              value={form.title}
              onChange={(e) => setForm((prev) => ({ ...prev, title: e.target.value }))}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Conversion Rate (relative to default)</label>
            <Input
              type="number"
              step="any"
              required
              value={form.conversion_rate}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, conversion_rate: parseFloat(e.target.value) || 1 }))
              }
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Display Format Template</label>
            <Input
              required
              placeholder="e.g. ¥ {{price}} or Rp {{price}}"
              value={form.format}
              onChange={(e) => setForm((prev) => ({ ...prev, format: e.target.value }))}
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Creating...' : 'Save Currency'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
