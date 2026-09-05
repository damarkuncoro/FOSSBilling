import React from 'react';
import {
  Tag,
  RefreshCw,
  Trash2,
  CheckCircle2,
  XCircle,
} from 'lucide-react';
import { useCoupons } from '@/hooks/useCoupons';
import { AddCouponDialog } from '@/components/coupons/AddCouponDialog';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const Coupons: React.FC = () => {
  const {
    coupons,
    loading,
    openModal,
    setOpenModal,
    saving,
    form,
    setForm,
    fetchCoupons,
    generateRandomCode,
    handleCreate,
    handleDelete,
    toggleStatus,
  } = useCoupons();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Promo Codes & Discounts</h1>
          <p className="text-sm text-muted-foreground">
            Create promotional discount codes, percentage vouchers, and limit usage per customer campaign.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchCoupons} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          <AddCouponDialog
            open={openModal}
            onOpenChange={setOpenModal}
            form={form}
            setForm={setForm}
            onGenerateCode={generateRandomCode}
            onSave={handleCreate}
            saving={saving}
          />
        </div>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Promotional Vouchers ({coupons.length})</CardTitle>
          <CardDescription>Active and expired promotional vouchers configured for order discounts</CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Coupon Code</TableHead>
                <TableHead>Discount Value</TableHead>
                <TableHead>Usage Progress</TableHead>
                <TableHead>Expiration</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-10 text-muted-foreground">
                    <RefreshCw className="h-5 w-5 animate-spin mx-auto mb-2 text-primary" />
                    Loading coupons...
                  </TableCell>
                </TableRow>
              ) : coupons.map((c) => {
                const isExpired = c.expires_at && new Date(c.expires_at) < new Date();
                const usagePercent = Math.round((c.used_count / (c.max_uses || 1)) * 100);

                return (
                  <TableRow key={c.id} className="hover:bg-muted/40 transition-colors">
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Tag className="h-4 w-4 text-primary" />
                        <span className="font-mono font-bold text-sm tracking-wide bg-muted px-2 py-0.5 rounded border">
                          {c.code}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <span className="font-bold text-sm text-emerald-600 dark:text-emerald-400">
                        {c.type === 'percentage' ? `${c.value}% OFF` : `$${c.value.toFixed(2)} OFF`}
                      </span>
                    </TableCell>
                    <TableCell>
                      <div className="space-y-1 w-32">
                        <div className="flex justify-between text-[11px]">
                          <span>{c.used_count} / {c.max_uses}</span>
                          <span className="text-muted-foreground">{usagePercent}%</span>
                        </div>
                        <div className="w-full bg-muted rounded-full h-1.5 overflow-hidden">
                          <div
                            className={`h-full rounded-full transition-all ${
                              usagePercent >= 100 ? 'bg-destructive' : 'bg-primary'
                            }`}
                            style={{ width: `${Math.min(usagePercent, 100)}%` }}
                          />
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {c.expires_at ? (
                        <span className={isExpired ? 'text-destructive font-medium' : ''}>
                          {c.expires_at} {isExpired ? '(Expired)' : ''}
                        </span>
                      ) : (
                        'Never (Unlimited)'
                      )}
                    </TableCell>
                    <TableCell>
                      <button onClick={() => toggleStatus(c.id)} className="cursor-pointer focus:outline-none">
                        {c.is_active && !isExpired ? (
                          <Badge variant="success" className="gap-1 text-[11px]">
                            <CheckCircle2 className="h-3 w-3" />
                            Active
                          </Badge>
                        ) : (
                          <Badge variant="outline" className="gap-1 text-[11px]">
                            <XCircle className="h-3 w-3" />
                            Disabled
                          </Badge>
                        )}
                      </button>
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        size="icon"
                        variant="ghost"
                        className="h-7 w-7 text-muted-foreground hover:text-destructive"
                        onClick={() => handleDelete(c.id)}
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                      </Button>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default Coupons;
