import React, { useEffect, useState } from 'react';
import {
  DollarSign,
  TrendingUp,
  Users,
  Package,
  FileText,
  LifeBuoy,
  RefreshCw,
  ArrowUpRight,
  ShieldCheck,
} from 'lucide-react';
import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, Tooltip, CartesianGrid } from 'recharts';
import { api } from '@/lib/api';
import { formatMoney } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

const mockRevenueTrends = [
  { month: 'Jan', revenue: 12400, mrr: 8200 },
  { month: 'Feb', revenue: 15800, mrr: 9400 },
  { month: 'Mar', revenue: 19200, mrr: 11000 },
  { month: 'Apr', revenue: 24500, mrr: 14500 },
  { month: 'May', revenue: 31000, mrr: 18200 },
  { month: 'Jun', revenue: 42000, mrr: 23500 },
  { month: 'Jul', revenue: 56000, mrr: 29000 },
];

export const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await api.getDashboardStats();
      setStats(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load dashboard metrics');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, []);

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Executive Dashboard</h1>
          <p className="text-sm text-muted-foreground">
            Real-time financial performance and active operations overview.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchStats} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh Data
          </Button>
        </div>
      </div>

      {error && (
        <div className="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive text-sm">
          {error}
        </div>
      )}

      {/* KPI Cards Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Total Revenue */}
        <Card className="border-border/60 shadow-sm relative overflow-hidden bg-gradient-to-br from-card to-muted/20">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Total Revenue
            </CardTitle>
            <div className="h-8 w-8 rounded-lg bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
              <DollarSign className="h-4 w-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold tracking-tight">
              {formatMoney(stats?.total_revenue || 0)}
            </div>
            <p className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
              <span className="text-emerald-500 font-medium inline-flex items-center">
                +18.4% <ArrowUpRight className="h-3 w-3" />
              </span>
              from all settled invoices
            </p>
          </CardContent>
        </Card>

        {/* MRR */}
        <Card className="border-border/60 shadow-sm relative overflow-hidden bg-gradient-to-br from-card to-muted/20">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Monthly Recurring (MRR)
            </CardTitle>
            <div className="h-8 w-8 rounded-lg bg-primary/10 text-primary flex items-center justify-center">
              <TrendingUp className="h-4 w-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold tracking-tight">
              {formatMoney(stats?.mrr || 0)}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              ARR Run-rate: <span className="font-semibold text-foreground">{formatMoney(stats?.arr || 0)}</span>
            </p>
          </CardContent>
        </Card>

        {/* Total Clients */}
        <Card className="border-border/60 shadow-sm relative overflow-hidden bg-gradient-to-br from-card to-muted/20">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Active Clients
            </CardTitle>
            <div className="h-8 w-8 rounded-lg bg-blue-500/10 text-blue-500 flex items-center justify-center">
              <Users className="h-4 w-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold tracking-tight">
              {stats?.total_clients || 0}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              Registered verified customers
            </p>
          </CardContent>
        </Card>

        {/* Active Orders */}
        <Card className="border-border/60 shadow-sm relative overflow-hidden bg-gradient-to-br from-card to-muted/20">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Active Services
            </CardTitle>
            <div className="h-8 w-8 rounded-lg bg-purple-500/10 text-purple-500 flex items-center justify-center">
              <Package className="h-4 w-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold tracking-tight">
              {stats?.active_orders || 0}
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              {stats?.suspended_orders || 0} suspended • {stats?.pending_orders || 0} pending
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Revenue Growth Chart & Operational Status */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Chart */}
        <Card className="lg:col-span-2 border-border/60 shadow-sm">
          <CardHeader>
            <CardTitle className="text-base font-semibold">Revenue & MRR Trajectory</CardTitle>
            <CardDescription>Monthly billing growth trends across all gateways</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[280px] w-full">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={mockRevenueTrends} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
                  <defs>
                    <linearGradient id="colorRevenue" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.4} />
                      <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                    </linearGradient>
                    <linearGradient id="colorMrr" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#10b981" stopOpacity={0.4} />
                      <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" className="stroke-muted/40" />
                  <XAxis dataKey="month" className="text-[11px] text-muted-foreground" />
                  <YAxis className="text-[11px] text-muted-foreground" tickFormatter={(v) => `$${v / 1000}k`} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: 'hsl(var(--card))',
                      borderColor: 'hsl(var(--border))',
                      borderRadius: '0.5rem',
                      fontSize: '12px',
                    }}
                  />
                  <Area type="monotone" dataKey="revenue" stroke="#3b82f6" strokeWidth={2} fillOpacity={1} fill="url(#colorRevenue)" name="Total Revenue" />
                  <Area type="monotone" dataKey="mrr" stroke="#10b981" strokeWidth={2} fillOpacity={1} fill="url(#colorMrr)" name="MRR" />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        {/* Operational Status */}
        <Card className="border-border/60 shadow-sm flex flex-col justify-between">
          <CardHeader>
            <CardTitle className="text-base font-semibold">Service Health & Queue</CardTitle>
            <CardDescription>Live pending actions and ticket queue</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="p-3 rounded-lg border bg-muted/20 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <FileText className="h-5 w-5 text-amber-500" />
                <div>
                  <p className="text-xs font-semibold">Unpaid Invoices</p>
                  <p className="text-[11px] text-muted-foreground">Awaiting client settlement</p>
                </div>
              </div>
              <Badge variant="warning">{stats?.unpaid_invoices || 0}</Badge>
            </div>

            <div className="p-3 rounded-lg border bg-muted/20 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <LifeBuoy className="h-5 w-5 text-blue-500" />
                <div>
                  <p className="text-xs font-semibold">Open Support Tickets</p>
                  <p className="text-[11px] text-muted-foreground">Customer inquiries pending reply</p>
                </div>
              </div>
              <Badge variant="info">{stats?.open_tickets || 0}</Badge>
            </div>

            <div className="p-3 rounded-lg border bg-muted/20 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <ShieldCheck className="h-5 w-5 text-emerald-500" />
                <div>
                  <p className="text-xs font-semibold">Provisioning Daemon</p>
                  <p className="text-[11px] text-muted-foreground">cPanel, DirectAdmin, Plesk SPI</p>
                </div>
              </div>
              <Badge variant="success">Online</Badge>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
