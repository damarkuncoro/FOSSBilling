import React from 'react';
import { Trophy, Users, DollarSign, ArrowUpRight, Search, Settings2 } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { formatMoney } from '@/lib/utils';
import { useQuery } from '@tanstack/react-query';
import { adminStatsService } from '@/services/admin_stats.service';

export const Affiliates: React.FC = () => {
  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ['admin', 'affiliates', 'stats'],
    queryFn: () => adminStatsService.getDashboardStats(), // We can extend this for affiliate specific stats
  });

  return (
    <div className="space-y-8 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Affiliate Program Management</h1>
          <p className="text-sm text-muted-foreground">
            Manage your partners, track referrals, and handle commission payouts.
          </p>
        </div>
        <Button className="gap-2">
          <Settings2 className="h-4 w-4" /> Global Affiliate Settings
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card className="border-border/60">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">Total Partners</CardTitle>
            <Users className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">124</div>
            <p className="text-[10px] text-muted-foreground mt-1">+12 from last month</p>
          </CardContent>
        </Card>
        <Card className="border-border/60">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">Pending Payouts</CardTitle>
            <DollarSign className="h-4 w-4 text-amber-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatMoney(4500000)}</div>
            <p className="text-[10px] text-muted-foreground mt-1">8 requests awaiting approval</p>
          </CardContent>
        </Card>
        <Card className="border-border/60">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-muted-foreground uppercase">Total Commissions Paid</CardTitle>
            <Trophy className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatMoney(82400000)}</div>
            <div className="flex items-center gap-1 text-emerald-500 text-[10px] mt-1 font-bold uppercase">
               <ArrowUpRight className="h-3 w-3" /> lifetime payout
            </div>
          </CardContent>
        </Card>
      </div>

      <Card className="border-border/60">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-base">Top Performing Affiliates</CardTitle>
              <CardDescription>Most active referrers in the last 30 days</CardDescription>
            </div>
            <div className="relative w-64">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input placeholder="Search partners..." className="pl-9 h-9 text-xs" />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Partner Name</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Total Referrals</TableHead>
                <TableHead>Current Balance</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell>
                   <div className="font-medium">Andi Wijaya</div>
                   <div className="text-[10px] text-muted-foreground">andi@wijaya-tech.com</div>
                </TableCell>
                <TableCell><Badge variant="success">Active</Badge></TableCell>
                <TableCell>42 Clients</TableCell>
                <TableCell className="font-mono font-bold">{formatMoney(1200000)}</TableCell>
                <TableCell className="text-right">
                  <Button variant="ghost" size="sm">Manage</Button>
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>
                   <div className="font-medium">Dewi Lestari</div>
                   <div className="text-[10px] text-muted-foreground">dewi.l@bloghost.id</div>
                </TableCell>
                <TableCell><Badge variant="success">Active</Badge></TableCell>
                <TableCell>28 Clients</TableCell>
                <TableCell className="font-mono font-bold">{formatMoney(850000)}</TableCell>
                <TableCell className="text-right">
                  <Button variant="ghost" size="sm">Manage</Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default Affiliates;
