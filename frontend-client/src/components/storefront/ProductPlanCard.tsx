import React from 'react';
import { Server, Cpu, Key, Download, Check, Zap } from 'lucide-react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatMoney } from '@/lib/utils';
import { HostingPlan } from '@/types/api';

interface ProductPlanCardProps {
  plan: HostingPlan;
  onAddToCart: (plan: HostingPlan) => void;
}

export const ProductPlanCard: React.FC<ProductPlanCardProps> = ({ plan, onAddToCart }) => {
  const getIcon = () => {
    switch (plan.type) {
      case 'hosting':
        return Server;
      case 'vps':
        return Cpu;
      case 'license':
        return Key;
      case 'downloadable':
      default:
        return Download;
    }
  };

  const Icon = getIcon();

  return (
    <Card
      className={`relative flex flex-col border transition-all hover:shadow-xl hover:-translate-y-1 ${
        plan.popular ? 'border-primary shadow-lg shadow-primary/10 bg-primary/[0.02]' : 'border-border/70'
      }`}
    >
      {plan.popular && (
        <div className="absolute -top-3 left-1/2 -translate-x-1/2">
          <Badge className="bg-primary text-primary-foreground text-xs uppercase px-3 font-bold">
            Most Popular
          </Badge>
        </div>
      )}
      <CardHeader>
        <div className="h-10 w-10 rounded-xl bg-primary/10 text-primary flex items-center justify-center mb-2">
          <Icon className="h-5 w-5" />
        </div>
        <CardTitle className="text-lg">{plan.title}</CardTitle>
        <CardDescription className="text-xs min-h-[32px]">{plan.description}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 space-y-4">
        <div>
          <span className="text-3xl font-extrabold">{formatMoney(plan.price)}</span>
          <span className="text-xs text-muted-foreground font-medium">/{plan.period.toLowerCase()}</span>
        </div>

        <ul className="space-y-2 text-xs text-muted-foreground pt-2 border-t">
          {plan.features.map((feat, idx) => (
            <li key={idx} className="flex items-center gap-2">
              <Check className="h-3.5 w-3.5 text-emerald-500 shrink-0" />
              <span>{feat}</span>
            </li>
          ))}
        </ul>
      </CardContent>
      <CardFooter>
        <Button
          className="w-full font-semibold gap-1.5"
          variant={plan.popular ? 'default' : 'outline'}
          onClick={() => onAddToCart(plan)}
        >
          <Zap className="h-4 w-4" />
          Order Now
        </Button>
      </CardFooter>
    </Card>
  );
};
