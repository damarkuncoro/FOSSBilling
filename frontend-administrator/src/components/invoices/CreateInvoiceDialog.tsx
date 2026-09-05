import React, { useState } from 'react';
import { Plus, Trash2, Loader2, FilePlus } from 'lucide-react';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

interface CreateInvoiceDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  clients: Array<{ id: number; first_name: string; last_name: string; email: string; currency?: string }>;
  onInvoiceCreated: () => void;
  apiCreateInvoice: (data: any) => Promise<any>;
}

interface InvoiceLineItem {
  title: string;
  price: number;
  quantity: number;
}

export const CreateInvoiceDialog: React.FC<CreateInvoiceDialogProps> = ({
  open,
  onOpenChange,
  clients,
  onInvoiceCreated,
  apiCreateInvoice,
}) => {
  const [clientId, setClientId] = useState<number>(clients[0]?.id || 1);
  const [currency, setCurrency] = useState<string>('USD');
  const [dueDays, setDueDays] = useState<number>(14);
  const [items, setItems] = useState<InvoiceLineItem[]>([
    { title: 'Custom Hosting Service', price: 29.99, quantity: 1 },
  ]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAddItem = () => setItems((prev) => [...prev, { title: '', price: 0, quantity: 1 }]);
  const handleRemoveItem = (index: number) => items.length > 1 && setItems((prev) => prev.filter((_, i) => i !== index));
  const handleItemChange = (index: number, field: keyof InvoiceLineItem, val: any) => {
    setItems((prev) => {
      const updated = [...prev];
      updated[index] = { ...updated[index], [field]: val };
      return updated;
    });
  };

  const total = items.reduce((acc, it) => acc + (Number(it.price) || 0) * (Number(it.quantity) || 1), 0);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const validItems = items.filter((it) => it.title.trim() !== '');
    if (validItems.length === 0) {
      setError('Please add at least one line item with a title');
      return;
    }

    setLoading(true);
    setError(null);
    try {
      await apiCreateInvoice({
        client_id: Number(clientId),
        currency,
        due_days: Number(dueDays) || 14,
        items: validItems.map((it) => ({
          title: it.title,
          price: Number(it.price) || 0,
          quantity: Number(it.quantity) || 1,
          taxable: false,
        })),
      });
      onOpenChange(false);
      onInvoiceCreated();
    } catch (err: any) {
      setError(err?.message || 'Failed to create invoice');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[580px] max-h-[90vh] overflow-y-auto">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <FilePlus className="h-5 w-5 text-primary" /> Create Custom Invoice
            </DialogTitle>
            <DialogDescription>Generate a new manual billing invoice with custom line items.</DialogDescription>
          </DialogHeader>

          {error && <div className="mt-3 p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm">{error}</div>}

          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="text-xs font-semibold text-muted-foreground block mb-1">Target Client</label>
                <select
                  value={clientId}
                  onChange={(e) => setClientId(Number(e.target.value))}
                  className="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm"
                >
                  {clients.map((c) => (
                    <option key={c.id} value={c.id}>#{c.id} - {c.first_name} {c.last_name} ({c.email})</option>
                  ))}
                </select>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="text-xs font-semibold text-muted-foreground block mb-1">Currency</label>
                  <Input value={currency} onChange={(e) => setCurrency(e.target.value.toUpperCase())} placeholder="USD" className="h-9 font-mono uppercase" maxLength={3} />
                </div>
                <div>
                  <label className="text-xs font-semibold text-muted-foreground block mb-1">Due (Days)</label>
                  <Input type="number" value={dueDays} onChange={(e) => setDueDays(Number(e.target.value))} min={1} className="h-9 font-mono" />
                </div>
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="text-xs font-semibold text-muted-foreground">Invoice Line Items</label>
                <Button type="button" variant="outline" size="sm" onClick={handleAddItem} className="h-7 text-xs gap-1">
                  <Plus className="h-3 w-3" /> Add Item
                </Button>
              </div>

              <div className="space-y-2 max-h-[180px] overflow-y-auto pr-1">
                {items.map((item, idx) => (
                  <div key={idx} className="flex items-center gap-2 p-2 rounded-md bg-muted/40 border border-border/50">
                    <Input placeholder="Item description" value={item.title} onChange={(e) => handleItemChange(idx, 'title', e.target.value)} className="h-8 text-xs flex-1" required />
                    <Input type="number" step="0.01" placeholder="Price" value={item.price} onChange={(e) => handleItemChange(idx, 'price', Number(e.target.value))} className="h-8 text-xs w-24 font-mono" required />
                    <Input type="number" placeholder="Qty" value={item.quantity} onChange={(e) => handleItemChange(idx, 'quantity', Number(e.target.value))} className="h-8 text-xs w-16 font-mono" min={1} required />
                    <Button type="button" variant="ghost" size="sm" onClick={() => handleRemoveItem(idx)} disabled={items.length <= 1} className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive">
                      <Trash2 className="h-3.5 w-3.5" />
                    </Button>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex items-center justify-between p-3 rounded-lg bg-muted/60 border border-border/60">
              <span className="text-sm font-medium">Estimated Subtotal</span>
              <span className="text-base font-bold font-mono text-primary">{total.toLocaleString('en-US', { minimumFractionDigits: 2 })} {currency}</span>
            </div>
          </div>

          <DialogFooter className="gap-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button>
            <Button type="submit" disabled={loading} className="gap-2">
              {loading && <Loader2 className="h-4 w-4 animate-spin" />} Generate Invoice
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};
