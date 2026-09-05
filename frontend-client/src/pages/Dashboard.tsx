import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import {
  Package,
  FileText,
  LifeBuoy,
  Wallet,
  ArrowRight,
  ShieldCheck,
  AlertTriangle,
  Server,
  Zap,
} from 'lucide-react';
import { api } from '@/lib/api';
import { useClientAuth } from '@/lib/auth';
import { formatMoney, formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const Dashboard: React.FC = () => {
  const { user, balance } = useClientAuth();
  const [orders, setOrders] = useState<any[]>([]);
  const [invoices, setInvoices] = useState<any[]>([]);
  const [tickets, setTickets] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      api.getOrders().catch(() => []),
      api.getInvoices().catch(() => []),
      api.getTickets().catch(() => []),
    ]).then(([ordersData, invoicesData, ticketsData]) => {
      setOrders(ordersData || []);
      setInvoices(invoicesData || []);
      setTickets(ticketsData || []);
      setLoading(false);
    });
  }, []);

  const unpaidInvoices = invoices.filter((inv) => inv.status === 'unpaid');
  const activeOrders = orders.filter((o) => o.status === 'active');

  return (
    <div className="space-y-8 animate-in fade-in-50 duration-300">
      {/* Welcome Banner */}
      <div className="p-6 md:p-8 rounded-2xl bg-gradient-to-r from-primary/10 via-primary/5 to-transparent border border-primary/20 flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
        <div className="space-y-1">
          <h1 className="text-2xl md:text-3xl font-bold tracking-tight">
            Welcome back, {user?.first_name} {user?.last_name}! 👋
          </h1>
          <p className="text-sm text-muted-foreground">
            Manage your cloud servers, track invoices, and access support tickets.
          </p>
        </div>
        <Link to="/">
          <Button className="gap-2 font-semibold shadow-md shadow-primary/20">
            <Zap className="h-4 w-4" />
            Order New Cloud Service
          </Button>
        </Link>
      </div>

      {/* Unpaid Invoice Alert Banner */}
      {unpaidInvoices.length > 0 && (
        <div className="p-4 rounded-xl bg-amber-500/10 border border-amber-500/30 text-amber-900 dark:text-amber-200 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 animate-in fade-in">
          <div className="flex items-center gap-3">
            <AlertTriangle className="h-5 w-5 text-amber-500 shrink-0" />
            <div>
              <p className="text-sm font-semibold">
                You have {unpaidInvoices.length} unpaid invoice pending payment!
              </p>
              <p className="text-xs opacity-90">Please settle your invoice to avoid service interruption.</p>
            </div>
          </div>
          <Link to="/invoices">
            <Button size="sm" variant="default" className="bg-amber-600 hover:bg-amber-700 text-white gap-1 text-xs">
              View & Pay Now <ArrowRight className="h-3.5 w-3.5" />
            </Button>
          </Link>
        </div>
      )}

      {/* Stats Cards */}
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
            <div className="text-2xl font-bold">{formatMoney(balance, user?.currency || 'IDR')}</div>
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

      {/* Services & Invoices Quick View */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Active Hosting Services */}
        <Card className="border-border/60">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base font-semibold">My Active Hosting & Services</CardTitle>
              <CardDescription>Quick view of your deployed cloud servers</CardDescription>
            </div>
            <Link to="/services">
              <Button variant="ghost" size="sm" className="text-xs gap-1">
                View All <ArrowRight className="h-3 w-3" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent className="space-y-3">
            {orders.length === 0 ? (
              <p className="text-sm text-muted-foreground py-4 text-center">No active services yet.</p>
            ) : (
              orders.slice(0, 3).map((order) => (
                <div key={order.id} className="p-3 rounded-lg border bg-muted/20 flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Server className="h-5 w-5 text-primary shrink-0" />
                    <div>
                      <p className="font-semibold text-sm">{order.title || `Service #${order.id}`}</p>
                      <p className="text-xs text-muted-foreground">Period: {order.period} • {formatMoney(order.price, order.currency)}</p>
                    </div>
                  </div>
                  <Badge variant={order.status === 'active' ? 'success' : 'warning'}>
                    {order.status}
                  </Badge>
                </div>
              ))
            )}
          </CardContent>
        </Card>

        {/* Recent Invoices */}
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
      </div>
    </div>
  );
};
