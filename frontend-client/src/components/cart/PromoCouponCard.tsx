import React from 'react';
import { Tag, CheckCircle, AlertCircle } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

interface PromoCouponCardProps {
  couponInput: string;
  onCouponInputChange: (val: string) => void;
  onApplyCoupon: (e: React.FormEvent) => void;
  couponSuccess: string | null;
  couponError: string | null;
}

export const PromoCouponCard: React.FC<PromoCouponCardProps> = ({
  couponInput,
  onCouponInputChange,
  onApplyCoupon,
  couponSuccess,
  couponError,
}) => {
  return (
    <Card className="border-border/60 shadow-sm">
      <CardContent className="p-4 space-y-3">
        <form onSubmit={onApplyCoupon} className="flex gap-2">
          <div className="relative flex-1">
            <Tag className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Enter Promo Code (e.g. MERDEKA20)..."
              value={couponInput}
              onChange={(e) => onCouponInputChange(e.target.value)}
              className="pl-9 h-9 uppercase font-mono text-xs"
            />
          </div>
          <Button type="submit" variant="secondary" size="sm" className="font-semibold">
            Apply
          </Button>
        </form>

        {couponSuccess && (
          <p className="text-xs text-emerald-600 dark:text-emerald-400 font-semibold flex items-center gap-1.5">
            <CheckCircle className="h-3.5 w-3.5" />
            {couponSuccess}
          </p>
        )}
        {couponError && (
          <p className="text-xs text-destructive font-semibold flex items-center gap-1.5">
            <AlertCircle className="h-3.5 w-3.5" />
            {couponError}
          </p>
        )}
      </CardContent>
    </Card>
  );
};
