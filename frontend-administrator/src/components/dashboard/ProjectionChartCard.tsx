import React from 'react';
import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
} from 'recharts';
import { formatMoney } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface ProjectionChartCardProps {
  data: Record<string, number>;
}

export const ProjectionChartCard: React.FC<ProjectionChartCardProps> = ({ data }) => {
  // Convert map to cumulative chart data
  const sortedDays = Object.keys(data).sort();
  let cumulative = 0;
  const chartData = sortedDays.map(day => {
    cumulative += data[day];
    return {
      day: day.split('-').slice(2).join('/'), // DD
      amount: data[day],
      cumulative,
    };
  });

  return (
    <Card className="border-border/60 shadow-sm h-full">
      <CardHeader>
        <CardTitle className="text-base font-semibold text-amber-600">30-Day Cash Flow Projection</CardTitle>
        <CardDescription>Expected incoming revenue from automated service renewals</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="h-[200px] w-full mt-2">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} opacity={0.1} />
              <XAxis
                dataKey="day"
                axisLine={false}
                tickLine={false}
                fontSize={10}
                tick={{ fill: '#888' }}
              />
              <YAxis
                axisLine={false}
                tickLine={false}
                fontSize={10}
                tick={{ fill: '#888' }}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: 'var(--card)',
                  borderColor: 'var(--border)',
                  borderRadius: '12px',
                  fontSize: '11px',
                }}
                formatter={(val: any) => [formatMoney(val), 'Income']}
              />
              <Line
                type="stepAfter"
                dataKey="cumulative"
                stroke="#f59e0b"
                strokeWidth={3}
                dot={false}
                activeDot={{ r: 4, strokeWidth: 0 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
        <div className="mt-4 flex justify-between items-center text-[11px]">
           <span className="text-muted-foreground uppercase font-bold tracking-tighter">Projected Total</span>
           <span className="font-mono font-bold text-amber-600">{formatMoney(cumulative)}</span>
        </div>
      </CardContent>
    </Card>
  );
};
