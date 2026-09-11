import React from 'react';
import { Globe, Search, CheckCircle2, XCircle, ShoppingCart, RefreshCw, MoreHorizontal, Settings2, Trash2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useClientDomains } from '../hooks/useClientDomains';
import { useCart } from '../lib/cart';
import { ManageDnsDialog } from '../components/domains/ManageDnsDialog';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

export const Domains: React.FC = () => {
  const { t } = useTranslation();
  const {
    domains,
    loading,
    search,
    setSearch,
    checkQuery,
    setCheckQuery,
    checkResult,
    isSearching,
    editingDomain,
    setEditingDomain,
    checkAvailability,
    updateNameservers,
    toggleAutoRenew,
  } = useClientDomains();

  const { addItem } = useCart();

  const handleAddToCart = () => {
    if (!checkResult) return;
    addItem({
      id: `cart_domain_${Date.now()}`,
      product_id: 99,
      title: `Domain Registration: ${checkResult.domain}`,
      type: 'domain',
      domain_name: checkResult.domain,
      period: '1Y',
      price: checkResult.price,
    });
  };

  return (
    <div className="space-y-8 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Domains & DNS Management</h1>
          <p className="text-sm text-muted-foreground">
            Search for new domain names or manage your active registrations and DNS zones.
          </p>
        </div>
      </div>

      {/* Domain Lookup Card */}
      <Card className="bg-indigo-900 text-indigo-50 border-none shadow-xl overflow-hidden relative">
        <div className="absolute top-0 right-0 p-8 opacity-10">
           <Globe className="w-40 h-40" />
        </div>
        <CardHeader className="relative z-10">
          <CardTitle className="text-indigo-300 text-xs font-bold uppercase tracking-widest">WHOIS Checker</CardTitle>
          <h2 className="text-2xl font-bold text-white">Register Your Next Domain</h2>
        </CardHeader>
        <CardContent className="relative z-10 space-y-6">
          <div className="flex flex-col sm:flex-row gap-3">
            <Input
              placeholder="Find your new domain (e.g. mycompany.com)..."
              value={checkQuery}
              onChange={(e) => setCheckQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && checkAvailability()}
              className="flex-1 bg-white/10 border-white/20 text-white placeholder:text-indigo-300/60 h-12 rounded-xl"
            />
            <Button
              onClick={checkAvailability}
              disabled={isSearching}
              className="bg-indigo-500 hover:bg-indigo-400 text-white font-bold h-12 px-8 rounded-xl shadow-lg"
            >
              {isSearching ? <RefreshCw className="w-4 h-4 animate-spin mr-2" /> : <Search className="w-4 h-4 mr-2" />}
              Check Availability
            </Button>
          </div>

          {checkResult && (
            <div className="p-4 bg-white/5 rounded-xl border border-white/10 backdrop-blur-md flex items-center justify-between animate-in zoom-in-95 duration-200">
              <div className="flex items-center gap-4">
                {checkResult.available ? <CheckCircle2 className="w-6 h-6 text-emerald-400" /> : <XCircle className="w-6 h-6 text-rose-400" />}
                <div>
                  <span className="font-mono font-bold text-lg text-white">{checkResult.domain}</span>
                  <p className="text-sm text-indigo-200">
                    {checkResult.available ? `Available for $${checkResult.price}/year` : 'Already registered by another owner.'}
                  </p>
                </div>
              </div>
              {checkResult.available && (
                <Button
                  onClick={handleAddToCart}
                  variant="secondary"
                  className="bg-emerald-500 hover:bg-emerald-400 text-white border-none font-bold"
                >
                  <ShoppingCart className="w-4 h-4 mr-2" /> Add to Cart
                </Button>
              )}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Active Domains */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-bold flex items-center gap-2">
            <Globe className="w-5 h-5 text-primary" /> My Active Domains
          </h3>
          <Input
            placeholder="Filter domains..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-64 h-9 text-xs"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {loading ? (
            Array.from({ length: 3 }).map((_, i) => (
              <Card key={i} className="border-border/60">
                 <CardHeader><Skeleton className="h-4 w-32" /><Skeleton className="h-3 w-40 mt-2" /></CardHeader>
                 <CardContent><Skeleton className="h-10 w-full rounded-lg" /></CardContent>
              </Card>
            ))
          ) : domains.length === 0 ? (
            <Card className="col-span-full py-12 border-dashed flex flex-col items-center justify-center text-muted-foreground">
               <Globe className="w-10 h-10 mb-3 opacity-20" />
               <p className="text-sm font-medium">No domains found in your account.</p>
            </Card>
          ) : (
            domains.map((d) => (
              <Card key={d.id} className="border-border/60 hover:shadow-md transition-all group overflow-hidden">
                <CardHeader className="pb-3">
                   <div className="flex justify-between items-start">
                      <Badge variant={d.status === 'active' ? 'success' : 'warning'} className="uppercase text-[10px]">
                        {d.status}
                      </Badge>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                           <Button variant="ghost" size="icon" className="h-8 w-8">
                              <MoreHorizontal className="h-4 w-4" />
                           </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                           <DropdownMenuItem className="gap-2" onClick={() => toggleAutoRenew(d.id)}>
                              <RefreshCw className="h-3.5 w-3.5" />
                              {d.auto_renew ? 'Disable Auto-Renew' : 'Enable Auto-Renew'}
                           </DropdownMenuItem>
                           <DropdownMenuItem className="gap-2 text-destructive">
                              <Trash2 className="h-3.5 w-3.5" />
                              Request Transfer
                           </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                   </div>
                   <CardTitle className="text-lg font-mono font-bold mt-2">{d.domain_name}</CardTitle>
                   <CardDescription className="text-xs">
                      Expires: {new Date(d.expires_at).toLocaleDateString()}
                   </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                   <div className="p-3 rounded-lg bg-muted/40 border text-[10px] font-mono space-y-1">
                      <div className="text-muted-foreground uppercase font-bold text-[9px] mb-1">Nameservers</div>
                      {d.nameservers.map((ns, idx) => (
                        <div key={idx}>{ns}</div>
                      ))}
                   </div>
                   <Button
                    onClick={() => setEditingDomain(d)}
                    className="w-full gap-2 text-xs font-semibold"
                    variant="outline"
                   >
                    <Settings2 className="w-3.5 h-3.5" /> Manage DNS Zone
                   </Button>
                </CardContent>
              </Card>
            ))
          )}
        </div>
      </div>

      <ManageDnsDialog
        domain={editingDomain}
        onClose={() => setEditingDomain(null)}
        onSave={updateNameservers}
      />
    </div>
  );
};
