import React from 'react';
import { Link } from 'react-router-dom';
import { ShoppingCart, Zap } from 'lucide-react';
import { useClientCartPage } from '@/hooks/useClientCartPage';
import { Button } from '@/components/ui/button';
import { CartItemsList } from '@/components/cart/CartItemsList';
import { PromoCouponCard } from '@/components/cart/PromoCouponCard';
import { OrderSummaryCard } from '@/components/cart/OrderSummaryCard';

export const Cart: React.FC = () => {
  const {
    items,
    removeItem,
    clearCart,
    promoCode,
    subtotal,
    discount,
    tax,
    total,
    couponInput,
    setCouponInput,
    couponError,
    couponSuccess,
    checkoutLoading,
    handleApplyCoupon,
    handleCheckout,
  } = useClientCartPage();

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
        <p className="text-sm text-muted-foreground">
          Review your cloud services and apply promotional voucher codes.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
        <div className="lg:col-span-2 space-y-4">
          <CartItemsList
            items={items}
            onRemoveItem={removeItem}
            onClearCart={clearCart}
          />
          <PromoCouponCard
            couponInput={couponInput}
            onCouponInputChange={setCouponInput}
            onApplyCoupon={handleApplyCoupon}
            couponSuccess={couponSuccess}
            couponError={couponError}
          />
        </div>

        <OrderSummaryCard
          subtotal={subtotal}
          discount={discount}
          promoCode={promoCode}
          tax={tax}
          total={total}
          checkoutLoading={checkoutLoading}
          onCheckout={handleCheckout}
        />
      </div>
    </div>
  );
};
