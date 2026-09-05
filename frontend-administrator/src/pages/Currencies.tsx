import React, { useEffect, useState } from 'react';
import { Coins, RefreshCw, Plus, CheckCircle, Trash2 } from 'lucide-react';
import { api } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';

export const Currencies: React.FC = () => {
  const [currencies, setCurrencies] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState({
    code: '',
    title: '',
    conversion_rate: 1.0,
    format: '$ {{price}}',
  });
  const [saving, setSaving] = useState(false);

  const fetchCurrencies = async () => {
    setLoading(true);
    try {
      const data = await api.getCurrencies();
      setCurrencies(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCurrencies();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await api.createCurrency({
        code: form.code.toUpperCase(),
        title: form.title,
        conversion_rate: Number(form.conversion_rate),
        format: form.format,
      });
      setOpenModal(false);
      setForm({ code: '', title: '', conversion_rate: 1.0, format: '$ {{price}}' });
      await fetchCurrencies();
    } catch (err) {
      console.error(err);
    } finally {
      setSaving(false);
    }
  };

  const handleSetDefault = async (code: string) => {
    try {
      await api.setDefaultCurrency(code);
      await fetchCurrencies();
    } catch (err) {
      console.error(err);
    }
  };

  const handleDelete = async (code: string) => {
    if (!confirm(`Are you sure you want to delete currency ${code}?`)) return;
    try {
      await api.deleteCurrency(code);
      await fetchCurrencies();
    } catch (err) {
      console.error(err);
    }
  };

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

          <Dialog open={openModal} onOpenChange={setOpenModal}>
            <DialogTrigger asChild>
              <Button size="sm" className="gap-1.5 shadow-sm">
                <Plus className="h-4 w-4" />
                Add Currency
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New Currency</DialogTitle>
                <DialogDescription>Define a new international currency and its conversion rate</DialogDescription>
              </DialogHeader>
              <form onSubmit={handleCreate} className="space-y-4 pt-2">
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">ISO Currency Code</label>
                  <Input
                    required
                    placeholder="e.g. JPY, GBP, AUD"
                    value={form.code}
                    onChange={(e) => setForm({ ...form, code: e.target.value })}
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">Title / Description</label>
                  <Input
                    required
                    placeholder="e.g. Japanese Yen"
                    value={form.title}
                    onChange={(e) => setForm({ ...form, title: e.target.value })}
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">Conversion Rate (relative to default)</label>
                  <Input
                    type="number"
                    step="any"
                    required
                    value={form.conversion_rate}
                    onChange={(e) => setForm({ ...form, conversion_rate: parseFloat(e.target.value) || 1 })}
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">Display Format Template</label>
                  <Input
                    required
                    placeholder="e.g. ¥ {{price}} or Rp {{price}}"
                    value={form.format}
                    onChange={(e) => setForm({ ...form, format: e.target.value })}
                  />
                </div>
                <div className="flex justify-end gap-2 pt-2">
                  <Button type="button" variant="outline" onClick={() => setOpenModal(false)}>
                    Cancel
                  </Button>
                  <Button type="submit" disabled={saving}>
                    {saving ? 'Creating...' : 'Save Currency'}
                  </Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
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
              ) : currencies.map((curr) => (
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
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};
