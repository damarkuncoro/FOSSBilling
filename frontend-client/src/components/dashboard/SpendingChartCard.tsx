import React from 'react';
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
} from 'recharts';
import { formatMoney } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface SpendingChartCardProps {
  data: any[];
  currency?: string;
}

export const SpendingChartCard: React.FC<SpendingChartCardProps> = ({ data, currency = 'USD' }) => {
  return (
    <Card className="border-border/60 shadow-sm overflow-hidden">
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-semibold">Spending Trends</CardTitle>
        <CardDescription>Your monthly cloud infrastructure investment (last 6 months)</CardDescription>
      </CardHeader>
      <CardContent className="pt-4">
        <div className="h-[220px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
              <defs>
                <linearGradient id="colorAmt" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="hsl(var(--primary))" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity={0.0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" opacity={0.1} vertical={false} />
              <XAxis
                dataKey="month"
                axisLine={false}
                tickLine={false}
                fontSize={11}
                tick={{ fill: 'hsl(var(--muted-foreground))' }}
              />
              <YAxis
                axisLine={false}
                tickLine={false}
                fontSize={11}
                tick={{ fill: 'hsl(var(--muted-foreground))' }}
                tickFormatter={(v) => `$${v}`}
              />
              <Tooltip
                formatter={(val: any) => [formatMoney(val, currency), 'Total Paid']}
                contentStyle={{
                  backgroundColor: 'var(--card)',
                  borderColor: 'var(--border)',
                  borderRadius: '12px',
                  fontSize: '11px',
                }}
              />
              <Area
                type="monotone"
                dataKey="amount"
                stroke="hsl(var(--primary))"
                strokeWidth={2.5}
                fillOpacity={1}
                fill="url(#colorAmt)"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  );
};
