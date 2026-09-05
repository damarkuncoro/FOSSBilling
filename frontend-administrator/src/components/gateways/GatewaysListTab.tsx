import React from 'react';
import { CreditCard, CheckCircle2, XCircle, Sliders, ShieldCheck, Key } from 'lucide-react';
import { PaymentGatewayItem } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

interface GatewaysListTabProps {
  gateways: PaymentGatewayItem[];
  onToggleEnabled: (id: string) => void;
  onConfigure: (gw: PaymentGatewayItem) => void;
}

export const GatewaysListTab: React.FC<GatewaysListTabProps> = ({
  gateways,
  onToggleEnabled,
  onConfigure,
}) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {gateways.map((gw) => (
        <Card key={gw.id} className="border-border/60 hover:shadow-md transition-shadow">
          <CardHeader className="pb-3">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-xl bg-primary/10 text-primary flex items-center justify-center">
                  <CreditCard className="h-5 w-5" />
                </div>
                <div>
                  <CardTitle className="text-base font-semibold">{gw.name}</CardTitle>
                  <CardDescription className="text-xs">{gw.description}</CardDescription>
                </div>
              </div>
              <button
                onClick={() => onToggleEnabled(gw.id)}
                className="cursor-pointer focus:outline-none"
                title="Toggle Gateway Status"
              >
                {gw.enabled ? (
                  <Badge variant="success" className="gap-1 text-[11px]">
                    <CheckCircle2 className="h-3 w-3" />
                    Enabled
                  </Badge>
                ) : (
                  <Badge variant="outline" className="gap-1 text-[11px]">
                    <XCircle className="h-3 w-3" />
                    Disabled
                  </Badge>
                )}
              </button>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between text-xs border-t pt-3">
              <span className="text-muted-foreground flex items-center gap-1.5">
                <ShieldCheck className="h-3.5 w-3.5" />
                Mode
              </span>
              {gw.test_mode ? (
                <Badge variant="secondary" className="text-[10px] uppercase font-mono">
                  Sandbox / Test Mode
                </Badge>
              ) : (
                <Badge variant="outline" className="text-[10px] uppercase font-mono border-emerald-500 text-emerald-600">
                  Live Production
                </Badge>
              )}
            </div>

            {gw.public_key && (
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground flex items-center gap-1.5">
                  <Key className="h-3.5 w-3.5" />
                  Client Key
                </span>
                <span className="font-mono text-[11px] text-muted-foreground truncate max-w-[180px]">
                  {gw.public_key.slice(0, 14)}...
                </span>
              </div>
            )}

            <div className="flex justify-end gap-2 pt-2 border-t">
              <Button
                variant="outline"
                size="sm"
                className="w-full gap-1.5 text-xs"
                onClick={() => onConfigure(gw)}
              >
                <Sliders className="h-3.5 w-3.5" />
                Configure Credentials
              </Button>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
};
