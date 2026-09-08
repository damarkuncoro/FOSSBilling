import React from 'react';
import { RefreshCw, Download } from 'lucide-react';
import { useReports } from '@/hooks/useReports';
import { Button } from '@/components/ui/button';
import { RevenueMetricsGrid } from '@/components/reports/RevenueMetricsGrid';
import { TaxReportCard } from '@/components/reports/TaxReportCard';

export const Reports: React.FC = () => {
  const { report, loading, fetchReports, handleExportCsv } = useReports();

  const handleDetailedExport = () => {
    const token = localStorage.getItem('fossbilling_admin_token') || sessionStorage.getItem('fossbilling_admin_token') || '';
    const url = `/api/v1/admin/reports/invoices/csv?token=${encodeURIComponent(token)}`;
    window.open(url, '_blank');
  };

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Financial Reports & Tax Analytics</h1>
          <p className="text-sm text-muted-foreground">
            Monthly recurring revenue, fiscal VAT tax summaries, and accounting export tools.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={handleDetailedExport} className="gap-2">
            <Download className="h-4 w-4" />
            Detailed Export (CSV)
          </Button>
          <Button variant="outline" size="sm" onClick={fetchReports} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      <RevenueMetricsGrid report={report} />

      <TaxReportCard report={report} onExport={handleExportCsv} />
    </div>
  );
};
