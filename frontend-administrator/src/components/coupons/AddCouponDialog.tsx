import React from 'react';
import { Plus, Sparkles } from 'lucide-react';
import { CouponItem } from '@/lib/api';
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

interface AddCouponDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: Partial<CouponItem>;
  setForm: React.Dispatch<React.SetStateAction<Partial<CouponItem>>>;
  onGenerateCode: () => void;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const AddCouponDialog: React.FC<AddCouponDialogProps> = ({
  open,
  onOpenChange,
  form,
  setForm,
  onGenerateCode,
  onSave,
  saving,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button size="sm" className="gap-1.5 shadow-sm">
          <Plus className="h-4 w-4" />
          New Promo Code
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Create Promo Coupon</DialogTitle>
          <DialogDescription>Define discount values and usage restrictions</DialogDescription>
        </DialogHeader>

        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <label className="text-xs font-semibold">Promo Code</label>
              <button
                type="button"
                onClick={onGenerateCode}
                className="text-[11px] text-primary hover:underline flex items-center gap-1"
              >
                <Sparkles className="h-3 w-3" />
                Generate Code
              </button>
            </div>
            <Input
              required
              placeholder="e.g. SUMMER50"
              value={form.code || ''}
              onChange={(e) => setForm((prev) => ({ ...prev, code: e.target.value.toUpperCase() }))}
              className="font-mono uppercase font-bold"
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Discount Type</label>
              <select
                className="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                value={form.type}
                onChange={(e) => setForm((prev) => ({ ...prev, type: e.target.value as any }))}
              >
                <option value="percentage">Percentage (%)</option>
                <option value="fixed">Fixed Amount ($)</option>
              </select>
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">
                Discount Value {form.type === 'percentage' ? '(%)' : '($)'}
              </label>
              <Input
                type="number"
                step="any"
                required
                value={form.value ?? 0}
                onChange={(e) => setForm((prev) => ({ ...prev, value: parseFloat(e.target.value) || 0 }))}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Usage Limit (Max Uses)</label>
              <Input
                type="number"
                required
                value={form.max_uses ?? 100}
                onChange={(e) => setForm((prev) => ({ ...prev, max_uses: parseInt(e.target.value) || 0 }))}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Expiry Date (Optional)</label>
              <Input
                type="date"
                value={form.expires_at || ''}
                onChange={(e) => setForm((prev) => ({ ...prev, expires_at: e.target.value }))}
              />
            </div>
          </div>

          <div className="flex items-center gap-2 pt-2">
            <input
              type="checkbox"
              id="is_coupon_active"
              checked={form.is_active ?? true}
              onChange={(e) => setForm((prev) => ({ ...prev, is_active: e.target.checked }))}
              className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
            />
            <label htmlFor="is_coupon_active" className="text-xs font-medium cursor-pointer">
              Enable and activate immediately for client checkouts
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-3 border-t">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Creating...' : 'Create Coupon'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
