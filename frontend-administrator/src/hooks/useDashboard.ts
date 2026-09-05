import { useState, useEffect } from 'react';
import { api } from '@/lib/api';

export const mockRevenueTrends = [
  { month: 'Jan', revenue: 12400, mrr: 8200 },
  { month: 'Feb', revenue: 15800, mrr: 9400 },
  { month: 'Mar', revenue: 19200, mrr: 11000 },
  { month: 'Apr', revenue: 24500, mrr: 14500 },
  { month: 'May', revenue: 31000, mrr: 18200 },
  { month: 'Jun', revenue: 42000, mrr: 23500 },
  { month: 'Jul', revenue: 56000, mrr: 29000 },
];

export function useDashboard() {
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await api.getDashboardStats();
      setStats(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load dashboard metrics');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, []);

  return {
    stats,
    loading,
    error,
    fetchStats,
    mockRevenueTrends,
  };
}
