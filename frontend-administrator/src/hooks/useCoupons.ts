import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminCouponService } from '@/services/admin_coupon.service';
import type { CouponItem } from '@/types/api';

export const defaultCoupons: CouponItem[] = [
  { id: 1, code: 'WELCOME50', type: 'percentage', value: 50, max_uses: 500, used_count: 142, expires_at: '2026-12-31', is_active: true },
  { id: 2, code: 'HOSTING10OFF', type: 'fixed', value: 10, max_uses: 200, used_count: 65, expires_at: '2026-10-15', is_active: true },
];

export const initialCouponForm: Partial<CouponItem> = { code: '', type: 'percentage', value: 20, max_uses: 100, expires_at: '', is_active: true };

export function useCoupons() {
  const queryClient = useQueryClient();
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState(initialCouponForm);

  const { data: coupons = [], isLoading: loading, refetch: fetchCoupons } = useQuery({
    queryKey: ['admin', 'coupons'],
    queryFn: async () => {
      const data = await adminCouponService.listCoupons();
      return data && data.length > 0 ? data : defaultCoupons;
    },
  });

  const createMutation = useMutation({
    mutationFn: (input: Partial<CouponItem>) => adminCouponService.createCoupon(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'coupons'] });
      setOpenModal(false);
      setForm(initialCouponForm);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminCouponService.deleteCoupon(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'coupons'] }),
  });

  return {
    coupons,
    loading,
    fetchCoupons,
    openModal,
    setOpenModal,
    saving: createMutation.isPending || deleteMutation.isPending,
    form,
    setForm,
    generateRandomCode: () => setForm((p) => ({ ...p, code: adminCouponService.generateRandomCode() })),
    handleCreate: (e: React.FormEvent) => { e.preventDefault(); createMutation.mutate(form); },
    handleDelete: (id: number) => { if (confirm('Delete coupon?')) deleteMutation.mutate(id); },
    toggleStatus: (id: number) => { console.log('Toggle status for', id); },
  };
}
