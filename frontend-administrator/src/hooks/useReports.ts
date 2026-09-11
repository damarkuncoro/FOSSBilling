import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { FinancialReportSummary } from '@/types/modules';

export function useReports() {
  const { data: report = null, isLoading: loading, refetch: fetchReports } = useQuery({
    queryKey: ['admin', 'reports', 'financial'],
    queryFn: () => api.getFinancialReports(),
  });

  const handleExportCsv = () => {
    if (!report) return;
    let csv = 'Month,Revenue,Tax,Invoices\n';
    report.monthly_breakdown.forEach((row) => {
      csv += `${row.month},${row.revenue},${row.tax},${row.invoices_count}\n`;
    });
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `fossbilling-financial-report-${new Date().toISOString().slice(0, 10)}.csv`;
    a.click();
  };

  return {
    report,
    loading,
    fetchReports,
    handleExportCsv,
  };
}
