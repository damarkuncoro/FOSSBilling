import React from 'react';
import { ShieldCheck, Heart } from 'lucide-react';

export const Footer: React.FC = () => {
  return (
    <footer className="border-t bg-muted/20 py-10 mt-20">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-muted-foreground">
        <div className="flex items-center gap-2">
          <ShieldCheck className="h-4 w-4 text-primary" />
          <span>FOSSBilling Next-Gen Powered by Golang Engine</span>
        </div>
        <p className="flex items-center gap-1">
          Built with <Heart className="h-3 w-3 text-red-500 fill-red-500" /> for hosting & cloud providers
        </p>
        <p>© {new Date().getFullYear()} FOSSBilling. Open Source Software under Apache-2.0.</p>
      </div>
    </footer>
  );
};
