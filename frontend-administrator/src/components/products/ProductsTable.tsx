import React from 'react';
import {
  RefreshCw,
  Server,
  Globe,
  KeyRound,
  Download,
  Boxes,
  Trash2,
  CheckCircle2,
  XCircle,
} from 'lucide-react';
import { ProductItem } from '@/lib/api';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

interface ProductsTableProps {
  products: ProductItem[];
  loading: boolean;
  onToggleStatus: (id: number) => void;
  onDelete: (id: number) => void;
}

export const ProductsTable: React.FC<ProductsTableProps> = ({
  products,
  loading,
  onToggleStatus,
  onDelete,
}) => {
  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'hosting':
        return <Server className="h-4 w-4 text-sky-500" />;
      case 'domain':
        return <Globe className="h-4 w-4 text-emerald-500" />;
      case 'license':
        return <KeyRound className="h-4 w-4 text-amber-500" />;
      case 'downloadable':
        return <Download className="h-4 w-4 text-purple-500" />;
      default:
        return <Boxes className="h-4 w-4 text-muted-foreground" />;
    }
  };

  if (loading) {
    return (
      <div className="text-center py-10 text-muted-foreground">
        <RefreshCw className="h-5 w-5 animate-spin mx-auto mb-2 text-primary" />
        Loading product catalog...
      </div>
    );
  }

  if (products.length === 0) {
    return <div className="text-center py-10 text-muted-foreground">No products found.</div>;
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Product Name</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Category</TableHead>
          <TableHead>Monthly</TableHead>
          <TableHead>Annually</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {products.map((prod) => (
          <TableRow key={prod.id} className="hover:bg-muted/40 transition-colors">
            <TableCell>
              <div>
                <div className="font-semibold text-sm flex items-center gap-1.5">{prod.title}</div>
                <span className="font-mono text-xs text-muted-foreground">/{prod.slug}</span>
              </div>
            </TableCell>
            <TableCell>
              <div className="flex items-center gap-1.5 text-xs font-medium capitalize">
                {getTypeIcon(prod.type)}
                {prod.type}
              </div>
            </TableCell>
            <TableCell className="text-xs text-muted-foreground">{prod.category_name || 'General'}</TableCell>
            <TableCell className="font-mono text-xs font-semibold">
              {prod.price_monthly > 0 ? `$${prod.price_monthly.toFixed(2)}/mo` : '—'}
            </TableCell>
            <TableCell className="font-mono text-xs font-semibold">
              {prod.price_annually && prod.price_annually > 0 ? `$${prod.price_annually.toFixed(2)}/yr` : '—'}
            </TableCell>
            <TableCell>
              <button onClick={() => onToggleStatus(prod.id)} className="cursor-pointer focus:outline-none">
                {prod.is_active ? (
                  <Badge variant="success" className="gap-1 text-[11px]">
                    <CheckCircle2 className="h-3 w-3" />
                    Active
                  </Badge>
                ) : (
                  <Badge variant="outline" className="gap-1 text-[11px]">
                    <XCircle className="h-3 w-3" />
                    Disabled
                  </Badge>
                )}
              </button>
            </TableCell>
            <TableCell className="text-right">
              <Button
                size="icon"
                variant="ghost"
                className="h-7 w-7 text-muted-foreground hover:text-destructive"
                onClick={() => onDelete(prod.id)}
              >
                <Trash2 className="h-3.5 w-3.5" />
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};
