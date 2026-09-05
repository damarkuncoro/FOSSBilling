import React from 'react';
import { Link } from 'react-router-dom';
import { ArrowRight } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatMoney, formatDate } from '@/lib/utils';
import { Invoice } from '@/types/api';

interface RecentInvoicesListProps {
  invoices: Invoice[];
}

export const RecentInvoicesList: React.FC<RecentInvoicesListProps> = ({ invoices }) => {
  return (
    <Card className="border-border/60">
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle className="text-base font-semibold">Billing Invoices</CardTitle>
          <CardDescription>Recent transaction records</CardDescription>
        </div>
        <Link to="/invoices">
          <Button variant="ghost" size="sm" className="text-xs gap-1">
            View All <ArrowRight className="h-3 w-3" />
          </Button>
        </Link>
      </CardHeader>
      <CardContent className="space-y-3">
        {invoices.length === 0 ? (
          <p className="text-sm text-muted-foreground py-4 text-center">No invoices recorded.</p>
        ) : (
          invoices.slice(0, 3).map((inv) => (
            <div key={inv.id} className="p-3 rounded-lg border bg-muted/20 flex items-center justify-between">
              <div>
                <p className="font-semibold text-sm">Invoice #{inv.serie_nr || inv.id}</p>
                <p className="text-xs text-muted-foreground">{formatDate(inv.created_at)}</p>
              </div>
              <div className="text-right flex items-center gap-3">
                <span className="font-bold text-sm">{formatMoney(inv.total, inv.currency)}</span>
                <Badge variant={inv.status === 'paid' ? 'success' : 'warning'}>
                  {inv.status}
                </Badge>
              </div>
            </div>
          ))
        )}
      </CardContent>
    </Card>
  );
};
