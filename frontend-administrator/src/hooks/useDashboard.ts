import { useState, useEffect } from 'react';
import { adminStatsService } from '@/services/admin_stats.service';
import { adminSystemService } from '@/services/admin_system.service';
import type { DashboardStats } from '@/types/api';

export function useDashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [recentLogs, setRecentLogs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = async () => {
    setLoading(true);
    setError(null);
    try {
      const [statsData, logsData] = await Promise.all([
        adminStatsService.getDashboardMetrics(),
        adminSystemService.getAuditLogs(5, 0), // Get top 5 recent logs
      ]);
      setStats(statsData);
      setRecentLogs(logsData || []);
    } catch (err: any) {
      setError(err.message || 'Failed to load dashboard metrics');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, []);

  const revenueTrends = stats?.revenue_trends && stats.revenue_trends.length > 0
    ? stats.revenue_trends
    : [];

  return {
    stats,
    recentLogs,
    loading,
    error,
    fetchStats,
    revenueTrends,
  };
}
