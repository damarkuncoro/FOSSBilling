import React from 'react';
import { DollarSign, TrendingUp, Users, Package, ArrowUpRight } from 'lucide-react';
import { formatMoney } from '@/lib/utils';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface DashboardKpiGridProps {
  stats: any;
}

export const DashboardKpiGrid: React.FC<DashboardKpiGridProps> = ({ stats }) => {
  return (
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
            from settled invoices
          </p>
        </CardContent>
      </Card>

      {/* Monthly Recurring Revenue */}
      <Card className="border-border/60 shadow-sm relative overflow-hidden bg-gradient-to-br from-card to-muted/20">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Monthly Recurring (MRR)
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-indigo-500/10 text-indigo-500 flex items-center justify-center">
            <TrendingUp className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold tracking-tight">
            {formatMoney(stats?.mrr || 0)}
          </div>
          <p className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
            <span className="text-indigo-500 font-medium inline-flex items-center">
              +12.1% <ArrowUpRight className="h-3 w-3" />
            </span>
            active subscription run-rate
          </p>
        </CardContent>
      </Card>

      {/* Total Active Clients */}
      <Card className="border-border/60 shadow-sm relative overflow-hidden bg-gradient-to-br from-card to-muted/20">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Total Clients
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-sky-500/10 text-sky-500 flex items-center justify-center">
            <Users className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold tracking-tight">
            {stats?.total_clients || 0}
          </div>
          <p className="text-xs text-muted-foreground mt-1">Verified customer accounts</p>
        </CardContent>
      </Card>

      {/* Active Hosting Orders */}
      <Card className="border-border/60 shadow-sm relative overflow-hidden bg-gradient-to-br from-card to-muted/20">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Active Services
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-amber-500/10 text-amber-500 flex items-center justify-center">
            <Package className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold tracking-tight">
            {stats?.active_orders || 0}
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            {stats?.suspended_orders || 0} suspended accounts
          </p>
        </CardContent>
      </Card>
    </div>
  );
};
