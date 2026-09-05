import React from 'react';
import { ShieldCheck, Server } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { RegistrarConfig } from '@/types/modules';

interface RegistrarsConfigTabProps {
  registrars: RegistrarConfig[];
}

export const RegistrarsConfigTab: React.FC<RegistrarsConfigTabProps> = ({ registrars }) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      {registrars.map((reg) => (
        <Card key={reg.id} className="border-border/60 shadow-sm">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2.5">
                <Server className="h-5 w-5 text-primary" />
                <CardTitle className="text-base">{reg.name}</CardTitle>
              </div>
              <Badge variant={reg.enabled ? 'success' : 'secondary'}>
                {reg.enabled ? 'Enabled' : 'Disabled'}
              </Badge>
            </div>
            <CardDescription className="text-xs">
              Automated API provisioning for domain registrations and renewals
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="p-3 rounded-lg bg-muted/40 border text-xs space-y-1 font-mono">
              <div className="flex justify-between">
                <span className="text-muted-foreground">API User:</span>
                <span className="font-semibold">{reg.api_user || 'Not configured'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Mode:</span>
                <span className="font-semibold">{reg.test_mode ? 'Sandbox / Test' : 'Live Production'}</span>
              </div>
            </div>

            <Button
              variant="outline"
              size="sm"
              className="w-full text-xs font-semibold"
              onClick={() => alert(`Opening configuration settings for ${reg.name}...`)}
            >
              Configure Credentials
            </Button>
          </CardContent>
        </Card>
      ))}
    </div>
  );
};
