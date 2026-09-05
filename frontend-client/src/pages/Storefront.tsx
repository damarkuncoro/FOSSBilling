import React from 'react';
import { Sparkles } from 'lucide-react';
import { useStorefront } from '@/hooks/useStorefront';
import { DomainSearchCard } from '@/components/storefront/DomainSearchCard';
import { ProductPlanCard } from '@/components/storefront/ProductPlanCard';
import { NewsAnnouncements } from '@/components/storefront/NewsAnnouncements';

export const Storefront: React.FC = () => {
  const {
    domainSearch,
    setDomainSearch,
    domainResult,
    searching,
    news,
    handleDomainCheck,
    handleAddToCart,
    handleAddDomainToCart,
    plans,
  } = useStorefront();

  return (
    <div className="space-y-16">
      {/* Hero Section */}
      <section className="text-center py-12 md:py-16 space-y-6 relative overflow-hidden">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-xs font-semibold">
          <Sparkles className="h-3.5 w-3.5" />
          <span>Next-Generation Cloud Infrastructure Powered by Go</span>
        </div>

        <h1 className="text-4xl md:6xl font-extrabold tracking-tight max-w-3xl mx-auto leading-tight">
          Superfast Cloud Hosting & Digital Services
        </h1>
        <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
          Deploy high-speed cPanel & DirectAdmin hosting, VPS cloud instances, and enterprise software licenses with instant automated provisioning.
        </p>

        {/* Domain Search Card */}
        <DomainSearchCard
          domainSearch={domainSearch}
          onDomainSearchChange={setDomainSearch}
          searching={searching}
          domainResult={domainResult}
          onCheckDomain={handleDomainCheck}
          onAddToCart={handleAddDomainToCart}
        />
      </section>

      {/* Pricing / Products Section */}
      <section className="space-y-8">
        <div className="text-center space-y-2">
          <h2 className="text-3xl font-bold tracking-tight">Select Your Service Package</h2>
          <p className="text-sm text-muted-foreground">
            Transparent pricing with instant automated server provisioning.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {plans.map((plan) => (
            <ProductPlanCard
              key={plan.id}
              plan={plan}
              onAddToCart={handleAddToCart}
            />
          ))}
        </div>
      </section>

      {/* Latest Announcements */}
      <NewsAnnouncements news={news} />
    </div>
  );
};
