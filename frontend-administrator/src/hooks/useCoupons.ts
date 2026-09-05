import { useState, useEffect } from 'react';
import { api, CouponItem } from '@/lib/api';

export const defaultCoupons: CouponItem[] = [
  {
    id: 1,
    code: 'WELCOME50',
    type: 'percentage',
    value: 50,
    max_uses: 500,
    used_count: 142,
    expires_at: '2026-12-31',
    is_active: true,
  },
  {
    id: 2,
    code: 'HOSTING10OFF',
    type: 'fixed',
    value: 10,
    max_uses: 200,
    used_count: 65,
    expires_at: '2026-10-15',
    is_active: true,
  },
  {
    id: 3,
    code: 'FLASH2026',
    type: 'percentage',
    value: 25,
    max_uses: 100,
    used_count: 100,
    expires_at: '2026-06-01',
    is_active: false,
  },
];

export const initialCouponForm: Partial<CouponItem> = {
  code: '',
  type: 'percentage',
  value: 20,
  max_uses: 100,
  expires_at: '',
  is_active: true,
};

export function useCoupons() {
  const [coupons, setCoupons] = useState<CouponItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [openModal, setOpenModal] = useState(false);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState(initialCouponForm);

  const fetchCoupons = async () => {
    setLoading(true);
    try {
      const data = await api.getCoupons().catch(() => null);
      setCoupons(data && data.length > 0 ? data : defaultCoupons);
    } catch {
      setCoupons(defaultCoupons);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCoupons();
  }, []);

  const generateRandomCode = () => {
    const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789';
    let code = 'PROMO';
    for (let i = 0; i < 5; i++) {
      code += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    setForm((prev) => ({ ...prev, code }));
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      const newCoupon: CouponItem = {
        id: Date.now(),
        code: (form.code || 'COUPON').toUpperCase(),
        type: form.type || 'percentage',
        value: Number(form.value) || 0,
        max_uses: Number(form.max_uses) || 0,
        used_count: 0,
        expires_at: form.expires_at || undefined,
        is_active: form.is_active ?? true,
      };

      await api.createCoupon(newCoupon).catch(() => null);
      setCoupons((prev) => [newCoupon, ...prev]);
      setOpenModal(false);
      setForm(initialCouponForm);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this coupon?')) return;
    try {
      await api.deleteCoupon(id).catch(() => null);
      setCoupons((prev) => prev.filter((c) => c.id !== id));
    } catch (err) {
      console.error(err);
    }
  };

  const toggleStatus = (id: number) => {
    setCoupons((prev) =>
      prev.map((c) => (c.id === id ? { ...c, is_active: !c.is_active } : c))
    );
  };

  return {
    coupons,
    loading,
    openModal,
    setOpenModal,
    saving,
    form,
    setForm,
    fetchCoupons,
    generateRandomCode,
    handleCreate,
    handleDelete,
    toggleStatus,
  };
}
