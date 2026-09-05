import React, { useEffect, useState } from 'react';
import { FileText, RefreshCw, CheckCircle, Clock, AlertTriangle, Download } from 'lucide-react';
import { api } from '@/lib/api';
import { formatMoney, formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const Invoices: React.FC = () => {
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  const fetchInvoices = async () => {
    setLoading(true);
    try {
      const data = await api.getDashboardStats();
      setStats(data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchInvoices();
  }, []);

  const mockInvoices = [
    {
      id: 1,
      serie_nr: 'INV00001',
      client_id: 1,
      client_name: 'Budi Santoso',
      subtotal: 760000,
      tax: 83600,
      total: 843600,
      currency: 'IDR',
      status: 'paid',
      created_at: '2026-09-05T05:48:10Z',
      due_at: '2026-09-12T05:48:10Z',
    },
    {
      id: 2,
      serie_nr: 'INV00002',
      client_id: 2,
      client_name: 'Demo Customer',
      subtotal: 49.99,
      tax: 0,
      total: 49.99,
      currency: 'USD',
      status: 'unpaid',
      created_at: '2026-09-05T06:00:00Z',
      due_at: '2026-09-12T06:00:00Z',
    },
  ];

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Billing & Invoices</h1>
          <p className="text-sm text-muted-foreground">
            Monitor client billing statements, automated tax calculations, and payment settlements.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchInvoices} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {/* Overview Metric Pills */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <Card className="border-border/60">
          <CardContent className="p-4 flex items-center gap-3">
            <div className="h-10 w-10 rounded-lg bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
              <CheckCircle className="h-5 w-5" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground">Settled Invoices</p>
              <p className="text-xl font-bold">{stats?.paid_invoices || 1}</p>
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
              <p className="text-xl font-bold">{stats?.unpaid_invoices || 1}</p>
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
          <CardTitle className="text-base font-semibold">Recent Invoices ({mockInvoices.length})</CardTitle>
          <CardDescription>Automated PDF generation and multi-gateway tracking (Midtrans, Stripe, PayPal)</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Invoice No.</TableHead>
                <TableHead>Customer</TableHead>
                <TableHead>Subtotal</TableHead>
                <TableHead>Tax</TableHead>
                <TableHead>Total Due</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Issued Date</TableHead>
                <TableHead className="text-right">PDF</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {mockInvoices.map((inv) => (
                <TableRow key={inv.id}>
                  <TableCell className="font-mono text-xs font-semibold flex items-center gap-1.5">
                    <FileText className="h-3.5 w-3.5 text-muted-foreground" />
                    #{inv.serie_nr}
                  </TableCell>
                  <TableCell>
                    <p className="font-medium text-sm">{inv.client_name}</p>
                    <p className="text-[11px] text-muted-foreground">Client ID: #{inv.client_id}</p>
                  </TableCell>
                  <TableCell className="text-xs text-muted-foreground">
                    {formatMoney(inv.subtotal, inv.currency)}
                  </TableCell>
                  <TableCell className="text-xs text-muted-foreground">
                    {inv.tax > 0 ? formatMoney(inv.tax, inv.currency) : '-'}
                  </TableCell>
                  <TableCell className="font-bold text-sm">
                    {formatMoney(inv.total, inv.currency)}
                  </TableCell>
                  <TableCell>
                    <Badge variant={inv.status === 'paid' ? 'success' : 'warning'}>
                      {inv.status.toUpperCase()}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-xs text-muted-foreground">
                    {formatDate(inv.created_at)}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button variant="ghost" size="sm" className="h-7 text-xs gap-1">
                      <Download className="h-3.5 w-3.5" />
                      PDF
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};
