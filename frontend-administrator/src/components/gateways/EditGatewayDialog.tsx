import React from 'react';
import { PaymentGatewayItem } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

interface EditGatewayDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  selectedGw: PaymentGatewayItem | null;
  gwForm: Partial<PaymentGatewayItem>;
  setGwForm: React.Dispatch<React.SetStateAction<Partial<PaymentGatewayItem>>>;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const EditGatewayDialog: React.FC<EditGatewayDialogProps> = ({
  open,
  onOpenChange,
  selectedGw,
  gwForm,
  setGwForm,
  onSave,
  saving,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Configure {selectedGw?.name}</DialogTitle>
          <DialogDescription>
            Update API credentials, sandbox options, and webhook parameters.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="flex items-center justify-between p-3 rounded-lg border bg-muted/30">
            <div>
              <p className="text-xs font-semibold">Enable Test / Sandbox Mode</p>
              <p className="text-[11px] text-muted-foreground">Process test transactions without real billing</p>
            </div>
            <input
              type="checkbox"
              checked={gwForm.test_mode ?? false}
              onChange={(e) => setGwForm((prev) => ({ ...prev, test_mode: e.target.checked }))}
              className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
            />
          </div>

          {selectedGw?.type !== 'bank_transfer' ? (
            <>
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">Public Key / Client ID</label>
                <Input
                  placeholder="pk_..."
                  value={gwForm.public_key || ''}
                  onChange={(e) => setGwForm((prev) => ({ ...prev, public_key: e.target.value }))}
                />
              </div>
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">Secret Key / Secret Token</label>
                <Input
                  type="password"
                  placeholder="sk_..."
                  value={gwForm.secret_key || ''}
                  onChange={(e) => setGwForm((prev) => ({ ...prev, secret_key: e.target.value }))}
                />
              </div>
              <div className="space-y-1.5">
                <label className="text-xs font-semibold">Webhook Signing Secret (Optional)</label>
                <Input
                  placeholder="whsec_..."
                  value={gwForm.webhook_secret || ''}
                  onChange={(e) => setGwForm((prev) => ({ ...prev, webhook_secret: e.target.value }))}
                />
              </div>
            </>
          ) : (
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Transfer Instructions & Bank Details</label>
              <Textarea
                rows={4}
                value={gwForm.instructions || ''}
                onChange={(e) => setGwForm((prev) => ({ ...prev, instructions: e.target.value }))}
              />
            </div>
          )}

          <div className="flex justify-end gap-2 pt-3 border-t">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Saving...' : 'Save Configuration'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
