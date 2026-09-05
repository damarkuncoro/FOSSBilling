import React, { useState } from 'react';
import { FileText, RefreshCw, CheckCircle, Clock, Download, Plus, FileSpreadsheet } from 'lucide-react';
import { useInvoices } from '@/hooks/useInvoices';
import { formatMoney, formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { CreateInvoiceDialog } from '@/components/invoices/CreateInvoiceDialog';

export const Invoices: React.FC = () => {
  const { stats, invoices, clients, loading, fetchInvoices, exportInvoicesToCSV, createInvoice } = useInvoices();
  const [createOpen, setCreateOpen] = useState(false);

  const handleDownloadPDF = async (inv: { id: number; nr?: string; serie_nr?: string }) => {
    try {
      const token = localStorage.getItem('fossbilling_admin_token') || sessionStorage.getItem('fossbilling_admin_token') || '';
      const url = `/api/v1/client/invoices/${inv.id}/pdf`;
      const response = await fetch(url, { headers: token ? { Authorization: `Bearer ${token}` } : {} });
      if (!response.ok) {
        window.open(token ? `${url}?token=${encodeURIComponent(token)}` : url, '_blank');
        return;
      }
      const blob = await response.blob();
      const blobUrl = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = blobUrl;
      a.download = `Invoice-${inv.nr || inv.serie_nr || inv.id}.pdf`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(blobUrl);
      document.body.removeChild(a);
    } catch {
      const token = localStorage.getItem('fossbilling_admin_token') || '';
      window.open(`/api/v1/client/invoices/${inv.id}/pdf?token=${encodeURIComponent(token)}`, '_blank');
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Billing & Invoices</h1>
          <p className="text-sm text-muted-foreground">
            Monitor client billing statements, generate custom invoices, and export financial records.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={exportInvoicesToCSV} disabled={invoices.length === 0} className="gap-2">
            <FileSpreadsheet className="h-4 w-4" /> Export CSV
          </Button>
          <Button variant="outline" size="sm" onClick={fetchInvoices} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} /> Refresh
          </Button>
          <Button size="sm" onClick={() => setCreateOpen(true)} className="gap-2">
            <Plus className="h-4 w-4" /> Create Invoice
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <Card className="border-border/60">
          <CardContent className="p-4 flex items-center gap-3">
            <div className="h-10 w-10 rounded-lg bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
              <CheckCircle className="h-5 w-5" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground">Settled Invoices</p>
              <p className="text-xl font-bold">{stats?.paid_invoices || invoices.filter((i) => i.status === 'paid').length}</p>
            </div>
          </CardContent>
        </Card>
        <Card className="border-border/60">
          <CardContent className="p-4 flex items-center gap-3">
            <div className="h-10 w-10 rounded-lg bg-amber-500/10 text-amber-500 flex items-center justify-center">
              <Clock className="h-5 w-5" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground">Unpaid Invoices</p>
              <p className="text-xl font-bold">{stats?.unpaid_invoices || invoices.filter((i) => i.status === 'unpaid').length}</p>
            </div>
          </CardContent>
        </Card>
        <Card className="border-border/60">
          <CardContent className="p-4 flex items-center gap-3">
            <div className="h-10 w-10 rounded-lg bg-blue-500/10 text-blue-500 flex items-center justify-center">
              <FileText className="h-5 w-5" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground">Total Revenue Collected</p>
              <p className="text-xl font-bold">{formatMoney(stats?.total_revenue || 843600, 'IDR')}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Recent Invoices ({invoices.length})</CardTitle>
          <CardDescription>Automated PDF generation and multi-gateway tracking (Midtrans, Stripe, PayPal)</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Invoice No.</TableHead>
                <TableHead>Client ID</TableHead>
                <TableHead>Subtotal</TableHead>
                <TableHead>Tax</TableHead>
                <TableHead>Total Due</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Issued Date</TableHead>
                <TableHead className="text-right">PDF</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">Loading invoices...</TableCell>
                </TableRow>
              ) : invoices.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">No invoices found.</TableCell>
                </TableRow>
              ) : (
                invoices.map((inv) => (
                  <TableRow key={inv.id}>
                    <TableCell className="font-mono text-xs font-semibold flex items-center gap-1.5">
                      <FileText className="h-3.5 w-3.5 text-muted-foreground" />
                      #{inv.nr || inv.serie_nr || inv.id}
                    </TableCell>
                    <TableCell className="text-xs font-mono">Client #{inv.client_id}</TableCell>
                    <TableCell className="text-xs text-muted-foreground">{formatMoney(inv.subtotal, inv.currency)}</TableCell>
                    <TableCell className="text-xs text-muted-foreground">{inv.tax > 0 ? formatMoney(inv.tax, inv.currency) : '-'}</TableCell>
                    <TableCell className="font-bold text-sm">{formatMoney(inv.total, inv.currency)}</TableCell>
                    <TableCell>
                      <Badge variant={inv.status === 'paid' ? 'success' : inv.status === 'unpaid' ? 'warning' : 'destructive'}>
                        {(inv.status || 'unpaid').toUpperCase()}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">{formatDate(inv.created_at)}</TableCell>
                    <TableCell className="text-right">
                      <Button variant="ghost" size="sm" className="h-7 text-xs gap-1" onClick={() => handleDownloadPDF(inv)}>
                        <Download className="h-3.5 w-3.5" /> PDF
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <CreateInvoiceDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        clients={clients}
        onInvoiceCreated={fetchInvoices}
        apiCreateInvoice={createInvoice}
      />
    </div>
  );
};

export default Invoices;
