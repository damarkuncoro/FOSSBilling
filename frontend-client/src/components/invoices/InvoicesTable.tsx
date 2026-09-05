import React from 'react';
import { FileText, Download, CreditCard } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatMoney, formatDate } from '@/lib/utils';
import { Invoice, ClientProfile } from '@/types/api';

interface InvoicesTableProps {
  invoices: Invoice[];
  loading: boolean;
  user: ClientProfile | null;
  onPayModal: (invoice: Invoice) => void;
}

export const InvoicesTable: React.FC<InvoicesTableProps> = ({
  invoices,
  loading,
  user,
  onPayModal,
}) => {
  const handleDownloadPDF = async (inv: Invoice) => {
    try {
      const token = localStorage.getItem('auth_token') || sessionStorage.getItem('auth_token') || '';
      const url = `/api/v1/client/invoices/${inv.id}/pdf`;
      const response = await fetch(url, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!response.ok) {
        window.open(token ? `${url}?token=${encodeURIComponent(token)}` : url, '_blank');
        return;
      }
      const blob = await response.blob();
      const blobUrl = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = blobUrl;
      a.download = `Invoice-${inv.serie_nr || inv.id}.pdf`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(blobUrl);
      document.body.removeChild(a);
    } catch {
      const token = localStorage.getItem('auth_token') || '';
      const fallbackUrl = `/api/v1/client/invoices/${inv.id}/pdf`;
      window.open(token ? `${fallbackUrl}?token=${encodeURIComponent(token)}` : fallbackUrl, '_blank');
    }
  };

  return (
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
              <TableCell className="text-xs text-muted-foreground">
                {formatDate(inv.created_at)}
              </TableCell>
              <TableCell className="text-xs text-muted-foreground">
                {formatDate(inv.due_at || inv.created_at)}
              </TableCell>
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
                      onClick={() => onPayModal(inv)}
                    >
                      <CreditCard className="h-3 w-3" />
                      Pay Now
                    </Button>
                  )}
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-7 text-xs gap-1"
                    onClick={() => handleDownloadPDF(inv)}
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
  );
};
