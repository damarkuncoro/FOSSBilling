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
    updateItemConfig,
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

  const handleCheckout = async (gatewayID: string = 'midtrans') => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    setCheckoutLoading(true);
    try {
      const res = await api.checkoutCart({
        client_id: user?.id || 1,
        items: items.map((i) => ({
          product_id: i.product_id,
          title: i.title,
          period: i.period,
          price: i.price,
          quantity: 1,
          config: i.config || (i.domain_name ? { domain_name: i.domain_name } : undefined),
        })),
        promo_code: promoCode || undefined,
      });

      clearCart();

      // If checkout successful, initiate payment for the generated invoice
      if (res && res.invoice) {
        try {
          const payRes = await api.initiateInvoicePayment(res.invoice.id, gatewayID);
          if (payRes.redirect_url) {
            window.location.href = payRes.redirect_url;
            return;
          }
        } catch (payErr) {
           console.error('Payment initiation failed, falling back to invoices list', payErr);
        }
      }

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
    updateItemConfig,
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
