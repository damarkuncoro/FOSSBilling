import React from 'react';
import { Trash2 } from 'lucide-react';
import { TaxRuleItem } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

interface TaxRulesListTabProps {
  taxRules: TaxRuleItem[];
  onDelete: (id: number) => void;
}

export const TaxRulesListTab: React.FC<TaxRulesListTabProps> = ({ taxRules, onDelete }) => {
  return (
    <Card className="border-border/60 shadow-sm">
      <CardHeader>
        <CardTitle className="text-base font-semibold">Regional Tax & VAT Rules</CardTitle>
        <CardDescription>
          Automatic tax percentage calculated at checkout based on client billing country.
        </CardDescription>
      </CardHeader>
      <CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Tax Rule Name</TableHead>
              <TableHead>Country Code</TableHead>
              <TableHead>Tax Rate</TableHead>
              <TableHead>Scope</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {taxRules.map((rule) => (
              <TableRow key={rule.id}>
                <TableCell className="font-medium text-sm">{rule.name}</TableCell>
                <TableCell>
                  <Badge variant="outline" className="font-mono text-xs">{rule.country}</Badge>
                </TableCell>
                <TableCell className="font-mono text-sm font-bold text-primary">{rule.rate}%</TableCell>
                <TableCell className="text-xs text-muted-foreground">
                  {rule.apply_to_all_clients ? 'All Clients Globally' : 'Matching Country Only'}
                </TableCell>
                <TableCell>
                  {rule.is_active ? (
                    <Badge variant="success" className="text-[10px]">Active</Badge>
                  ) : (
                    <Badge variant="outline" className="text-[10px]">Inactive</Badge>
                  )}
                </TableCell>
                <TableCell className="text-right">
                  <Button
                    size="icon"
                    variant="ghost"
                    className="h-7 w-7 text-muted-foreground hover:text-destructive"
                    onClick={() => onDelete(rule.id)}
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};
