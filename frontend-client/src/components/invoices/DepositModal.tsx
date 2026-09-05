import React, { useState } from 'react';
import { Wallet, Loader2, CheckCircle2, ArrowRight } from 'lucide-react';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { api } from '@/lib/api';

interface DepositModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  currency?: string;
  onDepositSuccess?: () => void;
}

const PRESET_AMOUNTS = [10, 25, 50, 100, 250];

export const DepositModal: React.FC<DepositModalProps> = ({
  open,
  onOpenChange,
  currency = 'USD',
  onDepositSuccess,
}) => {
  const [amount, setAmount] = useState<number>(50);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successResult, setSuccessResult] = useState<{ invoice_id: number; nr: string; total: number } | null>(null);

  const handleDeposit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!amount || amount <= 0) {
      setError('Please enter a valid deposit amount greater than 0');
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const res = await api.depositFunds(Number(amount), currency);
      setSuccessResult(res);
      if (onDepositSuccess) onDepositSuccess();
    } catch (err: any) {
      setError(err?.message || 'Failed to generate deposit invoice');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setSuccessResult(null);
    setError(null);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[480px]">
        {successResult ? (
          <div className="py-6 text-center space-y-4">
            <div className="h-12 w-12 rounded-full bg-emerald-500/10 text-emerald-500 flex items-center justify-center mx-auto">
              <CheckCircle2 className="h-6 w-6" />
            </div>
            <div>
              <h3 className="text-lg font-bold">Deposit Invoice Generated!</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Invoice <span className="font-mono font-semibold text-foreground">#{successResult.nr || successResult.invoice_id}</span> for{' '}
                <span className="font-bold text-foreground font-mono">{successResult.total} {currency}</span> has been issued to your account.
              </p>
            </div>
            <Button onClick={handleClose} className="w-full gap-2 mt-4">
              View Invoices <ArrowRight className="h-4 w-4" />
            </Button>
          </div>
        ) : (
          <form onSubmit={handleDeposit}>
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2">
                <Wallet className="h-5 w-5 text-primary" />
                Add Funds to Account Balance
              </DialogTitle>
              <DialogDescription>
                Deposit funds into your account balance for instant 1-click renewals and automated invoice settlements.
              </DialogDescription>
            </DialogHeader>

            {error && (
              <div className="mt-3 p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm">
                {error}
              </div>
            )}

            <div className="space-y-4 py-4">
              <div>
                <label className="text-xs font-semibold text-muted-foreground block mb-2">Quick Select Amount</label>
                <div className="grid grid-cols-5 gap-2">
                  {PRESET_AMOUNTS.map((val) => (
                    <Button
                      key={val}
                      type="button"
                      variant={amount === val ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => setAmount(val)}
                      className="font-mono text-xs"
                    >
                      ${val}
                    </Button>
                  ))}
                </div>
              </div>

              <div>
                <label className="text-xs font-semibold text-muted-foreground block mb-1">
                  Custom Deposit Amount ({currency})
                </label>
                <div className="relative">
                  <Input
                    type="number"
                    step="1"
                    min="1"
                    value={amount}
                    onChange={(e) => setAmount(Number(e.target.value))}
                    placeholder="Enter deposit amount..."
                    className="font-mono text-lg font-bold pl-8"
                    required
                  />
                  <span className="absolute left-3 top-2.5 text-muted-foreground font-bold text-sm">$</span>
                </div>
              </div>
            </div>

            <DialogFooter className="gap-2">
              <Button type="button" variant="outline" onClick={handleClose}>
                Cancel
              </Button>
              <Button type="submit" disabled={loading} className="gap-2">
                {loading && <Loader2 className="h-4 w-4 animate-spin" />}
                Generate Top-Up Invoice
              </Button>
            </DialogFooter>
          </form>
        )}
      </DialogContent>
    </Dialog>
  );
};
