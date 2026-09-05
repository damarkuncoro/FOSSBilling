import React, { useEffect, useState } from 'react';
import {
  FileText,
  Download,
  CreditCard,
  Wallet,
  CheckCircle,
  RefreshCw,
  Clock,
  AlertCircle,
} from 'lucide-react';
import { api } from '@/lib/api';
import { useClientAuth } from '@/lib/auth';
import { formatMoney, formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';

export const Invoices: React.FC = () => {
  const { user, balance, refreshProfile } = useClientAuth();
  const [invoices, setInvoices] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [payModal, setPayModal] = useState<any | null>(null);
  const [paying, setPaying] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  const fetchInvoices = async () => {
    setLoading(true);
    try {
      const data = await api.getInvoices();
      setInvoices(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchInvoices();
  }, []);

  const handlePayBalance = async (id: number) => {
    setPaying(true);
    setMessage(null);
    try {
      await api.payWithBalance(id);
      setMessage(`Invoice #${id} successfully paid with account balance!`);
      setPayModal(null);
      await Promise.all([fetchInvoices(), refreshProfile()]);
    } catch (err: any) {
      alert(`Payment failed: ${err.message}`);
    } finally {
      setPaying(false);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Billing & Invoices</h1>
          <p className="text-sm text-muted-foreground">
            View billing history, download tax invoices in PDF format, and pay pending statements.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchInvoices} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {message && (
        <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 text-sm flex items-center gap-2">
          <CheckCircle className="h-4 w-4 shrink-0" />
          <span>{message}</span>
        </div>
      )}

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Invoice History ({invoices.length})</CardTitle>
          <CardDescription>Instant receipts and tax compliance documentation</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Invoice ID</TableHead>
                <TableHead>Issued Date</TableHead>
                <TableHead>Due Date</TableHead>
                <TableHead>Total Amount</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                    Loading invoices...
                  </TableCell>
                </TableRow>
              ) : invoices.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                    No billing statements found.
                  </TableCell>
                </TableRow>
              ) : (
                invoices.map((inv) => (
                  <TableRow key={inv.id}>
                    <TableCell className="font-mono text-xs font-semibold flex items-center gap-2">
                      <FileText className="h-3.5 w-3.5 text-muted-foreground" />
                      #{inv.serie_nr || inv.id}
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">{formatDate(inv.created_at)}</TableCell>
                    <TableCell className="text-xs text-muted-foreground">{formatDate(inv.due_at || inv.created_at)}</TableCell>
                    <TableCell className="font-bold text-sm">
                      {formatMoney(inv.total, inv.currency || user?.currency || 'USD')}
                    </TableCell>
                    <TableCell>
                      <Badge variant={inv.status === 'paid' ? 'success' : 'warning'}>
                        {inv.status.toUpperCase()}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        {inv.status === 'unpaid' && (
                          <Button
                            size="sm"
                            className="h-7 text-xs gap-1 font-semibold"
                            onClick={() => setPayModal(inv)}
                          >
                            <CreditCard className="h-3 w-3" />
                            Pay Now
                          </Button>
                        )}
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-7 text-xs gap-1"
                          onClick={() => alert(`Downloading Invoice PDF #${inv.serie_nr || inv.id}...`)}
                        >
                          <Download className="h-3 w-3" />
                          PDF
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Payment Selection Modal */}
      {payModal && (
        <Dialog open={!!payModal} onOpenChange={() => setPayModal(null)}>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2">
                <CreditCard className="h-5 w-5 text-primary" />
                <span>Settle Invoice #{payModal.serie_nr || payModal.id}</span>
              </DialogTitle>
              <DialogDescription>
                Total to pay: <strong className="text-foreground">{formatMoney(payModal.total, payModal.currency)}</strong>
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
                  <span className="text-xs font-bold">{formatMoney(balance, user?.currency || 'USD')}</span>
                </div>
                <Button
                  className="w-full gap-2 text-xs font-semibold"
                  variant="default"
                  disabled={paying}
                  onClick={() => handlePayBalance(payModal.id)}
                >
                  <CheckCircle className="h-4 w-4" />
                  {paying ? 'Processing...' : 'Pay with Account Balance'}
                </Button>
              </div>

              {/* Online Payment Gateway */}
              <div className="p-4 rounded-xl border bg-muted/20 space-y-2">
                <p className="text-xs font-semibold">Online Gateway (Midtrans QRIS / Virtual Account / Stripe)</p>
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
      )}
    </div>
  );
};
