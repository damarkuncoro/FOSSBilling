import React from 'react';
import { ArrowRight, ShieldCheck } from 'lucide-react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';

interface OrderSummaryCardProps {
  subtotal: number;
  discount: number;
  promoCode: string | null;
  tax: number;
  total: number;
  checkoutLoading: boolean;
  onCheckout: () => void;
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
  return (
    <Card className="border-border/60 shadow-md">
      <CardHeader>
        <CardTitle className="text-base font-semibold">Payment Summary</CardTitle>
        <CardDescription>Instant automated activation on invoice settlement</CardDescription>
      </CardHeader>
      <CardContent className="space-y-3 text-sm">
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
      </CardContent>
      <CardFooter className="flex flex-col gap-3">
        <Button
          className="w-full gap-2 font-semibold shadow-md shadow-primary/25 h-11 text-base"
          onClick={onCheckout}
          disabled={checkoutLoading}
        >
          {checkoutLoading ? 'Generating Invoice...' : 'Proceed to Checkout'}
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
