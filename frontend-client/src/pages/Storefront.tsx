import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Server,
  Cpu,
  Globe,
  Download,
  Key,
  Search,
  Check,
  Zap,
  ShieldCheck,
  ArrowRight,
  Sparkles,
} from 'lucide-react';
import { api } from '@/lib/api';
import { useCart } from '@/lib/cart';
import { formatMoney } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';

const hostingPlans = [
  {
    id: 1,
    title: 'cPanel Starter Cloud',
    description: 'Perfect for personal blogs, portfolios, and small business sites.',
    price: 9.99,
    period: '1M',
    type: 'hosting',
    icon: Server,
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
    icon: Cpu,
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
    icon: Key,
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
    icon: Download,
    features: ['Instant Download Access', 'HMAC Signed URLs', 'Docker & Kubernetes Ready', 'Security CIS Benchmark', 'Free Updates'],
    popular: false,
  },
];

export const Storefront: React.FC = () => {
  const { addItem } = useCart();
  const navigate = useNavigate();
  const [domainSearch, setDomainSearch] = useState('');
  const [domainResult, setDomainResult] = useState<any | null>(null);
  const [searching, setSearching] = useState(false);
  const [news, setNews] = useState<any[]>([]);

  useEffect(() => {
    api.guestNews().then((data) => setNews(data || [])).catch(() => {});
  }, []);

  const handleDomainCheck = (e: React.FormEvent) => {
    e.preventDefault();
    if (!domainSearch.trim()) return;

    setSearching(true);
    setTimeout(() => {
      let clean = domainSearch.toLowerCase().replace(/https?:\/\//, '').replace(/\/$/, '');
      if (!clean.includes('.')) clean += '.com';

      setDomainResult({
        domain: clean,
        available: true,
        price: 12.99,
        currency: 'USD',
      });
      setSearching(false);
    }, 400);
  };

  const handleAddToCart = (plan: any) => {
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

  return (
    <div className="space-y-16">
      {/* Hero Section */}
      <section className="text-center py-12 md:py-16 space-y-6 relative overflow-hidden">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-xs font-semibold">
          <Sparkles className="h-3.5 w-3.5" />
          <span>Next-Generation Cloud Infrastructure Powered by Go</span>
        </div>

        <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight max-w-3xl mx-auto leading-tight">
          Superfast Cloud Hosting & Digital Services
        </h1>
        <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
          Deploy high-speed cPanel & DirectAdmin hosting, VPS cloud instances, and enterprise software licenses with instant automated provisioning.
        </p>

        {/* Domain Search Card */}
        <div className="max-w-2xl mx-auto pt-4">
          <Card className="p-2 border-border/80 shadow-xl bg-card/80 backdrop-blur-xl">
            <form onSubmit={handleDomainCheck} className="flex gap-2">
              <div className="relative flex-1">
                <Search className="absolute left-3.5 top-3.5 h-4 w-4 text-muted-foreground" />
                <Input
                  type="text"
                  placeholder="Search your perfect domain (e.g. yourcompany.com)..."
                  value={domainSearch}
                  onChange={(e) => setDomainSearch(e.target.value)}
                  className="pl-10 h-12 border-0 bg-transparent text-base shadow-none focus-visible:ring-0"
                />
              </div>
              <Button type="submit" size="lg" disabled={searching} className="h-12 px-6 gap-2 font-semibold">
                {searching ? 'Checking...' : 'Check Domain'}
                {!searching && <ArrowRight className="h-4 w-4" />}
              </Button>
            </form>
          </Card>

          {/* Domain Search Result Box */}
          {domainResult && (
            <div className="mt-4 p-4 rounded-xl border bg-card flex flex-col sm:flex-row items-center justify-between gap-4 animate-in fade-in zoom-in-95">
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-full bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 flex items-center justify-center font-bold">
                  <Check className="h-5 w-5" />
                </div>
                <div className="text-left">
                  <p className="font-bold text-base">{domainResult.domain}</p>
                  <p className="text-xs text-emerald-600 dark:text-emerald-400 font-semibold">
                    ✓ Domain is available for registration!
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-4">
                <span className="text-lg font-bold">{formatMoney(domainResult.price, domainResult.currency)}/year</span>
                <Button size="sm" onClick={handleAddDomainToCart} className="gap-1 font-semibold">
                  Add to Cart
                </Button>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Pricing / Products Section */}
      <section className="space-y-8">
        <div className="text-center space-y-2">
          <h2 className="text-3xl font-bold tracking-tight">Select Your Service Package</h2>
          <p className="text-sm text-muted-foreground">Transparent pricing with instant automated server provisioning.</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {hostingPlans.map((plan) => {
            const Icon = plan.icon;
            return (
              <Card
                key={plan.id}
                className={`relative flex flex-col border transition-all hover:shadow-xl hover:-translate-y-1 ${
                  plan.popular ? 'border-primary shadow-lg shadow-primary/10 bg-primary/[0.02]' : 'border-border/70'
                }`}
              >
                {plan.popular && (
                  <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                    <Badge className="bg-primary text-primary-foreground text-xs uppercase px-3 font-bold">
                      Most Popular
                    </Badge>
                  </div>
                )}
                <CardHeader>
                  <div className="h-10 w-10 rounded-xl bg-primary/10 text-primary flex items-center justify-center mb-2">
                    <Icon className="h-5 w-5" />
                  </div>
                  <CardTitle className="text-lg">{plan.title}</CardTitle>
                  <CardDescription className="text-xs min-h-[32px]">{plan.description}</CardDescription>
                </CardHeader>
                <CardContent className="flex-1 space-y-4">
                  <div>
                    <span className="text-3xl font-extrabold">{formatMoney(plan.price)}</span>
                    <span className="text-xs text-muted-foreground font-medium">/{plan.period.toLowerCase()}</span>
                  </div>

                  <ul className="space-y-2 text-xs text-muted-foreground pt-2 border-t">
                    {plan.features.map((feat, idx) => (
                      <li key={idx} className="flex items-center gap-2">
                        <Check className="h-3.5 w-3.5 text-emerald-500 shrink-0" />
                        <span>{feat}</span>
                      </li>
                    ))}
                  </ul>
                </CardContent>
                <CardFooter>
                  <Button
                    className="w-full font-semibold gap-1.5"
                    variant={plan.popular ? 'default' : 'outline'}
                    onClick={() => handleAddToCart(plan)}
                  >
                    <Zap className="h-4 w-4" />
                    Order Now
                  </Button>
                </CardFooter>
              </Card>
            );
          })}
        </div>
      </section>

      {/* Latest Announcements */}
      {news.length > 0 && (
        <section className="space-y-4 pt-8 border-t">
          <h3 className="text-xl font-bold tracking-tight">System Announcements & News</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {news.slice(0, 2).map((item) => (
              <Card key={item.id} className="border-border/60">
                <CardHeader className="pb-2">
                  <CardTitle className="text-base">{item.title}</CardTitle>
                  <CardDescription className="text-xs">Published on {new Date(item.created_at).toLocaleDateString()}</CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground line-clamp-2">{item.content}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>
      )}
    </div>
  );
};
