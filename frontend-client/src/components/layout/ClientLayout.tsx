import React, { useState, useEffect } from 'react';
import { Outlet, Link } from 'react-router-dom';
import { Menu, X, ShoppingCart, Sun, Moon, Shield } from 'lucide-react';
import { ClientSidebar } from './ClientSidebar';
import { Footer } from './Footer';
import { CookieConsentBanner } from '../common/CookieConsentBanner';
import { useClientAuth } from '@/lib/auth';
import { request } from '@/lib/api/client';
import { useCart } from '@/lib/cart';
import { Button } from '@/components/ui/button';
import { NavUserMenu } from './NavUserMenu';
import { DepositModal } from '../invoices/DepositModal';
import { AlertCircle, Undo2 } from 'lucide-react';

import { LanguageSwitcher } from '../common/LanguageSwitcher';

export const ClientLayout: React.FC = () => {
  const { user, balance, isAuthenticated, isImpersonated, logout, theme, toggleTheme } = useClientAuth();
  const { items } = useCart();
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);
  const [depositOpen, setDepositOpen] = useState(false);
  const [branding, setBranding] = useState<any>(null);

  useEffect(() => {
    request<any>('/guest/company').then(data => {
       if (data && data.branding) {
         setBranding(data.branding);
         if (data.branding.primary_color) {
            document.documentElement.style.setProperty('--primary', hexToHSL(data.branding.primary_color));
         }
       }
    }).catch(() => {});
  }, []);

  const hexToHSL = (hex: string) => {
    let r = 0, g = 0, b = 0;
    if (hex.length === 4) {
      r = parseInt(hex[1] + hex[1], 16);
      g = parseInt(hex[2] + hex[2], 16);
      b = parseInt(hex[3] + hex[3], 16);
    } else if (hex.length === 7) {
      r = parseInt(hex.substring(1, 3), 16);
      g = parseInt(hex.substring(3, 5), 16);
      b = parseInt(hex.substring(5, 7), 16);
    }
    r /= 255; g /= 255; b /= 255;
    let max = Math.max(r, g, b), min = Math.min(r, g, b);
    let h = 0, s = 0, l = (max + min) / 2;
    if (max !== min) {
      let d = max - min;
      s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
      switch (max) {
        case r: h = (g - b) / d + (g < b ? 6 : 0); break;
        case g: h = (b - r) / d + 2; break;
        case b: h = (r - g) / d + 4; break;
      }
      h /= 6;
    }
    return `${Math.round(h * 360)} ${Math.round(s * 100)}% ${Math.round(l * 100)}%`;
  };

  return (
    <div className="min-h-screen flex bg-background text-foreground">
      {/* Impersonation Banner */}
      {isAuthenticated && isImpersonated && (
        <div className="fixed top-0 left-0 right-0 z-[100] bg-indigo-600 text-white px-4 py-2 flex items-center justify-center gap-4 shadow-lg animate-in slide-in-from-top duration-300">
           <div className="flex items-center gap-2 text-xs font-bold uppercase tracking-wider">
              <AlertCircle className="h-4 w-4" />
              Viewing as Customer: {user?.first_name} {user?.last_name}
           </div>
           <Button
              size="sm"
              variant="secondary"
              className="h-7 text-[10px] font-black uppercase bg-white text-indigo-600 hover:bg-indigo-50 border-none"
              onClick={() => {
                logout();
                window.close(); // Close if opened in new tab, or redirect
                window.location.href = 'http://admin.fossbilling.test/clients';
              }}
           >
              <Undo2 className="h-3 w-3 mr-1" />
              Return to Admin
           </Button>
        </div>
      )}

      {/* Desktop Fixed Sidebar */}
      {isAuthenticated && (
        <div className="hidden lg:flex lg:w-64 lg:flex-col lg:fixed lg:inset-y-0 z-30">
          <ClientSidebar onOpenDeposit={() => setDepositOpen(true)} />
        </div>
      )}

      {/* Mobile Drawer Overlay */}
      {isAuthenticated && mobileSidebarOpen && (
        <div className="fixed inset-0 z-50 lg:hidden flex">
          <div
            className="fixed inset-0 bg-background/80 backdrop-blur-sm"
            onClick={() => setMobileSidebarOpen(false)}
          />
          <div className="relative flex-1 flex flex-col max-w-xs w-full bg-card shadow-2xl z-10 animate-in slide-in-from-left duration-200">
            <div className="absolute top-3 right-3">
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 text-muted-foreground"
                onClick={() => setMobileSidebarOpen(false)}
              >
                <X className="h-5 w-5" />
              </Button>
            </div>
            <ClientSidebar
              onItemClick={() => setMobileSidebarOpen(false)}
              onOpenDeposit={() => {
                setMobileSidebarOpen(false);
                setDepositOpen(true);
              }}
            />
          </div>
        </div>
      )}

      {/* Main Content Area */}
      <div className={`flex-1 flex flex-col min-w-0 ${isAuthenticated ? 'lg:pl-64' : ''} ${isImpersonated ? 'pt-10' : ''}`}>
        {/* Top Header Bar */}
        <header className="sticky top-0 z-20 h-16 border-b bg-background/80 backdrop-blur-xl flex items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-3">
            {isAuthenticated && (
              <Button
                variant="ghost"
                size="icon"
                className="lg:hidden h-9 w-9"
                onClick={() => setMobileSidebarOpen(true)}
              >
                <Menu className="h-5 w-5" />
              </Button>
            )}
            <Link to="/" className={`flex items-center gap-2 ${isAuthenticated ? 'lg:hidden' : ''}`}>
              {branding?.logo_url ? (
                <img src={branding.logo_url} alt={branding.company_name} className="h-7 object-contain" />
              ) : (
                <>
                  <div className="h-7 w-7 rounded-lg bg-primary flex items-center justify-center text-primary-foreground">
                    <Shield className="h-4 w-4" />
                  </div>
                  <span className="font-bold text-sm">{branding?.company_name || 'FOSSBilling'}</span>
                </>
              )}
            </Link>

            {/* Public Navigation (Visible when not logged in) */}
            {!isAuthenticated && (
              <nav className="hidden md:flex items-center gap-6 ml-8">
                <Link to="/" className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors">Store</Link>
                <Link to="/domains" className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors">Domains</Link>
                <Link to="/kb" className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors">Help Center</Link>
                <Link to="/news" className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors">News</Link>
              </nav>
            )}
          </div>

          <div className="flex items-center gap-3">
            {/* Language Switcher */}
            <LanguageSwitcher />

            {/* Theme Switcher */}
            <Button
              variant="ghost"
              size="icon"
              className="h-9 w-9"
              onClick={toggleTheme}
              title={theme === 'dark' ? 'Switch to Light' : 'Switch to Dark'}
            >
              {theme === 'dark' ? (
                <Sun className="h-4 w-4 text-amber-400" />
              ) : (
                <Moon className="h-4 w-4 text-slate-700" />
              )}
            </Button>

            {/* Shopping Cart */}
            <Link to="/cart">
              <Button variant="outline" size="icon" className="h-9 w-9 relative shadow-sm">
                <ShoppingCart className="h-4 w-4" />
                {items && items.length > 0 && (
                  <span className="absolute -top-1 -right-1 h-4 w-4 rounded-full bg-primary text-primary-foreground text-[10px] font-bold flex items-center justify-center animate-in zoom-in">
                    {items.length}
                  </span>
                )}
              </Button>
            </Link>

            {/* User Dropdown */}
            <NavUserMenu
              user={user}
              balance={balance}
              isAuthenticated={isAuthenticated}
              onLogout={logout}
            />
          </div>
        </header>

        {/* Dynamic Page Content */}
        <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8 animate-in fade-in-50 duration-300">
          <Outlet />
        </main>

        <Footer />
      </div>

      <CookieConsentBanner />

      <DepositModal
        open={depositOpen}
        onOpenChange={setDepositOpen}
        currency={user?.currency || 'USD'}
      />
    </div>
  );
};

export default ClientLayout;
