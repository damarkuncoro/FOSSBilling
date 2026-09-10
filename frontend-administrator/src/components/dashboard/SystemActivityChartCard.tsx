import React from 'react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  Cell,
} from 'recharts';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface SystemActivityChartCardProps {
  data: Record<string, number>;
}

export const SystemActivityChartCard: React.FC<SystemActivityChartCardProps> = ({ data }) => {
  // Convert object map to array for recharts
  const chartData = Object.entries(data).map(([day, count]) => ({
    day: day.split('-').slice(1).join('/'), // format as MM/DD
    count,
  }));

  return (
    <Card className="border-border/60 shadow-sm h-full">
      <CardHeader>
        <CardTitle className="text-base font-semibold">System Traffic & Activity</CardTitle>
        <CardDescription>Daily administrative and system events (7 days)</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="h-[200px] w-full mt-2">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData}>
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
                allowDecimals={false}
              />
              <Tooltip
                cursor={{ fill: 'rgba(99, 102, 241, 0.05)' }}
                contentStyle={{
                  backgroundColor: 'var(--card)',
                  borderColor: 'var(--border)',
                  borderRadius: '8px',
                  fontSize: '11px',
                }}
              />
              <Bar dataKey="count" radius={[4, 4, 0, 0]}>
                {chartData.map((entry, index) => (
                  <Cell
                    key={`cell-${index}`}
                    fill={index === chartData.length - 1 ? '#6366f1' : '#a5b4fc'}
                  />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  );
};
