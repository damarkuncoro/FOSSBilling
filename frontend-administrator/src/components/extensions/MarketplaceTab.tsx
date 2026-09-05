import React from 'react';
import { Download, Sparkles } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const MarketplaceTab: React.FC = () => {
  const marketplaceItems = [
    { id: 'telegram_bot', name: 'Telegram Notification Bot', author: 'BotStudio', type: 'plugin', price: 'Free', desc: 'Real-time invoice alerts and new order notifications sent to Telegram groups.' },
    { id: 'whmcs_migrator', name: 'WHMCS & Blesta One-Click Migrator', author: 'FOSSBilling Community', type: 'plugin', price: 'Free', desc: 'Import clients, products, servers, and unpaid invoices from legacy WHMCS databases.' },
    { id: 'dark_neon_theme', name: 'Cyberpunk Dark Glassmorphism Theme', author: 'UI Masters', type: 'theme', price: '$25.00', desc: 'High conversion vibrant theme with customizable CSS variables.' },
  ];

  return (
    <div className="space-y-6">
      <div className="p-4 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Sparkles className="h-5 w-5 text-primary shrink-0" />
          <p className="text-xs text-muted-foreground font-medium">
            Explore verified community extensions, themes, and server provisioners from the FOSSBilling Hub.
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {marketplaceItems.map((item) => (
          <Card key={item.id} className="border-border/60 shadow-sm flex flex-col justify-between">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <Badge variant="outline" className="text-[10px] uppercase font-mono">{item.type}</Badge>
                <span className="font-bold text-xs text-primary">{item.price}</span>
              </div>
              <CardTitle className="text-base mt-2">{item.name}</CardTitle>
              <CardDescription className="text-xs">{item.desc}</CardDescription>
            </CardHeader>
            <CardContent>
              <Button
                size="sm"
                className="w-full text-xs font-semibold gap-1.5"
                onClick={() => alert(`Installing ${item.name}...`)}
              >
                <Download className="h-3.5 w-3.5" />
                Install Extension
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
};
