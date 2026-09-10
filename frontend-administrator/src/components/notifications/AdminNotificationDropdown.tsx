import React, { useState } from 'react';
import { Bell, CheckCheck, Info, AlertTriangle, AlertCircle, ShoppingCart, LifeBuoy } from 'lucide-react';
import { useAdminNotifications } from '../../hooks/useAdminNotifications';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { ScrollArea } from '@/components/ui/scroll-area';

export const AdminNotificationDropdown: React.FC = () => {
  const { notifications, unreadCount, markAsRead, markAllRead } = useAdminNotifications();
  const [open, setOpen] = useState(false);

  const getIcon = (type: string, module: string) => {
    if (module === 'support') return <LifeBuoy className="h-4 w-4 text-rose-500" />;
    if (module === 'billing') return <ShoppingCart className="h-4 w-4 text-emerald-500" />;

    switch (type) {
      case 'success': return <Info className="h-4 w-4 text-emerald-500" />;
      case 'warning': return <AlertTriangle className="h-4 w-4 text-amber-500" />;
      case 'danger': return <AlertCircle className="h-4 w-4 text-rose-500" />;
      default: return <Info className="h-4 w-4 text-blue-500" />;
    }
  };

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="relative h-9 w-9 text-muted-foreground hover:text-foreground">
          <Bell className="h-4.5 w-4.5" />
          {unreadCount > 0 && (
            <span className="absolute top-1 right-1 flex h-4 w-4 items-center justify-center rounded-full bg-rose-500 text-[10px] font-bold text-white ring-2 ring-background">
              {unreadCount > 9 ? '9+' : unreadCount}
            </span>
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-80 p-0 shadow-2xl rounded-2xl overflow-hidden">
        <div className="p-4 bg-muted/30 border-b flex items-center justify-between">
           <h3 className="text-sm font-black tracking-tight">System Alerts</h3>
           {unreadCount > 0 && (
             <button
               onClick={() => markAllRead()}
               className="text-[10px] font-bold text-indigo-600 hover:text-indigo-700 flex items-center gap-1 uppercase tracking-wider"
             >
                <CheckCheck className="h-3 w-3" /> Mark all read
             </button>
           )}
        </div>

        <ScrollArea className="h-[350px]">
          {!notifications || notifications.length === 0 ? (
            <div className="py-20 text-center space-y-2">
               <Bell className="h-8 w-8 text-muted-foreground/20 mx-auto" />
               <p className="text-xs text-muted-foreground font-medium">No recent alerts.</p>
            </div>
          ) : (
            notifications.map((n) => (
              <DropdownMenuItem
                key={n.id}
                className={`flex flex-col items-start gap-1 p-4 cursor-default focus:bg-muted/50 border-b border-border/40 last:border-0 ${!n.is_read ? 'bg-indigo-50/20' : ''}`}
                onSelect={(e: any) => {
                  e.preventDefault();
                  if (!n.is_read) markAsRead(n.id);
                }}
              >
                <div className="flex items-center gap-2 w-full">
                  <div className="shrink-0 p-1.5 rounded-lg bg-background border border-border/50">
                    {getIcon(n.type, n.module)}
                  </div>
                  <span className="font-bold text-xs truncate flex-1">{n.title}</span>
                  {!n.is_read && <div className="h-1.5 w-1.5 rounded-full bg-indigo-600" />}
                </div>
                <p className="text-[11px] text-muted-foreground leading-relaxed pl-8">
                  {n.message}
                </p>
                <span className="text-[9px] text-muted-foreground/60 font-medium pl-8 mt-0.5">
                  {new Date(n.created_at).toLocaleString()}
                </span>
              </DropdownMenuItem>
            ))
          )}
        </ScrollArea>

        <div className="p-2.5 bg-muted/10 border-t">
           <Button variant="ghost" size="sm" className="w-full text-[11px] h-8 font-semibold text-muted-foreground hover:text-foreground">
              View All History
           </Button>
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
