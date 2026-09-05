import React from 'react';
import { Globe, Trash2 } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatMoney } from '@/lib/utils';
import { TldPricingItem } from '@/types/modules';

interface TldPricingTableProps {
  tlds: TldPricingItem[];
  loading: boolean;
  onDelete: (id: number) => void;
}

export const TldPricingTable: React.FC<TldPricingTableProps> = ({ tlds, loading, onDelete }) => {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Extension (TLD)</TableHead>
          <TableHead>Registrar Provider</TableHead>
          <TableHead>Register</TableHead>
          <TableHead>Renew</TableHead>
          <TableHead>Transfer</TableHead>
          <TableHead>Min. Period</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {loading ? (
          <TableRow>
            <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">
              Loading TLD pricing matrix...
            </TableCell>
          </TableRow>
        ) : tlds.length === 0 ? (
          <TableRow>
            <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">
              No TLDs configured yet.
            </TableCell>
          </TableRow>
        ) : (
          tlds.map((tld) => (
            <TableRow key={tld.id}>
              <TableCell className="font-bold text-sm flex items-center gap-2">
                <Globe className="h-4 w-4 text-primary" />
                {tld.tld}
              </TableCell>
              <TableCell className="font-mono text-xs capitalize">{tld.registrar}</TableCell>
              <TableCell className="font-semibold text-xs">{formatMoney(tld.price_registration)}/yr</TableCell>
              <TableCell className="text-xs text-muted-foreground">{formatMoney(tld.price_renewal)}/yr</TableCell>
              <TableCell className="text-xs text-muted-foreground">{formatMoney(tld.price_transfer)}</TableCell>
              <TableCell className="text-xs">{tld.min_years} Year(s)</TableCell>
              <TableCell>
                <Badge variant={tld.is_active ? 'success' : 'secondary'}>
                  {tld.is_active ? 'Active' : 'Disabled'}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7 text-muted-foreground hover:text-destructive"
                  onClick={() => onDelete(tld.id)}
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              </TableCell>
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  );
};
