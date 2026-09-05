import React, { useState } from 'react';
import { RefreshCw, CheckCircle, Wallet, Plus } from 'lucide-react';
import { useClientInvoices } from '@/hooks/useClientInvoices';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { InvoicesTable } from '@/components/invoices/InvoicesTable';
import { InvoicePayDialog } from '@/components/invoices/InvoicePayDialog';
import { DepositModal } from '@/components/invoices/DepositModal';
import { formatMoney } from '@/lib/utils';

export const Invoices: React.FC = () => {
  const {
    user,
    balance,
    invoices,
    loading,
    payModal,
    setPayModal,
    paying,
    message,
    fetchInvoices,
    handlePayBalance,
  } = useClientInvoices();
  const [depositOpen, setDepositOpen] = useState(false);

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Billing & Invoices</h1>
          <p className="text-sm text-muted-foreground">
            View billing history, download tax invoices in PDF format, and top up account balance.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={() => setDepositOpen(true)} className="gap-2">
            <Wallet className="h-4 w-4 text-emerald-500" />
            <span>Balance: {formatMoney(balance || 0, user?.currency || 'USD')}</span>
            <Plus className="h-3 w-3 ml-1" />
          </Button>
          <Button variant="outline" size="sm" onClick={fetchInvoices} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {message && (
        <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 text-sm flex items-center gap-2">
          <CheckCircle className="h-4 w-4 shrink-0" />
          <span>{message}</span>
        </div>
      )}

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">
            Invoice History ({invoices.length})
          </CardTitle>
          <CardDescription>Instant receipts and tax compliance documentation</CardDescription>
        </CardHeader>
        <CardContent>
          <InvoicesTable
            invoices={invoices}
            loading={loading}
            user={user}
            onPayModal={setPayModal}
          />
        </CardContent>
      </Card>

      <InvoicePayDialog
        invoice={payModal}
        onClose={() => setPayModal(null)}
        balance={balance}
        user={user}
        paying={paying}
        onPayBalance={handlePayBalance}
      />

      <DepositModal
        open={depositOpen}
        onOpenChange={setDepositOpen}
        currency={user?.currency || 'USD'}
        onDepositSuccess={fetchInvoices}
      />
    </div>
  );
};

export default Invoices;
