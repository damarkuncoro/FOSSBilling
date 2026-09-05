import React from 'react';
import { Link } from 'react-router-dom';
import { Server, ArrowRight } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatMoney } from '@/lib/utils';
import { Order } from '@/types/api';

interface ActiveServicesListProps {
  orders: Order[];
}

export const ActiveServicesList: React.FC<ActiveServicesListProps> = ({ orders }) => {
  return (
    <Card className="border-border/60">
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle className="text-base font-semibold">My Active Hosting & Services</CardTitle>
          <CardDescription>Quick view of your deployed cloud servers</CardDescription>
        </div>
        <Link to="/services">
          <Button variant="ghost" size="sm" className="text-xs gap-1">
            View All <ArrowRight className="h-3 w-3" />
          </Button>
        </Link>
      </CardHeader>
      <CardContent className="space-y-3">
        {orders.length === 0 ? (
          <p className="text-sm text-muted-foreground py-4 text-center">No active services yet.</p>
        ) : (
          orders.slice(0, 3).map((order) => (
            <div key={order.id} className="p-3 rounded-lg border bg-muted/20 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Server className="h-5 w-5 text-primary shrink-0" />
                <div>
                  <p className="font-semibold text-sm">{order.title || `Service #${order.id}`}</p>
                  <p className="text-xs text-muted-foreground">
                    Period: {order.period} • {formatMoney(order.price, order.currency)}
                  </p>
                </div>
              </div>
              <Badge variant={order.status === 'active' ? 'success' : 'warning'}>
                {order.status}
              </Badge>
            </div>
          ))
        )}
      </CardContent>
    </Card>
  );
};
