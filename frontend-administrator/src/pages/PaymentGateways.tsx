import React from 'react';
import { CreditCard, Plus, RefreshCw, Percent } from 'lucide-react';
import { usePaymentGateways } from '@/hooks/usePaymentGateways';
import { EditGatewayDialog } from '@/components/gateways/EditGatewayDialog';
import { AddTaxRuleDialog } from '@/components/gateways/AddTaxRuleDialog';
import { GatewaysListTab } from '@/components/gateways/GatewaysListTab';
import { TaxRulesListTab } from '@/components/gateways/TaxRulesListTab';
import { Button } from '@/components/ui/button';

export const PaymentGateways: React.FC = () => {
  const {
    activeTab,
    setActiveTab,
    gateways,
    taxRules,
    loading,
    selectedGw,
    editGwOpen,
    setEditGwOpen,
    gwForm,
    setGwForm,
    savingGw,
    taxModalOpen,
    setTaxModalOpen,
    taxForm,
    setTaxForm,
    savingTax,
    fetchData,
    openEditGateway,
    handleSaveGateway,
    toggleGatewayEnabled,
    handleAddTaxRule,
    handleDeleteTaxRule,
  } = usePaymentGateways();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Payment Gateways & Tax Rules</h1>
          <p className="text-sm text-muted-foreground">
            Configure checkout payment processors, API credentials, webhooks, and regional tax calculations.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchData} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          {activeTab === 'tax' && (
            <Button size="sm" className="gap-1.5 shadow-sm" onClick={() => setTaxModalOpen(true)}>
              <Plus className="h-4 w-4" />
              Add Tax Rule
            </Button>
          )}
        </div>
      </div>

      <div className="flex border-b">
        <button
          onClick={() => setActiveTab('gateways')}
          className={`pb-3 px-4 text-sm font-semibold border-b-2 transition-all flex items-center gap-2 ${
            activeTab === 'gateways'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <CreditCard className="h-4 w-4" />
          Active Gateways ({gateways.length})
        </button>
        <button
          onClick={() => setActiveTab('tax')}
          className={`pb-3 px-4 text-sm font-semibold border-b-2 transition-all flex items-center gap-2 ${
            activeTab === 'tax'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Percent className="h-4 w-4" />
          Tax & VAT Rules ({taxRules.length})
        </button>
      </div>

      {activeTab === 'gateways' && (
        <GatewaysListTab
          gateways={gateways}
          onToggleEnabled={toggleGatewayEnabled}
          onConfigure={openEditGateway}
        />
      )}

      {activeTab === 'tax' && (
        <TaxRulesListTab taxRules={taxRules} onDelete={handleDeleteTaxRule} />
      )}

      <EditGatewayDialog
        open={editGwOpen}
        onOpenChange={setEditGwOpen}
        selectedGw={selectedGw}
        gwForm={gwForm}
        setGwForm={setGwForm}
        onSave={handleSaveGateway}
        saving={savingGw}
      />

      <AddTaxRuleDialog
        open={taxModalOpen}
        onOpenChange={setTaxModalOpen}
        taxForm={taxForm}
        setTaxForm={setTaxForm}
        onSave={handleAddTaxRule}
        saving={savingTax}
      />
    </div>
  );
};

export default PaymentGateways;
