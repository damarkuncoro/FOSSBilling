import React from 'react';
import { DollarSign, TrendingUp, Users, Activity } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { formatMoney } from '@/lib/utils';
import { FinancialReportSummary } from '@/types/modules';

interface RevenueMetricsGridProps {
  report: FinancialReportSummary | null;
}

export const RevenueMetricsGrid: React.FC<RevenueMetricsGridProps> = ({ report }) => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {/* MRR */}
      <Card className="border-border/60">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">
            Monthly Recurring (MRR)
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
            <DollarSign className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatMoney(report?.mrr || 0)}</div>
          <p className="text-xs text-muted-foreground mt-1">+8.4% from last month</p>
        </CardContent>
      </Card>

      {/* ARR */}
      <Card className="border-border/60">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">
            Annual Run Rate (ARR)
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-primary/10 text-primary flex items-center justify-center">
            <TrendingUp className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatMoney(report?.arr || 0)}</div>
          <p className="text-xs text-muted-foreground mt-1">Projected annual gross</p>
        </CardContent>
      </Card>

      {/* Active Subscriptions */}
      <Card className="border-border/60">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">
            Active Subscriptions
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-blue-500/10 text-blue-500 flex items-center justify-center">
            <Users className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{report?.active_subscriptions || 0}</div>
          <p className="text-xs text-muted-foreground mt-1">Hosting, VPS & Licenses</p>
        </CardContent>
      </Card>

      {/* Churn Rate */}
      <Card className="border-border/60">
        <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">
            Customer Churn Rate
          </CardTitle>
          <div className="h-8 w-8 rounded-lg bg-amber-500/10 text-amber-500 flex items-center justify-center">
            <Activity className="h-4 w-4" />
          </div>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{report?.churn_rate || 0}%</div>
          <p className="text-xs text-muted-foreground mt-1">Healthy SaaS baseline (&lt; 2%)</p>
        </CardContent>
      </Card>
    </div>
  );
};
