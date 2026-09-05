import React, { useState } from 'react';
import { Outlet, Link } from 'react-router-dom';
import { Menu, X, ShoppingCart, Sun, Moon, Shield } from 'lucide-react';
import { ClientSidebar } from './ClientSidebar';
import { Footer } from './Footer';
import { CookieConsentBanner } from '../common/CookieConsentBanner';
import { useClientAuth } from '@/lib/auth';
import { useCart } from '@/lib/cart';
import { Button } from '@/components/ui/button';
import { NavUserMenu } from './NavUserMenu';
import { DepositModal } from '../invoices/DepositModal';

export const ClientLayout: React.FC = () => {
  const { user, balance, isAuthenticated, logout, theme, toggleTheme } = useClientAuth();
  const { items } = useCart();
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);
  const [depositOpen, setDepositOpen] = useState(false);

  return (
    <div className="min-h-screen flex bg-background text-foreground">
      {/* Desktop Fixed Sidebar */}
      <div className="hidden lg:flex lg:w-64 lg:flex-col lg:fixed lg:inset-y-0 z-30">
        <ClientSidebar onOpenDeposit={() => setDepositOpen(true)} />
      </div>

      {/* Mobile Drawer Overlay */}
      {mobileSidebarOpen && (
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
      <div className="flex-1 flex flex-col lg:pl-64 min-w-0">
        {/* Top Header Bar */}
        <header className="sticky top-0 z-20 h-16 border-b bg-background/80 backdrop-blur-xl flex items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="icon"
              className="lg:hidden h-9 w-9"
              onClick={() => setMobileSidebarOpen(true)}
            >
              <Menu className="h-5 w-5" />
            </Button>
            <Link to="/" className="flex items-center gap-2 lg:hidden">
              <div className="h-7 w-7 rounded-lg bg-primary flex items-center justify-center text-primary-foreground">
                <Shield className="h-4 w-4" />
              </div>
              <span className="font-bold text-sm">FOSSBilling</span>
            </Link>
          </div>

          <div className="flex items-center gap-3">
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
                {items.length > 0 && (
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
