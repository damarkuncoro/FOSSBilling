import React from 'react';
import { Layers, CheckCircle, XCircle } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ExtensionModuleItem } from '@/types/modules';

interface InstalledModulesTabProps {
  extensions: ExtensionModuleItem[];
  onToggle: (id: string, enabled: boolean) => void;
}

export const InstalledModulesTab: React.FC<InstalledModulesTabProps> = ({ extensions, onToggle }) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {extensions.map((ext) => (
        <Card key={ext.id} className="border-border/60 shadow-sm flex flex-col justify-between">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <Badge variant="outline" className="text-[10px] uppercase font-mono">
                {ext.type}
              </Badge>
              <Badge variant={ext.is_enabled ? 'success' : 'secondary'}>
                {ext.is_enabled ? 'Enabled' : 'Disabled'}
              </Badge>
            </div>
            <CardTitle className="text-base mt-2">{ext.name}</CardTitle>
            <CardDescription className="text-xs line-clamp-2">{ext.description}</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-center justify-between text-xs text-muted-foreground pt-2 border-t">
              <span>v{ext.version}</span>
              <span>By {ext.author}</span>
            </div>

            <Button
              variant={ext.is_enabled ? 'outline' : 'default'}
              size="sm"
              className="w-full text-xs font-semibold gap-1.5"
              onClick={() => onToggle(ext.id, ext.is_enabled)}
            >
              {ext.is_enabled ? (
                <>
                  <XCircle className="h-3.5 w-3.5 text-destructive" />
                  Disable Module
                </>
              ) : (
                <>
                  <CheckCircle className="h-3.5 w-3.5" />
                  Enable Module
                </>
              )}
            </Button>
          </CardContent>
        </Card>
      ))}
    </div>
  );
};
