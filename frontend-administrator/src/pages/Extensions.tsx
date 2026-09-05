import React from 'react';
import { RefreshCw, Search } from 'lucide-react';
import { useExtensions } from '@/hooks/useExtensions';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { InstalledModulesTab } from '@/components/extensions/InstalledModulesTab';
import { MarketplaceTab } from '@/components/extensions/MarketplaceTab';

export const Extensions: React.FC = () => {
  const {
    extensions,
    loading,
    searchQuery,
    setSearchQuery,
    selectedType,
    setSelectedType,
    fetchExtensions,
    handleToggle,
  } = useExtensions();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Extensions & Modules Manager</h1>
          <p className="text-sm text-muted-foreground">
            Manage provisioners, payment gateways, theme assets, and install new community plugins.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchExtensions} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      <Tabs defaultValue="installed" className="space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <TabsList>
            <TabsTrigger value="installed">Installed Extensions</TabsTrigger>
            <TabsTrigger value="marketplace">Extension Marketplace</TabsTrigger>
          </TabsList>

          <div className="flex items-center gap-2">
            <div className="relative">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search modules..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-8 h-9 text-xs w-48"
              />
            </div>
            <select
              value={selectedType}
              onChange={(e) => setSelectedType(e.target.value)}
              className="h-9 rounded-md border border-input bg-background px-2.5 py-1 text-xs shadow-sm focus:outline-none"
            >
              <option value="all">All Types</option>
              <option value="service">Service Provisioners</option>
              <option value="gateway">Payment Gateways</option>
              <option value="plugin">Plugins</option>
              <option value="theme">Themes</option>
            </select>
          </div>
        </div>

        <TabsContent value="installed">
          <InstalledModulesTab extensions={extensions} onToggle={handleToggle} />
        </TabsContent>

        <TabsContent value="marketplace">
          <MarketplaceTab />
        </TabsContent>
      </Tabs>
    </div>
  );
};
