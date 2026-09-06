import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '@/lib/cart';
import { domainService } from '@/services/domain.service';
import { newsService } from '@/services/news.service';
import { DomainSearchResult, HostingPlan } from '@/types/api';

export const defaultHostingPlans: HostingPlan[] = [
  {
    id: 1,
    title: 'cPanel Starter Cloud',
    description: 'Perfect for personal blogs, portfolios, and small business sites.',
    price: 9.99,
    period: '1M',
    type: 'hosting',
    features: ['1 Website Domain', '10 GB NVMe Storage', 'Unlimited Bandwidth', 'Free SSL Certificate', 'cPanel Control Panel'],
    popular: false,
  },
  {
    id: 2,
    title: 'Cloud VPS Pro (DirectAdmin)',
    description: 'High performance dedicated computing for high-traffic SaaS apps.',
    price: 29.99,
    period: '1M',
    type: 'vps',
    features: ['4 vCPU Cores', '8 GB RAM Memory', '80 GB NVMe Storage', 'Dedicated IPv4 Address', 'Root SSH & DirectAdmin'],
    popular: true,
  },
  {
    id: 3,
    title: 'FOSSBilling Enterprise License',
    description: 'Self-hosted enterprise license with priority 24/7 SLA.',
    price: 199.0,
    period: '1Y',
    type: 'license',
    features: ['Unlimited Clients', 'Multi-Server Clusters', 'Whitelabel Branding', 'Priority Ticket SLA', 'Automated Upgrades'],
    popular: false,
  },
  {
    id: 4,
    title: 'Nusantara Cloud OS Template',
    description: 'Pre-hardened Linux image with automated docker deployments.',
    price: 15.0,
    period: 'ONETIME',
    type: 'downloadable',
    features: ['Instant Download Access', 'HMAC Signed URLs', 'Docker & Kubernetes Ready', 'Security CIS Benchmark', 'Free Updates'],
    popular: false,
  },
];

export function useStorefront() {
  const { addItem } = useCart();
  const navigate = useNavigate();
  const [domainSearch, setDomainSearch] = useState('');
  const [domainResult, setDomainResult] = useState<DomainSearchResult | null>(null);
  const [searching, setSearching] = useState(false);
  const [news, setNews] = useState<any[]>([]);

  useEffect(() => {
    newsService.listPublishedNews().then((data) => setNews(data || [])).catch(() => {});
  }, []);

  const handleDomainCheck = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!domainSearch.trim()) return;

    setSearching(true);
    try {
      const res = await domainService.checkAvailability(domainSearch);
      setDomainResult({
        domain: res.domain,
        available: res.available,
        price: res.price,
        currency: res.currency,
      });
    } catch (err) {
      console.error('Failed to check domain availability:', err);
      setDomainResult(null);
    } finally {
      setSearching(false);
    }
  };

  const handleAddToCart = (plan: HostingPlan) => {
    addItem({
      id: `${plan.id}-${Date.now()}`,
      product_id: plan.id,
      title: plan.title,
      price: plan.price,
      period: plan.period,
      type: plan.type,
    });
    navigate('/cart');
  };

  const handleAddDomainToCart = () => {
    if (!domainResult) return;
    addItem({
      id: `domain-${Date.now()}`,
      product_id: 10,
      title: `Domain Registration: ${domainResult.domain}`,
      price: domainResult.price,
      period: '1Y',
      type: 'domain',
      domain_name: domainResult.domain,
    });
    navigate('/cart');
  };

  return {
    domainSearch,
    setDomainSearch,
    domainResult,
    searching,
    news,
    handleDomainCheck,
    handleAddToCart,
    handleAddDomainToCart,
    plans: defaultHostingPlans,
  };
}
