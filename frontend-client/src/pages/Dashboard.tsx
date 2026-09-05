import React from 'react';
import { Link } from 'react-router-dom';
import { Zap, AlertTriangle, ArrowRight } from 'lucide-react';
import { useClientDashboard } from '@/hooks/useClientDashboard';
import { Button } from '@/components/ui/button';
import { DashboardStatsGrid } from '@/components/dashboard/DashboardStatsGrid';
import { ActiveServicesList } from '@/components/dashboard/ActiveServicesList';
import { RecentInvoicesList } from '@/components/dashboard/RecentInvoicesList';

export const Dashboard: React.FC = () => {
  const {
    user,
    balance,
    orders,
    invoices,
    tickets,
    unpaidInvoices,
    activeOrders,
  } = useClientDashboard();

  return (
    <div className="space-y-8 animate-in fade-in-50 duration-300">
      {/* Welcome Banner */}
      <div className="p-6 md:p-8 rounded-2xl bg-gradient-to-r from-primary/10 via-primary/5 to-transparent border border-primary/20 flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
        <div className="space-y-1">
          <h1 className="text-2xl md:text-3xl font-bold tracking-tight">
            Welcome back, {user?.first_name} {user?.last_name}! 👋
          </h1>
          <p className="text-sm text-muted-foreground">
            Manage your cloud servers, track invoices, and access support tickets.
          </p>
        </div>
        <Link to="/">
          <Button className="gap-2 font-semibold shadow-md shadow-primary/20">
            <Zap className="h-4 w-4" />
            Order New Cloud Service
          </Button>
        </Link>
      </div>

      {/* Unpaid Invoice Alert Banner */}
      {unpaidInvoices.length > 0 && (
        <div className="p-4 rounded-xl bg-amber-500/10 border border-amber-500/30 text-amber-900 dark:text-amber-200 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 animate-in fade-in">
          <div className="flex items-center gap-3">
            <AlertTriangle className="h-5 w-5 text-amber-500 shrink-0" />
            <div>
              <p className="text-sm font-semibold">
                You have {unpaidInvoices.length} unpaid invoice pending payment!
              </p>
              <p className="text-xs opacity-90">
                Please settle your invoice to avoid service interruption.
              </p>
            </div>
          </div>
          <Link to="/invoices">
            <Button size="sm" variant="default" className="bg-amber-600 hover:bg-amber-700 text-white gap-1 text-xs">
              View & Pay Now <ArrowRight className="h-3.5 w-3.5" />
            </Button>
          </Link>
        </div>
      )}

      {/* Stats Cards */}
      <DashboardStatsGrid
        user={user}
        balance={balance}
        orders={orders}
        invoices={invoices}
        tickets={tickets}
        unpaidInvoices={unpaidInvoices}
        activeOrders={activeOrders}
      />

      {/* Services & Invoices Quick View */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ActiveServicesList orders={orders} />
        <RecentInvoicesList invoices={invoices} />
      </div>
    </div>
  );
};
