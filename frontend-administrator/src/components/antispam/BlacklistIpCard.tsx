import React from 'react';
import { Ban, Plus, Trash2 } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

interface BlacklistIpCardProps {
  blacklist: string[];
  newIp: string;
  onNewIpChange: (ip: string) => void;
  onAddIp: () => void;
  onRemoveIp: (ip: string) => void;
}

export const BlacklistIpCard: React.FC<BlacklistIpCardProps> = ({
  blacklist,
  newIp,
  onNewIpChange,
  onAddIp,
  onRemoveIp,
}) => {
  return (
    <Card className="border-border/60 shadow-sm">
      <CardHeader>
        <div className="flex items-center gap-2">
          <Ban className="h-5 w-5 text-destructive" />
          <CardTitle className="text-base font-semibold">IP Address Blacklist</CardTitle>
        </div>
        <CardDescription>
          Explicitly block abusive client IP addresses or fraudulent bot networks
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex gap-2">
          <Input
            placeholder="e.g. 198.51.100.25"
            value={newIp}
            onChange={(e) => onNewIpChange(e.target.value)}
            className="h-9 text-xs font-mono"
          />
          <Button size="sm" onClick={onAddIp} className="gap-1.5 shrink-0">
            <Plus className="h-3.5 w-3.5" />
            Block IP
          </Button>
        </div>

        <div className="border rounded-lg divide-y max-h-60 overflow-y-auto">
          {blacklist.length === 0 ? (
            <p className="p-4 text-center text-xs text-muted-foreground">No IP addresses blacklisted.</p>
          ) : (
            blacklist.map((ip) => (
              <div key={ip} className="p-2.5 flex items-center justify-between text-xs font-mono">
                <span className="text-foreground">{ip}</span>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6 text-muted-foreground hover:text-destructive"
                  onClick={() => onRemoveIp(ip)}
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              </div>
            ))
          )}
        </div>
      </CardContent>
    </Card>
  );
};
