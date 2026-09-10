import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { useCart } from '@/lib/cart';
import { useClientAuth } from '@/lib/auth';

export function useClientCartPage() {
  const { items, removeItem, clearCart, updateItemConfig, promoCode, applyPromo, subtotal, discount, tax, total } = useCart();
  const { user, isAuthenticated } = useClientAuth();
  const navigate = useNavigate();
  const [couponInput, setCouponInput] = useState(promoCode);
  const [couponError, setCouponError] = useState<string | null>(null);
  const [couponSuccess, setCouponSuccess] = useState<string | null>(null);

  const checkoutMutation = useMutation({
    mutationFn: (gid: string) => api.checkoutCart({ client_id: user?.id || 1, items: items.map(i => ({ product_id: i.product_id, title: i.title, period: i.period, price: i.price, quantity: 1, config: i.config || (i.domain_name ? { domain_name: i.domain_name } : undefined) })), promo_code: promoCode || undefined }),
    onSuccess: async (res, gid) => {
      clearCart();
      if (res?.invoice) {
        try { const p = await api.initiateInvoicePayment(res.invoice.id, gid); if (p.redirect_url) { window.location.href = p.redirect_url; return; } } catch {}
      }
      navigate('/invoices');
    },
  });

  return {
    items, removeItem, clearCart, updateItemConfig, promoCode, subtotal, discount, tax, total, couponInput, setCouponInput, couponError, couponSuccess,
    checkoutLoading: checkoutMutation.isPending,
    handleApplyCoupon: async (e: any) => {
      e.preventDefault(); setCouponError(null); setCouponSuccess(null);
      try { await applyPromo(couponInput.toUpperCase()); setCouponSuccess('Coupon applied!'); } catch (err: any) { setCouponError(err.message); }
    },
    handleCheckout: (gid = 'midtrans') => isAuthenticated ? checkoutMutation.mutate(gid) : navigate('/login'),
  };
}
