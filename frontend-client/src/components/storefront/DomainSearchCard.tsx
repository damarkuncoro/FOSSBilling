import React from 'react';
import { Search, Check, X, ArrowRight } from 'lucide-react';
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
            <div
              className={`h-10 w-10 rounded-full flex items-center justify-center font-bold ${
                domainResult.available
                  ? 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400'
                  : 'bg-rose-500/15 text-rose-600 dark:text-rose-400'
              }`}
            >
              {domainResult.available ? <Check className="h-5 w-5" /> : <X className="h-5 w-5" />}
            </div>
            <div className="text-left">
              <p className="font-bold text-base">{domainResult.domain}</p>
              <p
                className={`text-xs font-semibold ${
                  domainResult.available
                    ? 'text-emerald-600 dark:text-emerald-400'
                    : 'text-rose-600 dark:text-rose-400'
                }`}
              >
                {domainResult.available
                  ? '✓ Domain is available for registration!'
                  : '✕ Domain is already registered / taken.'}
              </p>
            </div>
          </div>
          {domainResult.available ? (
            <div className="flex items-center gap-4">
              <span className="text-lg font-bold">
                {formatMoney(domainResult.price, domainResult.currency)}/year
              </span>
              <Button size="sm" onClick={onAddToCart} className="gap-1 font-semibold">
                Add to Cart
              </Button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <span className="text-xs px-3 py-1.5 rounded-lg bg-rose-100 text-rose-700 dark:bg-rose-950 dark:text-rose-300 font-medium">
                Unavailable
              </span>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

