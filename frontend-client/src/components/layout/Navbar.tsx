import React, { useState } from 'react';
import { NavLink, Link, useNavigate } from 'react-router-dom';
import { Shield, ShoppingCart, Moon, Sun, Menu, X } from 'lucide-react';
import { useClientAuth } from '@/lib/auth';
import { useCart } from '@/lib/cart';
import { Button } from '@/components/ui/button';
import { NavMobileMenu } from './NavMobileMenu';
import { NavUserMenu } from './NavUserMenu';

export const Navbar: React.FC = () => {
  const { user, balance, isAuthenticated, logout, theme, toggleTheme } = useClientAuth();
  const { items } = useCart();
  const navigate = useNavigate();
  const [mobileMenu, setMobileMenu] = useState(false);

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
              <span className="text-[10px] text-muted-foreground font-semibold">
                Client Cloud Portal
              </span>
            </div>
          </Link>

          {/* Desktop Nav */}
          <nav className="hidden md:flex items-center gap-5 text-sm font-medium">
            <NavLink
              to="/"
              className={({ isActive }) =>
                `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
              }
            >
              Store
            </NavLink>
            <NavLink
              to="/domains"
              className={({ isActive }) =>
                `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
              }
            >
              Domains
            </NavLink>
            <NavLink
              to="/kb"
              className={({ isActive }) =>
                `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
              }
            >
              Knowledgebase
            </NavLink>
            <NavLink
              to="/news"
              className={({ isActive }) =>
                `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
              }
            >
              News
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
                  Services
                </NavLink>
                <NavLink
                  to="/licenses"
                  className={({ isActive }) =>
                    `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
                  }
                >
                  Licenses
                </NavLink>
                <NavLink
                  to="/downloads"
                  className={({ isActive }) =>
                    `transition-colors hover:text-primary ${isActive ? 'text-primary' : 'text-muted-foreground'}`
                  }
                >
                  Downloads
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
            {theme === 'dark' ? (
              <Sun className="h-4 w-4 text-amber-400" />
            ) : (
              <Moon className="h-4 w-4 text-slate-700" />
            )}
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

          {/* User Profile / Auth State */}
          <NavUserMenu
            user={user}
            balance={balance}
            isAuthenticated={isAuthenticated}
            onLogout={handleLogout}
          />

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
      <NavMobileMenu
        open={mobileMenu}
        onClose={() => setMobileMenu(false)}
        isAuthenticated={isAuthenticated}
        onLogout={handleLogout}
      />
    </header>
  );
};
