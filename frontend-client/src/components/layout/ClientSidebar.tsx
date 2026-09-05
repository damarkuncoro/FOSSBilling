import React from 'react';
import { NavLink, Link } from 'react-router-dom';
import { Shield, Wallet, LogOut, LogIn, Plus } from 'lucide-react';
import { useClientAuth } from '@/lib/auth';
import { clientNavGroups } from './clientNavConfig';
import { formatMoney } from '@/lib/utils';
import { Button } from '@/components/ui/button';

interface ClientSidebarProps {
  onItemClick?: () => void;
  onOpenDeposit?: () => void;
}

export const ClientSidebar: React.FC<ClientSidebarProps> = ({ onItemClick, onOpenDeposit }) => {
  const { user, balance, isAuthenticated, logout } = useClientAuth();

  return (
    <aside className="w-64 h-full bg-card border-r border-border/70 flex flex-col justify-between shrink-0 select-none">
      {/* Top Brand Header */}
      <div>
        <div className="h-16 flex items-center px-6 border-b border-border/50 gap-3">
          <Link to="/" className="flex items-center gap-2.5 group" onClick={onItemClick}>
            <div className="h-8 w-8 rounded-lg bg-gradient-to-tr from-primary to-indigo-500 flex items-center justify-center text-white shadow-md shadow-primary/20 group-hover:scale-105 transition-transform">
              <Shield className="h-4 w-4" />
            </div>
            <div>
              <span className="font-bold text-sm tracking-tight block leading-none">FOSSBilling</span>
              <span className="text-[10px] text-muted-foreground font-semibold">Client Portal</span>
            </div>
          </Link>
        </div>

        {/* Balance Card Widget */}
        {isAuthenticated && (
          <div className="p-3 mx-3 my-3 rounded-xl bg-primary/5 border border-primary/15">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-1.5 text-xs text-muted-foreground font-medium">
                <Wallet className="h-3.5 w-3.5 text-primary" />
                <span>Account Balance</span>
              </div>
              {onOpenDeposit && (
                <Button variant="ghost" size="sm" onClick={onOpenDeposit} className="h-5 px-1.5 text-[10px] text-primary gap-0.5 font-semibold">
                  <Plus className="h-3 w-3" /> Add
                </Button>
              )}
            </div>
            <p className="text-base font-bold font-mono text-primary mt-1">
              {formatMoney(balance || 0, user?.currency || 'USD')}
            </p>
          </div>
        )}

        {/* Navigation Groups */}
        <div className="p-3 space-y-4 overflow-y-auto max-h-[calc(100vh-230px)]">
          {clientNavGroups.map((group, gIdx) => {
            const visibleItems = group.items.filter((it) => !it.requiresAuth || isAuthenticated);
            if (visibleItems.length === 0) return null;

            return (
              <div key={gIdx} className="space-y-1">
                {group.groupName && (
                  <p className="px-3 text-[10px] font-bold uppercase tracking-wider text-muted-foreground/70 mb-1">
                    {group.groupName}
                  </p>
                )}
                {visibleItems.map((item) => {
                  const Icon = item.icon;
                  return (
                    <NavLink
                      key={item.path}
                      to={item.path}
                      end={item.path === '/' || item.path === '/dashboard'}
                      onClick={onItemClick}
                      className={({ isActive }) =>
                        `flex items-center gap-3 px-3 py-2 rounded-lg text-xs font-medium transition-all ${
                          isActive
                            ? 'bg-primary text-primary-foreground shadow-sm shadow-primary/20 font-semibold'
                            : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                        }`
                      }
                    >
                      <Icon className="h-4 w-4 shrink-0" />
                      <span>{item.name}</span>
                    </NavLink>
                  );
                })}
              </div>
            );
          })}
        </div>
      </div>

      {/* User Footer Action */}
      <div className="p-3 border-t border-border/50">
        {isAuthenticated && user ? (
          <div className="flex items-center justify-between p-2 rounded-lg bg-muted/40 border border-border/40">
            <div className="flex items-center gap-2 overflow-hidden">
              <div className="h-7 w-7 rounded-full bg-primary/10 text-primary font-bold flex items-center justify-center text-xs shrink-0">
                {user.first_name?.[0] || 'C'}
              </div>
              <div className="truncate">
                <p className="text-xs font-semibold truncate leading-tight">{user.first_name} {user.last_name}</p>
                <p className="text-[10px] text-muted-foreground truncate">{user.email}</p>
              </div>
            </div>
            <Button variant="ghost" size="icon" onClick={logout} className="h-7 w-7 text-muted-foreground hover:text-destructive shrink-0" title="Logout">
              <LogOut className="h-3.5 w-3.5" />
            </Button>
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-2">
            <Link to="/login" onClick={onItemClick} className="w-full">
              <Button variant="outline" size="sm" className="w-full text-xs gap-1.5 h-8">
                <LogIn className="h-3.5 w-3.5" /> Login
              </Button>
            </Link>
            <Link to="/register" onClick={onItemClick} className="w-full">
              <Button size="sm" className="w-full text-xs h-8">Register</Button>
            </Link>
          </div>
        )}
      </div>
    </aside>
  );
};
