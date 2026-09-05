import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  ShoppingCart,
  Trash2,
  Tag,
  ArrowRight,
  ShieldCheck,
  Zap,
  CheckCircle,
  AlertCircle,
} from 'lucide-react';
import { api } from '@/lib/api';
import { useCart } from '@/lib/cart';
import { useClientAuth } from '@/lib/auth';
import { formatMoney } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';

export const Cart: React.FC = () => {
  const { items, removeItem, clearCart, promoCode, setPromoCode, applyPromo, subtotal, discount, tax, total } = useCart();
  const { user, isAuthenticated } = useClientAuth();
  const navigate = useNavigate();
  const [couponInput, setCouponInput] = useState(promoCode);
  const [couponError, setCouponError] = useState<string | null>(null);
  const [couponSuccess, setCouponSuccess] = useState<string | null>(null);
  const [checkoutLoading, setCheckoutLoading] = useState(false);

  const handleApplyCoupon = async (e: React.FormEvent) => {
    e.preventDefault();
    setCouponError(null);
    setCouponSuccess(null);
    if (!couponInput.trim()) return;

    try {
      await applyPromo(couponInput.trim().toUpperCase());
      setCouponSuccess(`Coupon "${couponInput.toUpperCase()}" applied successfully!`);
    } catch (err: any) {
      setCouponError(err.message || 'Invalid coupon code');
    }
  };

  const handleCheckout = async () => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    setCheckoutLoading(true);
    try {
      const res = await api.checkoutCart({
        client_id: user?.id || 1,
        items: items.map((i) => ({ product_id: i.product_id, period: i.period })),
        promo_code: promoCode || undefined,
        gateway: 'midtrans',
      });
      clearCart();
      navigate('/invoices');
    } catch (err: any) {
      alert(`Checkout failed: ${err.message}`);
    } finally {
      setCheckoutLoading(false);
    }
  };

  if (items.length === 0) {
    return (
      <div className="min-h-[50vh] flex flex-col items-center justify-center text-center space-y-4 py-16 animate-in fade-in">
        <div className="h-16 w-16 rounded-full bg-muted flex items-center justify-center text-muted-foreground">
          <ShoppingCart className="h-8 w-8" />
        </div>
        <h2 className="text-2xl font-bold">Your shopping cart is empty</h2>
        <p className="text-sm text-muted-foreground max-w-sm">
          Explore our cloud hosting packages, domain registrations, or digital downloads.
        </p>
        <Link to="/">
          <Button className="gap-2 font-semibold">
            <Zap className="h-4 w-4" />
            Browse Store Catalog
          </Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-in fade-in-50 duration-300">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Shopping Cart & Order Summary</h1>
        <p className="text-sm text-muted-foreground">Review your cloud services and apply promotional voucher codes.</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
        {/* Cart Items List */}
        <div className="lg:col-span-2 space-y-4">
          <Card className="border-border/60 shadow-sm">
            <CardHeader className="flex flex-row items-center justify-between pb-3">
              <CardTitle className="text-base font-semibold">Selected Services ({items.length})</CardTitle>
              <Button variant="ghost" size="sm" onClick={clearCart} className="text-xs text-destructive hover:bg-destructive/10">
                Clear Cart
              </Button>
            </CardHeader>
            <CardContent className="divide-y">
              {items.map((item) => (
                <div key={item.id} className="py-4 first:pt-0 last:pb-0 flex items-center justify-between gap-4">
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
                      onClick={() => removeItem(item.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </CardContent>
          </Card>

          {/* Promo Coupon Box */}
          <Card className="border-border/60 shadow-sm">
            <CardContent className="p-4 space-y-3">
              <form onSubmit={handleApplyCoupon} className="flex gap-2">
                <div className="relative flex-1">
                  <Tag className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    placeholder="Enter Promo Code (e.g. MERDEKA20)..."
                    value={couponInput}
                    onChange={(e) => setCouponInput(e.target.value)}
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
        </div>

        {/* Order Breakdown / Totals */}
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
              onClick={handleCheckout}
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
      </div>
    </div>
  );
};
