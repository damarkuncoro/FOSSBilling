import React from 'react';
import { Search, Check, ArrowRight } from 'lucide-react';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';
import { DomainSearchResult } from '@/types/api';

interface DomainSearchCardProps {
  domainSearch: string;
  onDomainSearchChange: (val: string) => void;
  searching: boolean;
  domainResult: DomainSearchResult | null;
  onCheckDomain: (e: React.FormEvent) => void;
  onAddToCart: () => void;
}

export const DomainSearchCard: React.FC<DomainSearchCardProps> = ({
  domainSearch,
  onDomainSearchChange,
  searching,
  domainResult,
  onCheckDomain,
  onAddToCart,
}) => {
  return (
    <div className="max-w-2xl mx-auto pt-4">
      <Card className="p-2 border-border/80 shadow-xl bg-card/80 backdrop-blur-xl">
        <form onSubmit={onCheckDomain} className="flex gap-2">
          <div className="relative flex-1">
            <Search className="absolute left-3.5 top-3.5 h-4 w-4 text-muted-foreground" />
            <Input
              type="text"
              placeholder="Search your perfect domain (e.g. yourcompany.com)..."
              value={domainSearch}
              onChange={(e) => onDomainSearchChange(e.target.value)}
              className="pl-10 h-12 border-0 bg-transparent text-base shadow-none focus-visible:ring-0"
            />
          </div>
          <Button type="submit" size="lg" disabled={searching} className="h-12 px-6 gap-2 font-semibold">
            {searching ? 'Checking...' : 'Check Domain'}
            {!searching && <ArrowRight className="h-4 w-4" />}
          </Button>
        </form>
      </Card>

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
            <span className="text-lg font-bold">
              {formatMoney(domainResult.price, domainResult.currency)}/year
            </span>
            <Button size="sm" onClick={onAddToCart} className="gap-1 font-semibold">
              Add to Cart
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};
