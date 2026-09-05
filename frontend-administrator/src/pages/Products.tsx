import React from 'react';
import { Package, Search, RefreshCw, Server, Globe, KeyRound, Download } from 'lucide-react';
import { useProducts } from '@/hooks/useProducts';
import { AddProductDialog } from '@/components/products/AddProductDialog';
import { ProductsTable } from '@/components/products/ProductsTable';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

const productTabs = [
  { key: 'all', label: 'All Products', icon: Package },
  { key: 'hosting', label: 'Hosting', icon: Server },
  { key: 'domain', label: 'Domains', icon: Globe },
  { key: 'license', label: 'Licenses', icon: KeyRound },
  { key: 'downloadable', label: 'Downloads', icon: Download },
];

export const Products: React.FC = () => {
  const {
    loading,
    searchQuery,
    setSearchQuery,
    selectedType,
    setSelectedType,
    openModal,
    setOpenModal,
    saving,
    form,
    setForm,
    fetchProducts,
    handleTitleChange,
    handleSaveProduct,
    handleDelete,
    toggleStatus,
    filteredProducts,
  } = useProducts();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Products & Services</h1>
          <p className="text-sm text-muted-foreground">
            Configure your catalog, hosting packages, domain pricing, software licenses, and digital downloads.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchProducts} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          <AddProductDialog
            open={openModal}
            onOpenChange={setOpenModal}
            form={form}
            setForm={setForm}
            onTitleChange={handleTitleChange}
            onSave={handleSaveProduct}
            saving={saving}
          />
        </div>
      </div>

      <div className="flex flex-col sm:flex-row items-center justify-between gap-3">
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search products, descriptions, slugs..."
            className="pl-9 h-9"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

        <div className="flex items-center gap-1.5 overflow-x-auto w-full sm:w-auto pb-1 sm:pb-0">
          {productTabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <Button
                key={tab.key}
                variant={selectedType === tab.key ? 'secondary' : 'ghost'}
                size="sm"
                onClick={() => setSelectedType(tab.key)}
                className="gap-1.5 text-xs font-medium h-8"
              >
                <Icon className="h-3.5 w-3.5" />
                {tab.label}
              </Button>
            );
          })}
        </div>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader className="pb-3">
          <CardTitle className="text-base font-semibold">
            Catalog Items ({filteredProducts.length})
          </CardTitle>
          <CardDescription>All active products synced with customer store ordering</CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          <ProductsTable
            products={filteredProducts}
            loading={loading}
            onToggleStatus={toggleStatus}
            onDelete={handleDelete}
          />
        </CardContent>
      </Card>
    </div>
  );
};

export default Products;
