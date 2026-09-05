import React, { useState } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import { LogOut, Moon, Sun, ShieldCheck, Menu, X } from 'lucide-react';
import { useAuth } from '@/lib/auth';
import { SidebarNav } from './SidebarNav';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';

export const AdminLayout: React.FC = () => {
  const { user, logout, theme, toggleTheme } = useAuth();
  const navigate = useNavigate();
  const [mobileOpen, setMobileOpen] = useState(false);

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
          <div className="flex items-center gap-3">
            <Avatar className="h-9 w-9 border border-border">
              <AvatarFallback className="bg-primary/10 text-primary font-bold text-xs">
                {user?.name?.slice(0, 2).toUpperCase() || 'AD'}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="text-xs font-semibold truncate leading-tight">{user?.name || 'Administrator'}</p>
              <p className="text-[10px] text-muted-foreground truncate">{user?.email || 'admin@fossbilling.org'}</p>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={handleLogout}
              className="h-8 w-8 text-muted-foreground hover:text-destructive shrink-0"
              title="Logout"
            >
              <LogOut className="h-4 w-4" />
            </Button>
          </div>
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
              <span className="inline-block h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
              <span className="text-xs font-medium text-muted-foreground">Admin API Engine: Active</span>
            </div>
          </div>

          <div className="flex items-center gap-3">
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
