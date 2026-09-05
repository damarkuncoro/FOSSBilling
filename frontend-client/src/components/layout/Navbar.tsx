import React from 'react';
import { NavLink, Link, useNavigate } from 'react-router-dom';
import {
  Shield,
  ShoppingCart,
  User,
  LayoutDashboard,
  Package,
  FileText,
  LifeBuoy,
  LogOut,
  Moon,
  Sun,
  Wallet,
  Menu,
  X,
} from 'lucide-react';
import { useClientAuth } from '@/lib/auth';
import { useCart } from '@/lib/cart';
import { formatMoney } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';

export const Navbar: React.FC = () => {
  const { user, balance, isAuthenticated, logout, theme, toggleTheme } = useClientAuth();
  const { items } = useCart();
  const navigate = useNavigate();
  const [mobileMenu, setMobileMenu] = React.useState(false);

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <header className="sticky top-0 z-40 w-full border-b bg-background/80 backdrop-blur-xl">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        {/* Brand */}
        <div className="flex items-center gap-8">
          <Link to="/" className="flex items-center gap-2.5 group">
            <div className="h-9 w-9 rounded-xl bg-gradient-to-tr from-primary to-indigo-500 flex items-center justify-center text-white shadow-md shadow-primary/25 group-hover:scale-105 transition-transform">
              <Shield className="h-5 w-5" />
            </div>
            <div>
              <span className="font-bold text-base tracking-tight block leading-none">
                FOSSBilling
              </span>
              <span className="text-[10px] text-muted-foreground font-semibold">Client Cloud Portal</span>
            </div>
          </Link>

          {/* Desktop Nav */}
          <nav className="hidden md:flex items-center gap-6 text-sm font-medium">
            <NavLink
              to="/"
              className={({ isActive }) =>
                `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
              }
            >
              Store & Hosting
            </NavLink>
            {isAuthenticated && (
              <>
                <NavLink
                  to="/dashboard"
                  className={({ isActive }) =>
                    `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
                  }
                >
                  Dashboard
                </NavLink>
                <NavLink
                  to="/services"
                  className={({ isActive }) =>
                    `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
                  }
                >
                  My Services
                </NavLink>
                <NavLink
                  to="/invoices"
                  className={({ isActive }) =>
                    `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
                  }
                >
                  Invoices
                </NavLink>
                <NavLink
                  to="/support"
                  className={({ isActive }) =>
                    `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
                  }
                >
                  Support
                </NavLink>
              </>
            )}
          </nav>
        </div>

        {/* Right Controls */}
        <div className="flex items-center gap-3">
          {/* Theme Toggle */}
          <Button
            variant="ghost"
            size="icon"
            className="h-9 w-9"
            onClick={toggleTheme}
            title={theme === 'dark' ? 'Switch to Light' : 'Switch to Dark'}
          >
            {theme === 'dark' ? <Sun className="h-4 w-4 text-amber-400" /> : <Moon className="h-4 w-4 text-slate-700" />}
          </Button>

          {/* Cart Icon */}
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

          {/* Auth State */}
          {isAuthenticated ? (
            <div className="flex items-center gap-3">
              {/* Balance Badge */}
              <div className="hidden sm:flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 text-xs font-semibold">
                <Wallet className="h-3.5 w-3.5" />
                <span>{formatMoney(balance, user?.currency || 'IDR')}</span>
              </div>

              {/* User Avatar */}
              <Link to="/settings" className="flex items-center gap-2 group">
                <Avatar className="h-8 w-8 border border-border group-hover:border-primary transition-colors">
                  <AvatarFallback className="bg-primary/10 text-primary font-bold text-xs">
                    {user?.first_name?.[0] || 'U'}
                  </AvatarFallback>
                </Avatar>
              </Link>

              <Button variant="ghost" size="icon" onClick={handleLogout} className="h-8 w-8 text-muted-foreground hover:text-destructive" title="Logout">
                <LogOut className="h-4 w-4" />
              </Button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <Link to="/login">
                <Button variant="ghost" size="sm" className="text-xs font-semibold">
                  Sign In
                </Button>
              </Link>
              <Link to="/register">
                <Button size="sm" className="text-xs font-semibold shadow-sm">
                  Register
                </Button>
              </Link>
            </div>
          )}

          {/* Mobile Menu Toggle */}
          <Button
            variant="ghost"
            size="icon"
            className="md:hidden h-9 w-9"
            onClick={() => setMobileMenu(!mobileMenu)}
          >
            {mobileMenu ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </Button>
        </div>
      </div>

      {/* Mobile Dropdown */}
      {mobileMenu && (
        <div className="md:hidden border-b bg-card p-4 space-y-3">
          <Link to="/" onClick={() => setMobileMenu(false)} className="block text-sm font-medium py-1">
            Store & Hosting
          </Link>
          {isAuthenticated ? (
            <>
              <Link to="/dashboard" onClick={() => setMobileMenu(false)} className="block text-sm font-medium py-1">
                Dashboard
              </Link>
              <Link to="/services" onClick={() => setMobileMenu(false)} className="block text-sm font-medium py-1">
                My Services
              </Link>
              <Link to="/invoices" onClick={() => setMobileMenu(false)} className="block text-sm font-medium py-1">
                Invoices
              </Link>
              <Link to="/support" onClick={() => setMobileMenu(false)} className="block text-sm font-medium py-1">
                Support Tickets
              </Link>
              <Link to="/settings" onClick={() => setMobileMenu(false)} className="block text-sm font-medium py-1">
                Account Settings
              </Link>
              <Button variant="destructive" size="sm" onClick={handleLogout} className="w-full mt-2">
                Logout
              </Button>
            </>
          ) : (
            <div className="flex gap-2 pt-2">
              <Link to="/login" className="flex-1" onClick={() => setMobileMenu(false)}>
                <Button variant="outline" className="w-full">Sign In</Button>
              </Link>
              <Link to="/register" className="flex-1" onClick={() => setMobileMenu(false)}>
                <Button className="w-full">Register</Button>
              </Link>
            </div>
          )}
        </div>
      )}
    </header>
  );
};
