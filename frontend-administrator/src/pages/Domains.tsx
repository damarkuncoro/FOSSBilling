import React from 'react';
import { RefreshCw } from 'lucide-react';
import { useDomains } from '@/hooks/useDomains';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { TldPricingTable } from '@/components/domains/TldPricingTable';
import { AddTldDialog } from '@/components/domains/AddTldDialog';
import { RegistrarsConfigTab } from '@/components/domains/RegistrarsConfigTab';

export const Domains: React.FC = () => {
  const {
    tlds,
    registrars,
    loading,
    saving,
    openAddModal,
    setOpenAddModal,
    tldForm,
    setTldForm,
    fetchData,
    handleCreateTld,
    handleDeleteTld,
  } = useDomains();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Domain Names & Registrar Pricing</h1>
          <p className="text-sm text-muted-foreground">
            Configure TLD pricing matrix, registration rules, and domain registrar API connections.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchData} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          <AddTldDialog
            open={openAddModal}
            onOpenChange={setOpenAddModal}
            form={tldForm}
            onChange={setTldForm}
            onSubmit={handleCreateTld}
            saving={saving}
          />
        </div>
      </div>

      <Tabs defaultValue="tlds" className="space-y-4">
        <TabsList>
          <TabsTrigger value="tlds">TLD Extensions & Pricing ({tlds.length})</TabsTrigger>
          <TabsTrigger value="registrars">Registrar Integrations ({registrars.length})</TabsTrigger>
        </TabsList>

        <TabsContent value="tlds" className="space-y-4">
          <Card className="border-border/60 shadow-sm">
            <CardHeader>
              <CardTitle className="text-base font-semibold">TLD Pricing Rules</CardTitle>
              <CardDescription>Pricing structures applied to domain search and renewals</CardDescription>
            </CardHeader>
            <CardContent>
              <TldPricingTable tlds={tlds} loading={loading} onDelete={handleDeleteTld} />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="registrars">
          <RegistrarsConfigTab registrars={registrars} />
        </TabsContent>
      </Tabs>
    </div>
  );
};
