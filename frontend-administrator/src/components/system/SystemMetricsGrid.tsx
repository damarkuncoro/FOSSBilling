import React from 'react';
import { Clock, Cpu, HardDrive, Database } from 'lucide-react';
import { SystemStatusInfo } from '@/lib/api';
import { Card, CardContent } from '@/components/ui/card';

interface SystemMetricsGridProps {
  status: SystemStatusInfo;
}

export const SystemMetricsGrid: React.FC<SystemMetricsGridProps> = ({ status }) => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <Card className="border-border/60">
        <CardContent className="p-5 flex items-center gap-4">
          <div className="h-10 w-10 rounded-xl bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
            <Clock className="h-5 w-5" />
          </div>
          <div>
            <p className="text-xs text-muted-foreground font-medium">Cron Scheduler</p>
            <div className="flex items-center gap-1.5 mt-0.5">
              <span className="inline-block h-2 w-2 rounded-full bg-emerald-500" />
              <span className="font-semibold text-sm">Active & Healthy</span>
            </div>
            <p className="text-[10px] text-muted-foreground mt-0.5">{status.cron_last_run}</p>
          </div>
        </CardContent>
      </Card>

      <Card className="border-border/60">
        <CardContent className="p-5 flex items-center gap-4">
          <div className="h-10 w-10 rounded-xl bg-sky-500/10 text-sky-500 flex items-center justify-center">
            <Cpu className="h-5 w-5" />
          </div>
          <div>
            <p className="text-xs text-muted-foreground font-medium">System Load</p>
            <p className="font-semibold text-sm">{status.system_load}</p>
            <p className="text-[10px] text-muted-foreground mt-0.5">Uptime: {status.uptime}</p>
          </div>
        </CardContent>
      </Card>

      <Card className="border-border/60">
        <CardContent className="p-5 flex items-center gap-4">
          <div className="h-10 w-10 rounded-xl bg-indigo-500/10 text-indigo-500 flex items-center justify-center">
            <HardDrive className="h-5 w-5" />
          </div>
          <div>
            <p className="text-xs text-muted-foreground font-medium">Memory Usage</p>
            <p className="font-semibold text-sm">{status.memory_usage}</p>
            <p className="text-[10px] text-muted-foreground mt-0.5">Low footprint</p>
          </div>
        </CardContent>
      </Card>

      <Card className="border-border/60">
        <CardContent className="p-5 flex items-center gap-4">
          <div className="h-10 w-10 rounded-xl bg-amber-500/10 text-amber-500 flex items-center justify-center">
            <Database className="h-5 w-5" />
          </div>
          <div>
            <p className="text-xs text-muted-foreground font-medium">Database Volume</p>
            <p className="font-semibold text-sm">{status.database_size}</p>
            <p className="text-[10px] text-muted-foreground mt-0.5">{status.database_type}</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
