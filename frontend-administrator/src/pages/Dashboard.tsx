import React from 'react';
import { FileText, LifeBuoy, RefreshCw, ShieldCheck, Activity, Database, Cpu } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useDashboard } from '@/hooks/useDashboard';
import { DashboardKpiGrid } from '@/components/dashboard/DashboardKpiGrid';
import { RevenueChartCard } from '@/components/dashboard/RevenueChartCard';
import { SystemActivityChartCard } from '@/components/dashboard/SystemActivityChartCard';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { formatDate } from '@/lib/utils';

export const Dashboard: React.FC = () => {
  const { t } = useTranslation();
  const { stats, recentLogs, systemStatus, activityTrend, loading, error, fetchStats, revenueTrends } = useDashboard();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Executive Dashboard</h1>
          <p className="text-sm text-muted-foreground">
            Real-time financial performance and active operations overview.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={() => fetchStats()} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          {t('refresh')}
        </Button>
      </div>

      {error && (
        <div className="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive text-sm">
          {error}
        </div>
      )}

      <DashboardKpiGrid stats={stats} loading={loading} />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          {loading ? (
             <Card className="border-border/60 shadow-sm h-[400px]">
                <CardHeader>
                  <Skeleton className="h-4 w-40" />
                  <Skeleton className="h-3 w-64" />
                </CardHeader>
                <CardContent className="flex items-end justify-between gap-2 h-[280px] pt-10">
                   {Array.from({ length: 12 }).map((_, i) => (
                     <Skeleton key={i} className="w-full" style={{ height: `${Math.random() * 100}%` }} />
                   ))}
                </CardContent>
             </Card>
          ) : (
            <RevenueChartCard data={revenueTrends} />
          )}
        </div>

        <div className="space-y-6">
          <Card className="border-border/60 shadow-sm">
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-semibold">{t('operational_health')}</CardTitle>
              <CardDescription>Live pending tasks requiring staff attention</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {loading ? (
                Array.from({ length: 3 }).map((_, i) => (
                  <div key={i} className="flex items-center justify-between p-3 rounded-lg border bg-muted/40">
                    <div className="flex items-center gap-3">
                      <Skeleton className="h-8 w-8 rounded-lg" />
                      <div className="space-y-2">
                        <Skeleton className="h-3 w-20" />
                        <Skeleton className="h-2 w-28" />
                      </div>
                    </div>
                    <Skeleton className="h-5 w-10 rounded-full" />
                  </div>
                ))
              ) : (
                <>
                  <div className="flex items-center justify-between p-3 rounded-lg border bg-muted/40">
                    <div className="flex items-center gap-3">
                      <div className="h-8 w-8 rounded-lg bg-amber-500/10 text-amber-500 flex items-center justify-center">
                        <FileText className="h-4 w-4" />
                      </div>
                      <div>
                        <p className="text-xs font-semibold">{t('unpaid_invoices')}</p>
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
                        <p className="text-xs font-semibold">{t('open_tickets')}</p>
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
                        <p className="text-xs font-semibold">{t('active_services')}</p>
                        <p className="text-[11px] text-muted-foreground">Automated provisioning</p>
                      </div>
                    </div>
                    <Badge variant="success" className="font-bold">
                      {stats?.active_orders || 0}
                    </Badge>
                  </div>
                </>
              )}
            </CardContent>
          </Card>

          <SystemActivityChartCard data={activityTrend} />

          <Card className="border-border/60 shadow-sm bg-indigo-600 text-white">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-base font-semibold">{t('system_engine')}</CardTitle>
                <Activity className="h-4 w-4 text-indigo-200" />
              </div>
              <CardDescription className="text-indigo-100">Go High-Performance Core</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center gap-3">
                <Cpu className="h-4 w-4 text-indigo-300" />
                <span className="text-xs font-mono">{systemStatus?.engine_version || 'Go 1.23 Cloud-Native'}</span>
              </div>
              <div className="flex items-center gap-3">
                <Database className="h-4 w-4 text-indigo-300" />
                <span className="text-xs font-mono">{systemStatus?.database_type || 'PostgreSQL 16'}</span>
              </div>
              <div className="mt-2 pt-2 border-t border-indigo-500/50 flex justify-between items-center">
                 <span className="text-[10px] uppercase font-bold text-indigo-200 tracking-wider">{t('uptime')}</span>
                 <span className="text-xs font-bold">{systemStatus?.uptime || 'N/A'}</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">{t('audit_trail')}</CardTitle>
          <CardDescription>Live audit trail of security-sensitive operations by staff</CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          <div className="divide-y border-t">
            {loading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="p-4 flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Skeleton className="h-4 w-24" />
                    <Skeleton className="h-4 w-48" />
                  </div>
                  <Skeleton className="h-4 w-20" />
                </div>
              ))
            ) : recentLogs.length === 0 ? (
               <div className="p-10 text-center text-muted-foreground text-sm">No recent activity recorded.</div>
            ) : (
              recentLogs.map((log) => (
                <div key={log.id} className="p-4 flex items-center justify-between hover:bg-muted/30 transition-colors">
                  <div className="flex items-center gap-4">
                    <div className="flex flex-col">
                       <span className="text-xs font-bold text-foreground">{log.staff_name || 'System'}</span>
                       <span className="text-[10px] text-muted-foreground font-mono">{log.ip_address}</span>
                    </div>
                    <div>
                      <p className="text-sm font-medium">
                        <span className="text-primary font-bold uppercase text-[10px] border rounded px-1.5 py-0.5 mr-2">{log.module}</span>
                        {log.action}
                      </p>
                      <p className="text-[11px] text-muted-foreground line-clamp-1">{log.details}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-[10px] text-muted-foreground">{formatDate(log.created_at)}</p>
                  </div>
                </div>
              ))
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Dashboard;
