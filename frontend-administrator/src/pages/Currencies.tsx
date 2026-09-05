import React from 'react';
import { RefreshCw, CheckCircle, Trash2 } from 'lucide-react';
import { useCurrencies } from '@/hooks/useCurrencies';
import { AddCurrencyDialog } from '@/components/currencies/AddCurrencyDialog';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const Currencies: React.FC = () => {
  const {
    currencies,
    loading,
    openModal,
    setOpenModal,
    form,
    setForm,
    saving,
    fetchCurrencies,
    handleCreate,
    handleSetDefault,
    handleDelete,
  } = useCurrencies();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Currencies & Exchange Rates</h1>
          <p className="text-sm text-muted-foreground">
            Manage supported multi-currencies, daily conversion rates, and localized price formatting.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchCurrencies} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          <AddCurrencyDialog
            open={openModal}
            onOpenChange={setOpenModal}
            form={form}
            setForm={setForm}
            onSave={handleCreate}
            saving={saving}
          />
        </div>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Active Currencies ({currencies.length})</CardTitle>
          <CardDescription>Real-time exchange conversion applied to carts, invoices, and catalog prices</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Code</TableHead>
                <TableHead>Title</TableHead>
                <TableHead>Exchange Rate</TableHead>
                <TableHead>Display Format</TableHead>
                <TableHead>Default</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                    Loading currencies...
                  </TableCell>
                </TableRow>
              ) : (
                currencies.map((curr) => (
                  <TableRow key={curr.code}>
                    <TableCell className="font-mono text-sm font-bold">{curr.code}</TableCell>
                    <TableCell className="font-medium">{curr.title}</TableCell>
                    <TableCell className="font-mono text-xs">{curr.conversion_rate}</TableCell>
                    <TableCell className="font-mono text-xs text-muted-foreground">{curr.format}</TableCell>
                    <TableCell>
                      {curr.is_default ? (
                        <Badge variant="success" className="gap-1">
                          <CheckCircle className="h-3 w-3" />
                          Default Primary
                        </Badge>
                      ) : (
                        <Badge variant="outline">Secondary</Badge>
                      )}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        {!curr.is_default && (
                          <>
                            <Button
                              size="sm"
                              variant="outline"
                              className="h-7 text-xs"
                              onClick={() => handleSetDefault(curr.code)}
                            >
                              Set Default
                            </Button>
                            <Button
                              size="icon"
                              variant="ghost"
                              className="h-7 w-7 text-muted-foreground hover:text-destructive"
                              onClick={() => handleDelete(curr.code)}
                            >
                              <Trash2 className="h-3.5 w-3.5" />
                            </Button>
                          </>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default Currencies;
