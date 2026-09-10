import React, { useState, useEffect } from 'react';
import { Link, Outlet, useNavigate } from 'react-router-dom';
import { LogOut, Moon, Sun, ShieldCheck, Menu, X, User } from 'lucide-react';
import { useAuth } from '@/lib/auth';
import { request } from '@/lib/api/client';
import { SidebarNav } from './SidebarNav';
import { useLiveNotifications } from '@/hooks/useLiveNotifications';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';

import { LanguageSwitcher } from '@/components/common/LanguageSwitcher';
import { AdminNotificationDropdown } from '../notifications/AdminNotificationDropdown';

export const AdminLayout: React.FC = () => {
  const { user, logout, theme, toggleTheme } = useAuth();
  const { isConnected } = useLiveNotifications();
  const navigate = useNavigate();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [branding, setBranding] = useState<any>(null);

  useEffect(() => {
    request<any>('/admin/settings/branding').then(data => {
       if (data) {
         setBranding(data);
         if (data.primary_color) {
            document.documentElement.style.setProperty('--primary', hexToHSL(data.primary_color));
         }
       }
    }).catch(() => {});
  }, []);

  // Helper to convert hex to HSL for shadcn compatibility
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

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="flex h-screen w-full bg-background overflow-hidden">
      {/* Sidebar Desktop */}
      <aside className="hidden md:flex flex-col w-64 border-r bg-card/60 backdrop-blur-xl shrink-0">
        <div className="h-16 flex items-center px-6 border-b gap-3">
          <div className="h-9 w-9 rounded-xl bg-gradient-to-tr from-primary to-indigo-500 flex items-center justify-center text-white font-bold shadow-lg shadow-primary/25">
            <ShieldCheck className="h-5 w-5" />
          </div>
          <div>
            <h1 className="font-bold text-sm tracking-tight flex items-center gap-1.5">
              FOSSBilling
              <Badge variant="outline" className="text-[10px] px-1.5 py-0 font-mono">
                Next-Gen
              </Badge>
            </h1>
            <p className="text-[11px] text-muted-foreground font-medium">Administrator Portal</p>
          </div>
        </div>

        <nav className="flex-1 overflow-y-auto px-3 py-4 custom-scrollbar">
          <SidebarNav />
        </nav>

        <div className="p-4 border-t bg-muted/20">
          <Link to="/profile" className="flex items-center gap-3 p-1 rounded-lg hover:bg-muted transition-colors group">
            <Avatar className="h-9 w-9 border border-border">
              <AvatarFallback className="bg-primary/10 text-primary font-bold text-xs">
                {user?.name?.slice(0, 2).toUpperCase() || 'AD'}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="text-xs font-semibold truncate leading-tight group-hover:text-primary transition-colors">{user?.name || 'Administrator'}</p>
              <p className="text-[10px] text-muted-foreground truncate">{user?.email || 'admin@fossbilling.org'}</p>
            </div>
            <button
              onClick={(e) => {
                 e.preventDefault();
                 handleLogout();
              }}
              className="h-8 w-8 text-muted-foreground hover:text-destructive flex items-center justify-center rounded-md"
              title="Logout"
            >
              <LogOut className="h-4 w-4" />
            </button>
          </Link>
        </div>
      </aside>

      {/* Mobile Drawer */}
      {mobileOpen && (
        <div
          className="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm md:hidden"
          onClick={() => setMobileOpen(false)}
        >
          <div
            className="w-64 h-full bg-card border-r flex flex-col"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="h-16 flex items-center justify-between px-6 border-b">
              <span className="font-bold text-sm">FOSSBilling Admin</span>
              <Button variant="ghost" size="icon" onClick={() => setMobileOpen(false)}>
                <X className="h-5 w-5" />
              </Button>
            </div>
            <nav className="flex-1 overflow-y-auto px-3 py-4">
              <SidebarNav onItemClick={() => setMobileOpen(false)} />
            </nav>
          </div>
        </div>
      )}

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col h-full overflow-hidden">
        <header className="h-16 border-b bg-card/40 backdrop-blur-md px-6 flex items-center justify-between shrink-0">
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={() => setMobileOpen(true)}
            >
              <Menu className="h-5 w-5" />
            </Button>
            <div className="flex items-center gap-2">
              <span className={`inline-block h-2 w-2 rounded-full animate-pulse ${isConnected ? 'bg-emerald-500' : 'bg-amber-500'}`} />
              <span className="text-xs font-medium text-muted-foreground">
                {isConnected ? 'Admin API Engine: Live' : 'Admin API Engine: Offline'}
              </span>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <LanguageSwitcher />

            <AdminNotificationDropdown />

            <Button
              variant="outline"
              size="icon"
              className="h-9 w-9"
              onClick={toggleTheme}
              title={theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
            >
              {theme === 'dark' ? <Sun className="h-4 w-4 text-amber-400" /> : <Moon className="h-4 w-4 text-slate-700" />}
            </Button>

            <a
              href="http://localhost:8080/docs"
              target="_blank"
              rel="noreferrer"
              className="inline-flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 rounded-lg border bg-secondary/60 hover:bg-secondary transition-colors"
            >
              <span>API Docs</span>
              <span className="text-[10px] text-muted-foreground">↗</span>
            </a>
          </div>
        </header>

        <main className="flex-1 overflow-y-auto p-6 md:p-8">
          <div className="max-w-7xl mx-auto space-y-6">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
};

export default AdminLayout;
