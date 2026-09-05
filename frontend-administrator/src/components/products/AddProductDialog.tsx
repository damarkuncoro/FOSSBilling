import React from 'react';
import { Plus } from 'lucide-react';
import { ProductItem } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

interface AddProductDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: Partial<ProductItem>;
  setForm: React.Dispatch<React.SetStateAction<Partial<ProductItem>>>;
  onTitleChange: (val: string) => void;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const AddProductDialog: React.FC<AddProductDialogProps> = ({
  open,
  onOpenChange,
  form,
  setForm,
  onTitleChange,
  onSave,
  saving,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button size="sm" className="gap-1.5 shadow-sm">
          <Plus className="h-4 w-4" />
          New Product
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-xl">
        <DialogHeader>
          <DialogTitle>Add New Product / Service</DialogTitle>
          <DialogDescription>
            Configure product type, pricing tiers, and client purchasing parameters.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Product Title</label>
              <Input
                required
                placeholder="e.g. Cloud NVMe Web Hosting"
                value={form.title || ''}
                onChange={(e) => onTitleChange(e.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">URL Slug / Identifier</label>
              <Input
                required
                placeholder="e.g. cloud-nvme-web-hosting"
                value={form.slug || ''}
                onChange={(e) => setForm((prev) => ({ ...prev, slug: e.target.value }))}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Product Type</label>
              <select
                className="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                value={form.type}
                onChange={(e) => setForm((prev) => ({ ...prev, type: e.target.value as any }))}
              >
                <option value="hosting">Web Hosting / Server</option>
                <option value="domain">Domain Registration</option>
                <option value="license">Software License Key</option>
                <option value="downloadable">Downloadable File / Asset</option>
                <option value="custom">Custom Service</option>
              </select>
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Category</label>
              <Input
                placeholder="e.g. Shared Web Hosting"
                value={form.category_name || ''}
                onChange={(e) => setForm((prev) => ({ ...prev, category_name: e.target.value }))}
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Description / Features</label>
            <Textarea
              rows={2}
              placeholder="Short description shown to clients during checkout..."
              value={form.description || ''}
              onChange={(e) => setForm((prev) => ({ ...prev, description: e.target.value }))}
            />
          </div>

          <div className="grid grid-cols-3 gap-3">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Monthly Price ($)</label>
              <Input
                type="number"
                step="0.01"
                value={form.price_monthly ?? 0}
                onChange={(e) => setForm((prev) => ({ ...prev, price_monthly: parseFloat(e.target.value) || 0 }))}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Annual Price ($)</label>
              <Input
                type="number"
                step="0.01"
                value={form.price_annually ?? 0}
                onChange={(e) => setForm((prev) => ({ ...prev, price_annually: parseFloat(e.target.value) || 0 }))}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Setup Fee ($)</label>
              <Input
                type="number"
                step="0.01"
                value={form.setup_fee ?? 0}
                onChange={(e) => setForm((prev) => ({ ...prev, setup_fee: parseFloat(e.target.value) || 0 }))}
              />
            </div>
          </div>

          <div className="flex items-center gap-2 pt-2">
            <input
              type="checkbox"
              id="is_active"
              checked={form.is_active ?? true}
              onChange={(e) => setForm((prev) => ({ ...prev, is_active: e.target.checked }))}
              className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
            />
            <label htmlFor="is_active" className="text-xs font-medium cursor-pointer">
              Publish & enable this product for immediate client ordering
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-3 border-t">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Saving...' : 'Create Product'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
