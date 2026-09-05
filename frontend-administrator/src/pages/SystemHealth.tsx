import React from 'react';
import {
  Activity,
  RefreshCw,
  Clock,
  Trash2,
  Play,
  CheckCircle2,
  Server,
  Download,
} from 'lucide-react';
import { useSystemHealth } from '@/hooks/useSystemHealth';
import { SystemMetricsGrid } from '@/components/system/SystemMetricsGrid';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const SystemHealth: React.FC = () => {
  const {
    status,
    loading,
    runningCron,
    clearingCache,
    actionMessage,
    setActionMessage,
    fetchStatus,
    handleRunCron,
    handleClearCache,
    handleExportBackup,
    cronTasks,
  } = useSystemHealth();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">System Health & Cron Diagnostics</h1>
          <p className="text-sm text-muted-foreground">
            Monitor background cron task schedulers, runtime engine performance, and database maintenance.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchStatus} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          <Button size="sm" className="gap-1.5 shadow-sm" onClick={handleRunCron} disabled={runningCron}>
            <Play className={`h-3.5 w-3.5 ${runningCron ? 'animate-spin' : ''}`} />
            {runningCron ? 'Running Cron...' : 'Run Cron Now'}
          </Button>
        </div>
      </div>

      {actionMessage && (
        <div
          className={`p-4 rounded-xl border flex items-center justify-between text-sm ${
            actionMessage.success
              ? 'bg-emerald-500/10 border-emerald-500/20 text-emerald-600 dark:text-emerald-400'
              : 'bg-destructive/10 border-destructive/20 text-destructive'
          }`}
        >
          <div className="flex items-center gap-2">
            <CheckCircle2 className="h-4 w-4 shrink-0" />
            <span>{actionMessage.text}</span>
          </div>
          <Button variant="ghost" size="sm" className="h-7 text-xs" onClick={() => setActionMessage(null)}>
            Dismiss
          </Button>
        </div>
      )}

      <SystemMetricsGrid status={status} />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2 border-border/60 shadow-sm">
          <CardHeader>
            <CardTitle className="text-base font-semibold flex items-center gap-2">
              <Activity className="h-4 w-4 text-primary" />
              Automated Background Schedulers
            </CardTitle>
            <CardDescription>
              Cron tasks executed automatically by <code>bin/cron.php</code> or the internal Go scheduler daemon
            </CardDescription>
          </CardHeader>
          <CardContent className="p-0">
            <div className="divide-y">
              {cronTasks.map((t, idx) => (
                <div key={idx} className="p-4 flex items-center justify-between hover:bg-muted/30 transition-colors">
                  <div>
                    <p className="text-sm font-semibold">{t.name}</p>
                    <p className="text-xs text-muted-foreground flex items-center gap-1.5 mt-0.5">
                      <Clock className="h-3 w-3" />
                      {t.schedule}
                    </p>
                  </div>
                  <Badge variant="success" className="gap-1 text-[10px]">
                    <CheckCircle2 className="h-3 w-3" />
                    {t.status}
                  </Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <div className="space-y-6">
          <Card className="border-border/60 shadow-sm">
            <CardHeader>
              <CardTitle className="text-base font-semibold flex items-center gap-2">
                <Server className="h-4 w-4 text-primary" />
                Environment Details
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3 text-xs">
              <div className="flex justify-between border-b pb-2">
                <span className="text-muted-foreground">Core Release</span>
                <span className="font-mono font-semibold">{status.engine_version}</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="text-muted-foreground">Go Runtime</span>
                <span className="font-mono">{status.go_version}</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="text-muted-foreground">Active Sessions</span>
                <span className="font-semibold">{status.active_sessions} Active</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Database Engine</span>
                <span className="font-mono">{status.database_type}</span>
              </div>
            </CardContent>
          </Card>

          <Card className="border-border/60 shadow-sm">
            <CardHeader>
              <CardTitle className="text-base font-semibold">Maintenance Tools</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <Button
                variant="outline"
                className="w-full justify-start gap-2 text-xs"
                onClick={handleClearCache}
                disabled={clearingCache}
              >
                <Trash2 className="h-3.5 w-3.5 text-amber-500" />
                {clearingCache ? 'Purging Cache...' : 'Purge All System Cache'}
              </Button>

              <Button
                variant="outline"
                className="w-full justify-start gap-2 text-xs"
                onClick={handleExportBackup}
              >
                <Download className="h-3.5 w-3.5 text-indigo-500" />
                Export System Diagnostics JSON
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default SystemHealth;
