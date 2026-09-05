import React, { useEffect, useState } from 'react';
import { Package, RefreshCw, CheckCircle, PauseCircle, PlayCircle, AlertCircle } from 'lucide-react';
import { api } from '@/lib/api';
import { formatMoney, formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const Orders: React.FC = () => {
  const [orders, setOrders] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState<number | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const data = await api.getOrders();
      setOrders(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, []);

  const handleActivate = async (id: number) => {
    setActionLoading(id);
    setMessage(null);
    try {
      await api.activateOrder(id);
      setMessage(`Order #${id} successfully activated and provisioned!`);
      await fetchOrders();
    } catch (err: any) {
      setMessage(`Error activating order: ${err.message}`);
    } finally {
      setActionLoading(null);
    }
  };

  const handleSuspend = async (id: number) => {
    setActionLoading(id);
    setMessage(null);
    try {
      await api.suspendOrder(id, 'Admin manual suspension');
      setMessage(`Order #${id} has been suspended.`);
      await fetchOrders();
    } catch (err: any) {
      setMessage(`Error suspending order: ${err.message}`);
    } finally {
      setActionLoading(null);
    }
  };

  const handleUnsuspend = async (id: number) => {
    setActionLoading(id);
    setMessage(null);
    try {
      await api.unsuspendOrder(id);
      setMessage(`Order #${id} has been reactivated.`);
      await fetchOrders();
    } catch (err: any) {
      setMessage(`Error unsuspending order: ${err.message}`);
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Orders & Service Provisioning</h1>
          <p className="text-sm text-muted-foreground">
            Manage customer service subscriptions, automated server provisioning, and lifecycle states.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchOrders} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {message && (
        <div className="p-4 rounded-xl bg-primary/10 border border-primary/20 text-primary text-sm flex items-center gap-2">
          <AlertCircle className="h-4 w-4 shrink-0" />
          <span>{message}</span>
        </div>
      )}

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Active & Pending Subscriptions ({orders.length})</CardTitle>
          <CardDescription>Instant automated provisioning with cPanel, DirectAdmin, Plesk, and digital licenses</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-16">ID</TableHead>
                <TableHead>Title / Service</TableHead>
                <TableHead>Client ID</TableHead>
                <TableHead>Billing Period</TableHead>
                <TableHead>Price</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Created</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">
                    Loading orders...
                  </TableCell>
                </TableRow>
              ) : orders.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">
                    No active orders found.
                  </TableCell>
                </TableRow>
              ) : (
                orders.map((order) => (
                  <TableRow key={order.id}>
                    <TableCell className="font-mono text-xs font-semibold">#{order.id}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2 font-medium">
                        <Package className="h-4 w-4 text-primary" />
                        <span>{order.title || `Service #${order.product_id}`}</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-xs font-mono">Client #{order.client_id}</TableCell>
                    <TableCell>
                      <Badge variant="outline" className="text-[11px] font-mono">
                        {order.period || '1M'}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-semibold text-sm">
                      {formatMoney(order.price || 0, order.currency || 'USD')}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          order.status === 'active'
                            ? 'success'
                            : order.status === 'suspended'
                            ? 'destructive'
                            : 'warning'
                        }
                      >
                        {order.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {formatDate(order.created_at)}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-1.5">
                        {order.status === 'pending_setup' && (
                          <Button
                            size="sm"
                            variant="default"
                            className="h-7 text-xs gap-1"
                            disabled={actionLoading === order.id}
                            onClick={() => handleActivate(order.id)}
                          >
                            <CheckCircle className="h-3 w-3" />
                            Activate
                          </Button>
                        )}
                        {order.status === 'active' && (
                          <Button
                            size="sm"
                            variant="destructive"
                            className="h-7 text-xs gap-1"
                            disabled={actionLoading === order.id}
                            onClick={() => handleSuspend(order.id)}
                          >
                            <PauseCircle className="h-3 w-3" />
                            Suspend
                          </Button>
                        )}
                        {order.status === 'suspended' && (
                          <Button
                            size="sm"
                            variant="outline"
                            className="h-7 text-xs gap-1"
                            disabled={actionLoading === order.id}
                            onClick={() => handleUnsuspend(order.id)}
                          >
                            <PlayCircle className="h-3 w-3" />
                            Unsuspend
                          </Button>
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
