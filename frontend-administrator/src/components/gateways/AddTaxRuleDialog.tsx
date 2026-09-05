import React from 'react';
import { TaxRuleItem } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

interface AddTaxRuleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  taxForm: Partial<TaxRuleItem>;
  setTaxForm: React.Dispatch<React.SetStateAction<Partial<TaxRuleItem>>>;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const AddTaxRuleDialog: React.FC<AddTaxRuleDialogProps> = ({
  open,
  onOpenChange,
  taxForm,
  setTaxForm,
  onSave,
  saving,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Tax / VAT Rule</DialogTitle>
          <DialogDescription>Define regional percentage rate for client invoices</DialogDescription>
        </DialogHeader>

        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Rule Title</label>
            <Input
              required
              placeholder="e.g. Indonesia PPN 11%"
              value={taxForm.name || ''}
              onChange={(e) => setTaxForm((prev) => ({ ...prev, name: e.target.value }))}
            />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Country ISO Code</label>
              <Input
                required
                placeholder="e.g. ID, US, GB"
                value={taxForm.country || ''}
                onChange={(e) => setTaxForm((prev) => ({ ...prev, country: e.target.value.toUpperCase() }))}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Tax Rate Percentage (%)</label>
              <Input
                type="number"
                step="0.1"
                required
                value={taxForm.rate ?? 0}
                onChange={(e) => setTaxForm((prev) => ({ ...prev, rate: parseFloat(e.target.value) || 0 }))}
              />
            </div>
          </div>

          <div className="flex items-center gap-2 pt-2">
            <input
              type="checkbox"
              id="apply_all"
              checked={taxForm.apply_to_all_clients ?? false}
              onChange={(e) => setTaxForm((prev) => ({ ...prev, apply_to_all_clients: e.target.checked }))}
              className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
            />
            <label htmlFor="apply_all" className="text-xs font-medium cursor-pointer">
              Apply globally to all clients regardless of country
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-3 border-t">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Saving...' : 'Create Tax Rule'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
