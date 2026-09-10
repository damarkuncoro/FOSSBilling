import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation } from '@tanstack/react-query';
import { useCart } from '@/lib/cart';
import { domainService } from '@/services/domain.service';
import { newsService } from '@/services/news.service';
import { catalogService } from '@/services/catalog.service';
import { DomainSearchResult, HostingPlan } from '@/types/api';

export const defaultHostingPlans: HostingPlan[] = [
  { id: 1, title: 'cPanel Starter Cloud', description: 'Perfect for small sites.', price: 9.99, period: '1M', type: 'hosting', features: ['Free SSL'], popular: false },
  { id: 2, title: 'Cloud VPS Pro', description: 'High performance.', price: 29.99, period: '1M', type: 'vps', features: ['Root SSH'], popular: true },
];

export function useStorefront() {
  const { addItem } = useCart();
  const navigate = useNavigate();
  const [domainSearch, setDomainSearch] = useState('');

  const { data: news = [] } = useQuery({
    queryKey: ['public', 'news'],
    queryFn: () => newsService.listPublishedNews(),
  });

  const { data: plans = defaultHostingPlans } = useQuery({
    queryKey: ['public', 'products'],
    queryFn: async () => {
      const data = await catalogService.listProducts();
      return data.length > 0 ? data : defaultHostingPlans;
    },
  });

  const domainMutation = useMutation({
    mutationFn: (domain: string) => domainService.checkAvailability(domain),
  });

  const handleDomainCheck = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!domainSearch.trim()) return;
    domainMutation.mutate(domainSearch);
  };

  const handleAddToCart = (plan: HostingPlan) => {
    addItem({ id: `${plan.id}-${Date.now()}`, product_id: plan.id, title: plan.title, price: plan.price, period: plan.period, type: plan.type, form_id: plan.form_id });
    navigate('/cart');
  };

  const handleAddDomainToCart = () => {
    const res = domainMutation.data;
    if (!res) return;
    addItem({ id: `domain-${Date.now()}`, product_id: 10, title: `Domain: ${res.domain}`, price: res.price, period: '1Y', type: 'domain', domain_name: res.domain });
    navigate('/cart');
  };

  return {
    domainSearch,
    setDomainSearch,
    domainResult: domainMutation.data as DomainSearchResult,
    searching: domainMutation.isPending,
    news,
    handleDomainCheck,
    handleAddToCart,
    handleAddDomainToCart,
    plans,
  };
}
