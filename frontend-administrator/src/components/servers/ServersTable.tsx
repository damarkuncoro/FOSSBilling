import React from 'react';
import { RefreshCw, Activity, Trash2 } from 'lucide-react';
import { ServerItem } from '@/lib/api';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

interface ServersTableProps {
  servers: ServerItem[];
  loading: boolean;
  testingId: number | null;
  onTestConnection: (id: number) => void;
  onDelete: (id: number) => void;
}

export const ServersTable: React.FC<ServersTableProps> = ({
  servers,
  loading,
  testingId,
  onTestConnection,
  onDelete,
}) => {
  if (loading) {
    return (
      <div className="text-center py-10 text-muted-foreground">
        <RefreshCw className="h-5 w-5 animate-spin mx-auto mb-2 text-primary" />
        Loading servers...
      </div>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Server Name & Host</TableHead>
          <TableHead>Manager</TableHead>
          <TableHead>IP Address</TableHead>
          <TableHead>Capacity Usage</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {servers.map((srv) => {
          const usagePct = Math.round((srv.active_accounts / (srv.max_accounts || 1)) * 100);
          return (
            <TableRow key={srv.id} className="hover:bg-muted/40 transition-colors">
              <TableCell>
                <div>
                  <div className="font-semibold text-sm flex items-center gap-2">
                    {srv.name}
                    {srv.is_default && (
                      <Badge variant="outline" className="text-[10px] px-1.5 py-0 border-primary text-primary">
                        Default
                      </Badge>
                    )}
                  </div>
                  <span className="font-mono text-xs text-muted-foreground">{srv.hostname}</span>
                </div>
              </TableCell>
              <TableCell>
                <Badge variant="secondary" className="uppercase text-[10px] font-mono">
                  {srv.manager}
                </Badge>
              </TableCell>
              <TableCell className="font-mono text-xs text-muted-foreground">{srv.ip}</TableCell>
              <TableCell>
                <div className="space-y-1 w-36">
                  <div className="flex justify-between text-[11px]">
                    <span>{srv.active_accounts} / {srv.max_accounts}</span>
                    <span className="font-medium text-muted-foreground">{usagePct}%</span>
                  </div>
                  <div className="w-full bg-muted rounded-full h-1.5 overflow-hidden">
                    <div
                      className={`h-full rounded-full transition-all ${
                        usagePct > 85 ? 'bg-amber-500' : 'bg-primary'
                      }`}
                      style={{ width: `${Math.min(usagePct, 100)}%` }}
                    />
                  </div>
                </div>
              </TableCell>
              <TableCell>
                <Badge variant="success" className="gap-1 text-[11px]">
                  <span className="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-ping" />
                  Online
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                <div className="flex items-center justify-end gap-1.5">
                  <Button
                    size="sm"
                    variant="outline"
                    className="h-7 text-xs gap-1"
                    onClick={() => onTestConnection(srv.id)}
                    disabled={testingId === srv.id}
                  >
                    <Activity className={`h-3 w-3 ${testingId === srv.id ? 'animate-spin' : ''}`} />
                    {testingId === srv.id ? 'Testing...' : 'Test'}
                  </Button>
                  <Button
                    size="icon"
                    variant="ghost"
                    className="h-7 w-7 text-muted-foreground hover:text-destructive"
                    onClick={() => onDelete(srv.id)}
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
};
