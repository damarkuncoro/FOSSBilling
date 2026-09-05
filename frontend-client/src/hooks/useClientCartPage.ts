import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '@/lib/api';
import { useCart } from '@/lib/cart';
import { useClientAuth } from '@/lib/auth';

export function useClientCartPage() {
  const {
    items,
    removeItem,
    clearCart,
    promoCode,
    applyPromo,
    subtotal,
    discount,
    tax,
    total,
  } = useCart();
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
      await api.checkoutCart({
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

  return {
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
  };
}
