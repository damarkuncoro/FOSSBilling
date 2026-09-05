import React from 'react';
import { Trash2 } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatMoney } from '@/lib/utils';
import { CartItem } from '@/lib/cart';

interface CartItemsListProps {
  items: CartItem[];
  onRemoveItem: (id: string) => void;
  onClearCart: () => void;
}

export const CartItemsList: React.FC<CartItemsListProps> = ({
  items,
  onRemoveItem,
  onClearCart,
}) => {
  return (
    <Card className="border-border/60 shadow-sm">
      <CardHeader className="flex flex-row items-center justify-between pb-3">
        <CardTitle className="text-base font-semibold">
          Selected Services ({items.length})
        </CardTitle>
        <Button
          variant="ghost"
          size="sm"
          onClick={onClearCart}
          className="text-xs text-destructive hover:bg-destructive/10"
        >
          Clear Cart
        </Button>
      </CardHeader>
      <CardContent className="divide-y">
        {items.map((item) => (
          <div
            key={item.id}
            className="py-4 first:pt-0 last:pb-0 flex items-center justify-between gap-4"
          >
            <div className="space-y-1">
              <p className="font-semibold text-sm">{item.title}</p>
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <Badge variant="outline" className="text-[10px] uppercase font-mono">
                  {item.period}
                </Badge>
                <span>Type: {item.type}</span>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <span className="font-bold text-base">{formatMoney(item.price)}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 text-muted-foreground hover:text-destructive"
                onClick={() => onRemoveItem(item.id)}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
};
