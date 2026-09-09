import React from 'react';
import { ArrowRight, ShieldCheck } from 'lucide-react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';
import { Label } from '@/components/ui/label';

interface OrderSummaryCardProps {
  subtotal: number;
  discount: number;
  promoCode: string | null;
  tax: number;
  total: number;
  checkoutLoading: boolean;
  onCheckout: (gateway: string) => void;
}

export const OrderSummaryCard: React.FC<OrderSummaryCardProps> = ({
  subtotal,
  discount,
  promoCode,
  tax,
  total,
  checkoutLoading,
  onCheckout,
}) => {
  const [gateway, setGateway] = React.useState('midtrans');

  return (
    <Card className="border-border/60 shadow-md">
      <CardHeader>
        <CardTitle className="text-base font-semibold">Payment Summary</CardTitle>
        <CardDescription>Instant automated activation on invoice settlement</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4 text-sm">
        <div className="space-y-3">
          <div className="flex justify-between text-muted-foreground">
            <span>Subtotal</span>
            <span className="font-medium text-foreground">{formatMoney(subtotal)}</span>
          </div>

          {discount > 0 && (
            <div className="flex justify-between text-emerald-600 dark:text-emerald-400 font-semibold">
              <span>Promo Discount ({promoCode})</span>
              <span>-{formatMoney(discount)}</span>
            </div>
          )}

          <div className="flex justify-between text-muted-foreground">
            <span>Tax / VAT (11%)</span>
            <span className="font-medium text-foreground">{formatMoney(tax)}</span>
          </div>

          <div className="pt-3 border-t flex justify-between items-baseline">
            <span className="font-bold text-base">Total Due</span>
            <span className="text-2xl font-extrabold text-primary">{formatMoney(total)}</span>
          </div>
        </div>

        <div className="pt-2 space-y-3">
          <Label className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Select Payment Method</Label>
          <div className="grid grid-cols-1 gap-2">
             <div
                onClick={() => setGateway('midtrans')}
                className={`flex items-center justify-between space-x-2 rounded-lg border p-3 transition-colors cursor-pointer ${gateway === 'midtrans' ? 'border-primary bg-primary/5 ring-1 ring-primary' : 'hover:bg-muted/50'}`}
             >
                <div className="flex items-center gap-3">
                  <div className={`h-4 w-4 rounded-full border border-primary flex items-center justify-center ${gateway === 'midtrans' ? 'bg-primary' : ''}`}>
                    {gateway === 'midtrans' && <div className="h-1.5 w-1.5 rounded-full bg-white" />}
                  </div>
                  <span className="font-semibold">Midtrans / Local VA</span>
                </div>
             </div>
             <div
                onClick={() => setGateway('stripe')}
                className={`flex items-center justify-between space-x-2 rounded-lg border p-3 transition-colors cursor-pointer ${gateway === 'stripe' ? 'border-primary bg-primary/5 ring-1 ring-primary' : 'hover:bg-muted/50'}`}
             >
                <div className="flex items-center gap-3">
                  <div className={`h-4 w-4 rounded-full border border-primary flex items-center justify-center ${gateway === 'stripe' ? 'bg-primary' : ''}`}>
                    {gateway === 'stripe' && <div className="h-1.5 w-1.5 rounded-full bg-white" />}
                  </div>
                  <span className="font-semibold">Stripe / Global Cards</span>
                </div>
             </div>
          </div>
        </div>
      </CardContent>
      <CardFooter className="flex flex-col gap-3">
        <Button
          className="w-full gap-2 font-semibold shadow-md shadow-primary/25 h-11 text-base"
          onClick={() => onCheckout(gateway)}
          disabled={checkoutLoading}
        >
          {checkoutLoading ? 'Processing Checkout...' : 'Secure Checkout Now'}
          {!checkoutLoading && <ArrowRight className="h-4 w-4" />}
        </Button>

        <div className="flex items-center justify-center gap-1.5 text-xs text-muted-foreground">
          <ShieldCheck className="h-4 w-4 text-emerald-500" />
          <span>256-Bit SSL Encrypted & Secure Checkout</span>
        </div>
      </CardFooter>
    </Card>
  );
};
