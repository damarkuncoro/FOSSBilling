import { useQuery } from '@tanstack/react-query';
import { adminStatsService } from '@/services/admin_stats.service';
import { adminSystemService } from '@/services/admin_system.service';
import type { DashboardStats, AuditLog } from '@/types/api';

export function useDashboard() {
  const { data: stats, isLoading: statsLoading, refetch: fetchStats, error: statsError } = useQuery({
    queryKey: ['admin', 'dashboard', 'stats'],
    queryFn: () => adminStatsService.getDashboardStats(),
    refetchInterval: 60000,
  });

  const { data: recentLogs = [], isLoading: logsLoading } = useQuery({
    queryKey: ['admin', 'dashboard', 'recent-logs'],
    queryFn: async () => {
      const logs = await adminSystemService.listAuditLogs(5, 0);
      return (logs as any).data || logs; // handle potential wrapper
    },
  });

  const { data: systemStatus, isLoading: statusLoading } = useQuery({
    queryKey: ['admin', 'dashboard', 'system-status'],
    queryFn: () => adminSystemService.getSystemStatus(),
  });

  const { data: activityTrend = {}, isLoading: trendLoading } = useQuery({
    queryKey: ['admin', 'dashboard', 'activity-trend'],
    queryFn: () => adminSystemService.getActivityTrend(7),
  });

  const { data: revenueProjections = {}, isLoading: projectionsLoading } = useQuery({
    queryKey: ['admin', 'dashboard', 'projections'],
    queryFn: () => adminStatsService.getRevenueProjections(),
  });

  return {
    stats: stats as DashboardStats,
    recentLogs: recentLogs as AuditLog[],
    systemStatus,
    activityTrend,
    revenueProjections,
    loading: statsLoading || logsLoading || statusLoading || trendLoading || projectionsLoading,
    error: statsError ? (statsError as Error).message : null,
    fetchStats,
    revenueTrends: stats?.revenue_trends || [],
  };
}
