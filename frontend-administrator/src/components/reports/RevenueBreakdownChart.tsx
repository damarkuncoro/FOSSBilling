import React from 'react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  Legend,
} from 'recharts';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { formatMoney } from '@/lib/utils';
import { FinancialReportSummary } from '@/types/modules';

interface RevenueBreakdownChartProps {
  report: FinancialReportSummary | null;
}

export const RevenueBreakdownChart: React.FC<RevenueBreakdownChartProps> = ({ report }) => {
  const data = report?.monthly_breakdown || [];

  return (
    <Card className="border-border/60 shadow-sm h-full">
      <CardHeader>
        <CardTitle className="text-base font-semibold">Revenue & Tax Breakdown</CardTitle>
        <CardDescription>Comparison of monthly gross revenue vs collected taxes</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="h-[300px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={data} margin={{ top: 10, right: 10, left: -10, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} opacity={0.1} />
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
                tickFormatter={(v) => `$${v >= 1000 ? v/1000 + 'k' : v}`}
              />
              <Tooltip
                cursor={{ fill: 'rgba(99, 102, 241, 0.05)' }}
                contentStyle={{
                  backgroundColor: 'var(--card)',
                  borderColor: 'var(--border)',
                  borderRadius: '12px',
                  fontSize: '12px',
                  boxShadow: '0 10px 15px -3px rgba(0, 0, 0, 0.1)',
                }}
                formatter={(val: any) => [formatMoney(val), '']}
              />
              <Legend
                verticalAlign="top"
                align="right"
                iconType="circle"
                iconSize={8}
                wrapperStyle={{ fontSize: '11px', paddingBottom: '20px' }}
              />
              <Bar
                name="Gross Revenue"
                dataKey="revenue"
                fill="#6366f1"
                radius={[4, 4, 0, 0]}
                barSize={32}
              />
              <Bar
                name="Tax Collected"
                dataKey="tax"
                fill="#10b981"
                radius={[4, 4, 0, 0]}
                barSize={32}
              />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  );
};
