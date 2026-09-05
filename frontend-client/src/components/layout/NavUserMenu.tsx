import React from 'react';
import { Link } from 'react-router-dom';
import { Wallet, LogOut } from 'lucide-react';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';
import { ClientProfile } from '@/types/api';

interface NavUserMenuProps {
  user: ClientProfile | null;
  balance: number;
  isAuthenticated: boolean;
  onLogout: () => void;
}

export const NavUserMenu: React.FC<NavUserMenuProps> = ({
  user,
  balance,
  isAuthenticated,
  onLogout,
}) => {
  if (!isAuthenticated) {
    return (
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
    );
  }

  return (
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

      <Button
        variant="ghost"
        size="icon"
        onClick={onLogout}
        className="h-8 w-8 text-muted-foreground hover:text-destructive"
        title="Logout"
      >
        <LogOut className="h-4 w-4" />
      </Button>
    </div>
  );
};
