import React from 'react';
import { Key, Plus, Trash2 } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { ApiKey } from '@/types/api';

interface ApiKeysCardProps {
  apiKeys: ApiKey[];
  keyName: string;
  onKeyNameChange: (val: string) => void;
  onGenerateKey: (e: React.FormEvent) => void;
  onRevokeKey: (id: number) => void;
  generating: boolean;
}

export const ApiKeysCard: React.FC<ApiKeysCardProps> = ({
  apiKeys,
  keyName,
  onKeyNameChange,
  onGenerateKey,
  onRevokeKey,
  generating,
}) => {
  return (
    <Card className="border-border/60 shadow-sm">
      <CardHeader>
        <div className="flex items-center gap-2">
          <Key className="h-5 w-5 text-primary" />
          <CardTitle className="text-base font-semibold">Developer API Keys</CardTitle>
        </div>
        <CardDescription>
          Programmatic keys for automated provisioning, CI/CD, and CLI integrations
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <form onSubmit={onGenerateKey} className="flex gap-2">
          <Input
            required
            placeholder="Key label (e.g. Production Terraform Runner)..."
            value={keyName}
            onChange={(e) => onKeyNameChange(e.target.value)}
            className="h-9 text-xs"
          />
          <Button type="submit" size="sm" disabled={generating} className="gap-1.5 shrink-0">
            <Plus className="h-3.5 w-3.5" />
            Generate Key
          </Button>
        </form>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Label</TableHead>
              <TableHead>API Key</TableHead>
              <TableHead className="text-right">Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {apiKeys.length === 0 ? (
              <TableRow>
                <TableCell colSpan={3} className="text-center text-muted-foreground py-4 text-xs">
                  No API keys generated.
                </TableCell>
              </TableRow>
            ) : (
              apiKeys.map((k) => (
                <TableRow key={k.id}>
                  <TableCell className="font-medium text-xs">{k.name}</TableCell>
                  <TableCell className="font-mono text-xs text-primary">
                    {k.key || `fb_key_${k.id}`}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-7 w-7 text-muted-foreground hover:text-destructive"
                      onClick={() => onRevokeKey(k.id)}
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};
