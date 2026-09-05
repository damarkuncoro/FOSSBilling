import React from 'react';
import { Wallet, Package, FileText, LifeBuoy } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { formatMoney } from '@/lib/utils';
import { ClientProfile, Order, Invoice, SupportTicket } from '@/types/api';

interface DashboardStatsGridProps {
  user: ClientProfile | null;
  balance: number;
  orders: Order[];
  invoices: Invoice[];
  tickets: SupportTicket[];
  unpaidInvoices: Invoice[];
  activeOrders: Order[];
}

export const DashboardStatsGrid: React.FC<DashboardStatsGridProps> = ({
  user,
  balance,
  invoices,
  tickets,
  unpaidInvoices,
  activeOrders,
}) => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {/* Balance Card */}
      <Card className="border-border/60">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">
            Account Credit Balance
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
            <Wallet className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {formatMoney(balance, user?.currency || 'IDR')}
          </div>
          <p className="text-xs text-muted-foreground mt-1">Available for automatic checkout</p>
        </CardContent>
      </Card>

      {/* Active Services */}
      <Card className="border-border/60">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">
            Active Services
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-primary/10 text-primary flex items-center justify-center">
            <Package className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{activeOrders.length}</div>
          <p className="text-xs text-muted-foreground mt-1">Hosting, VPS & Licenses active</p>
        </CardContent>
      </Card>

      {/* Invoices */}
      <Card className="border-border/60">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">
            Total Invoices
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-blue-500/10 text-blue-500 flex items-center justify-center">
            <FileText className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{invoices.length}</div>
          <p className="text-xs text-muted-foreground mt-1">{unpaidInvoices.length} unpaid</p>
        </CardContent>
      </Card>

      {/* Support Tickets */}
      <Card className="border-border/60">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">
            Support Tickets
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-purple-500/10 text-purple-500 flex items-center justify-center">
            <LifeBuoy className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{tickets.length}</div>
          <p className="text-xs text-muted-foreground mt-1">Helpdesk requests</p>
        </CardContent>
      </Card>
    </div>
  );
};
