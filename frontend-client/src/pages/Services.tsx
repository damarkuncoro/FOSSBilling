import React from 'react';
import { RefreshCw, Download, CheckCircle } from 'lucide-react';
import { useClientServices } from '@/hooks/useClientServices';
import { formatMoney } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';

export const Services: React.FC = () => {
  const {
    orders,
    loading,
    downloadLink,
    downloadModal,
    setDownloadModal,
    fetchServices,
    handleGetDownload,
  } = useClientServices();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">My Active Services & Products</h1>
          <p className="text-sm text-muted-foreground">
            Manage your web hosting accounts, access control panel logins, license keys, and digital downloads.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchServices} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {loading ? (
        <div className="py-16 text-center text-muted-foreground">Loading your cloud services...</div>
      ) : orders.length === 0 ? (
        <Card className="p-8 text-center border-dashed">
          <p className="text-sm text-muted-foreground">You don't have any active subscriptions yet.</p>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {orders.map((order) => (
            <Card key={order.id} className="border-border/60 shadow-sm flex flex-col justify-between">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <Badge variant={order.status === 'active' ? 'success' : 'warning'}>
                    {order.status.toUpperCase()}
                  </Badge>
                  <span className="text-xs font-mono text-muted-foreground">#{order.id}</span>
                </div>
                <CardTitle className="text-base mt-2">{order.title || `Service #${order.product_id}`}</CardTitle>
                <CardDescription className="text-xs">
                  Cycle: {order.period || 'Monthly'} • Price: {formatMoney(order.price, order.currency)}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-3 text-xs">
                {order.server_id ? (
                  <div className="p-3 rounded-lg bg-muted/40 border space-y-1.5 font-mono text-[11px]">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Server Host:</span>
                      <span className="font-semibold text-foreground">sg1.nusantara-cloud.com</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Username:</span>
                      <span className="font-semibold text-foreground">solusinu</span>
                    </div>
                  </div>
                ) : (
                  <div className="p-3 rounded-lg bg-muted/40 border space-y-1 text-muted-foreground">
                    <p className="font-semibold text-foreground">Digital Product / License</p>
                    <p>Instant access verified with HMAC cryptographic tokens.</p>
                  </div>
                )}
              </CardContent>
              <div className="p-6 pt-0 flex gap-2">
                <Button
                  size="sm"
                  variant="outline"
                  className="w-full gap-1.5 text-xs font-semibold"
                  onClick={() => handleGetDownload(order.id)}
                >
                  <Download className="h-3.5 w-3.5" />
                  Get Download / License
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Download Signed Link Modal */}
      {downloadModal && (
        <Dialog open={downloadModal} onOpenChange={setDownloadModal}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2">
                <CheckCircle className="h-5 w-5 text-emerald-500" />
                <span>HMAC Signed Download Ready</span>
              </DialogTitle>
              <DialogDescription>
                Your secure temporary download URL has been generated with SHA-256 signatures.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 my-2">
              <div className="p-3 rounded-lg bg-muted border font-mono text-xs break-all text-primary">
                {downloadLink}
              </div>
              <div className="flex justify-end gap-2">
                <Button
                  variant="outline"
                  onClick={() => {
                    navigator.clipboard.writeText(downloadLink || '');
                    alert('Copied to clipboard!');
                  }}
                >
                  Copy URL
                </Button>
                <a href={downloadLink || '#'} target="_blank" rel="noreferrer">
                  <Button className="gap-1.5">
                    <Download className="h-4 w-4" />
                    Download File
                  </Button>
                </a>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
};
