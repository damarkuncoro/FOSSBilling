import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';

export function useAuditLogs() {
  const { data: logs = [], isLoading: loading, refetch: fetchLogs } = useQuery({
    queryKey: ['admin', 'audit-logs'],
    queryFn: () => api.getAuditLogs(),
  });

  return { logs, loading, fetchLogs };
}
