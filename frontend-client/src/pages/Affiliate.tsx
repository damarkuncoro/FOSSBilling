import React from 'react';
import { Users, DollarSign, ExternalLink, Trophy, ArrowUpRight, CheckCircle2 } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { formatMoney } from '@/lib/utils';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';

export const Affiliate: React.FC = () => {
  const queryClient = useQueryClient();
  const { data: aff, isLoading, error } = useQuery({
    queryKey: ['client', 'affiliate'],
    queryFn: () => api.getAffiliate(),
    retry: false,
  });

  const joinMutation = useMutation({
    mutationFn: () => api.joinAffiliate(),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['client', 'affiliate'] }),
  });

  if (isLoading) return <div className="p-8 text-center">Loading affiliate portal...</div>;

  if (error || !aff) {
    return (
      <div className="max-w-2xl mx-auto py-12 text-center space-y-6 animate-in fade-in zoom-in-95">
        <div className="h-20 w-20 bg-indigo-100 text-indigo-600 rounded-full flex items-center justify-center mx-auto mb-4">
           <Trophy className="h-10 w-10" />
        </div>
        <h1 className="text-3xl font-black tracking-tight">Earn Money with FOSSBilling</h1>
        <p className="text-muted-foreground leading-relaxed">
          Join our partner program today and earn <span className="font-bold text-foreground">10% recurring commission</span>
          for every client you refer. There's no limit to how much you can earn!
        </p>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 text-sm py-4">
           <div className="space-y-1">
              <div className="font-bold">Fast Payouts</div>
              <div className="text-xs text-muted-foreground">Monthly via PayPal/Stripe</div>
           </div>
           <div className="space-y-1 border-x">
              <div className="font-bold">Real-time Tracking</div>
              <div className="text-xs text-muted-foreground">Instant referral alerts</div>
           </div>
           <div className="space-y-1">
              <div className="font-bold">90 Day Cookies</div>
              <div className="text-xs text-muted-foreground">Long conversion window</div>
           </div>
        </div>
        <Button
          size="lg"
          onClick={() => joinMutation.mutate()}
          disabled={joinMutation.isPending}
          className="px-12 rounded-full font-bold"
        >
          {joinMutation.isPending ? 'Joining...' : 'Join Affiliate Program'}
        </Button>
      </div>
    );
  }

  const referralLink = `${window.location.origin}/r/${aff.client_id}`;

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Partner Dashboard</h1>
        <p className="text-sm text-muted-foreground">Grow your earnings by referring new customers to our platform.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card className="bg-primary text-primary-foreground border-none shadow-lg overflow-hidden relative">
          <div className="absolute top-0 right-0 p-4 opacity-10">
             <DollarSign className="h-24 w-24" />
          </div>
          <CardHeader className="pb-2">
            <CardTitle className="text-xs uppercase tracking-widest opacity-80">Available Balance</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black">{formatMoney(aff.balance)}</div>
            <Button variant="secondary" size="sm" className="mt-4 w-full font-bold bg-white text-primary hover:bg-indigo-50">
              Request Payout
            </Button>
          </CardContent>
        </Card>

        <Card className="border-border/60">
          <CardHeader className="pb-2">
            <CardTitle className="text-xs uppercase tracking-widest text-muted-foreground">Total Earned</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black">{formatMoney(aff.total_earned)}</div>
            <div className="flex items-center gap-1 text-emerald-500 text-xs mt-2 font-bold uppercase">
               <ArrowUpRight className="h-3 w-3" /> lifetime revenue
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/60">
          <CardHeader className="pb-2">
            <CardTitle className="text-xs uppercase tracking-widest text-muted-foreground">Commission Rate</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black">{aff.commission_rate}%</div>
            <div className="mt-3">
               <Progress value={aff.commission_rate * 5} className="h-1.5" />
               <p className="text-[10px] text-muted-foreground mt-2 uppercase font-bold">Standard Tier Partner</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-bold flex items-center gap-2">
            <ExternalLink className="h-4 w-4 text-primary" /> Your Referral Link
          </CardTitle>
          <CardDescription>Share this link with your audience or friends to start earning commissions.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2">
            <div className="flex-1 bg-muted/50 border rounded-lg px-4 py-2.5 font-mono text-sm flex items-center justify-between">
               {referralLink}
               <Badge variant="outline" className="ml-2 bg-background">90-Day Cookie</Badge>
            </div>
            <Button onClick={() => navigator.clipboard.writeText(referralLink)}>Copy Link</Button>
          </div>
        </CardContent>
      </Card>

      <div className="space-y-4">
        <h3 className="text-lg font-bold flex items-center gap-2">
          <Users className="h-5 w-5 text-primary" /> Recent Referrals
        </h3>
        <Card className="border-border/60">
           <CardContent className="p-0">
              <div className="p-8 text-center text-muted-foreground text-sm space-y-2">
                 <CheckCircle2 className="h-8 w-8 mx-auto opacity-10" />
                 <p>You haven't referred any customers yet. Share your link to get started!</p>
              </div>
           </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Affiliate;
