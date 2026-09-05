import React from 'react';
import { CreditCard, Wallet, CheckCircle } from 'lucide-react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';
import { Invoice, ClientProfile } from '@/types/api';

interface InvoicePayDialogProps {
  invoice: Invoice | null;
  onClose: () => void;
  balance: number;
  user: ClientProfile | null;
  paying: boolean;
  onPayBalance: (id: number) => void;
}

export const InvoicePayDialog: React.FC<InvoicePayDialogProps> = ({
  invoice,
  onClose,
  balance,
  user,
  paying,
  onPayBalance,
}) => {
  if (!invoice) return null;

  return (
    <Dialog open={!!invoice} onOpenChange={onClose}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <CreditCard className="h-5 w-5 text-primary" />
            <span>Settle Invoice #{invoice.serie_nr || invoice.id}</span>
          </DialogTitle>
          <DialogDescription>
            Total to pay:{' '}
            <strong className="text-foreground">
              {formatMoney(invoice.total, invoice.currency)}
            </strong>
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 my-2">
          {/* Pay with balance option */}
          <div className="p-4 rounded-xl border bg-muted/20 space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Wallet className="h-4 w-4 text-emerald-500" />
                <span className="text-xs font-semibold">Account Credit Balance</span>
              </div>
              <span className="text-xs font-bold">
                {formatMoney(balance, user?.currency || 'USD')}
              </span>
            </div>
            <Button
              className="w-full gap-2 text-xs font-semibold"
              variant="default"
              disabled={paying}
              onClick={() => onPayBalance(invoice.id)}
            >
              <CheckCircle className="h-4 w-4" />
              {paying ? 'Processing...' : 'Pay with Account Balance'}
            </Button>
          </div>

          {/* Online Payment Gateway */}
          <div className="p-4 rounded-xl border bg-muted/20 space-y-2">
            <p className="text-xs font-semibold">
              Online Gateway (Midtrans QRIS / Virtual Account / Stripe)
            </p>
            <Button
              className="w-full gap-2 text-xs font-semibold"
              variant="outline"
              onClick={() => alert('Redirecting to secure Payment Gateway Checkout...')}
            >
              <CreditCard className="h-4 w-4" />
              Pay via Online Gateway
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};
