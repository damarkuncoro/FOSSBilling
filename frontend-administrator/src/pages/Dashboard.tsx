import React from 'react';
import { FileText, LifeBuoy, RefreshCw, ShieldCheck } from 'lucide-react';
import { useDashboard } from '@/hooks/useDashboard';
import { DashboardKpiGrid } from '@/components/dashboard/DashboardKpiGrid';
import { RevenueChartCard } from '@/components/dashboard/RevenueChartCard';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const Dashboard: React.FC = () => {
  const { stats, loading, error, fetchStats, revenueTrends } = useDashboard();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Executive Dashboard</h1>
          <p className="text-sm text-muted-foreground">
            Real-time financial performance and active operations overview.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchStats} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh Data
        </Button>
      </div>

      {error && (
        <div className="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive text-sm">
          {error}
        </div>
      )}

      <DashboardKpiGrid stats={stats} />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <RevenueChartCard data={revenueTrends} />
        </div>

        <div className="space-y-6">
          <Card className="border-border/60 shadow-sm">
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-semibold">Operational Health</CardTitle>
              <CardDescription>Live pending tasks requiring staff attention</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between p-3 rounded-lg border bg-muted/40">
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-lg bg-amber-500/10 text-amber-500 flex items-center justify-center">
                    <FileText className="h-4 w-4" />
                  </div>
                  <div>
                    <p className="text-xs font-semibold">Unpaid Invoices</p>
                    <p className="text-[11px] text-muted-foreground">Awaiting client settlement</p>
                  </div>
                </div>
                <Badge variant="secondary" className="font-bold text-amber-500">
                  {stats?.unpaid_invoices || 0}
                </Badge>
              </div>

              <div className="flex items-center justify-between p-3 rounded-lg border bg-muted/40">
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-lg bg-rose-500/10 text-rose-500 flex items-center justify-center">
                    <LifeBuoy className="h-4 w-4" />
                  </div>
                  <div>
                    <p className="text-xs font-semibold">Open Support Tickets</p>
                    <p className="text-[11px] text-muted-foreground">Client inquiries awaiting reply</p>
                  </div>
                </div>
                <Badge variant="destructive" className="font-bold">
                  {stats?.open_tickets || 0}
                </Badge>
              </div>

              <div className="flex items-center justify-between p-3 rounded-lg border bg-muted/40">
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-lg bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
                    <ShieldCheck className="h-4 w-4" />
                  </div>
                  <div>
                    <p className="text-xs font-semibold">Active Orders</p>
                    <p className="text-[11px] text-muted-foreground">Automated provisioning</p>
                  </div>
                </div>
                <Badge variant="success" className="font-bold">
                  {stats?.active_orders || 0}
                </Badge>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
