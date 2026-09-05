import React from 'react';
import { FileSpreadsheet, Download } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { formatMoney } from '@/lib/utils';
import { FinancialReportSummary } from '@/types/modules';

interface TaxReportCardProps {
  report: FinancialReportSummary | null;
  onExport: () => void;
}

export const TaxReportCard: React.FC<TaxReportCardProps> = ({ report, onExport }) => {
  return (
    <Card className="border-border/60 shadow-sm">
      <CardHeader className="flex flex-row items-center justify-between pb-3">
        <div>
          <CardTitle className="text-base font-semibold">Monthly VAT / Tax Collected</CardTitle>
          <CardDescription>Fiscal breakdown of transactions and collected sales tax</CardDescription>
        </div>
        <Button variant="outline" size="sm" onClick={onExport} className="gap-1.5 text-xs font-semibold">
          <Download className="h-3.5 w-3.5" />
          Export CSV
        </Button>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Billing Month</TableHead>
              <TableHead>Settled Invoices</TableHead>
              <TableHead>Gross Revenue</TableHead>
              <TableHead className="text-right">Tax (PPN 11%)</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {!report ? (
              <TableRow>
                <TableCell colSpan={4} className="text-center py-6 text-muted-foreground text-xs">
                  Loading tax reports...
                </TableCell>
              </TableRow>
            ) : (
              report.monthly_breakdown.map((row) => (
                <TableRow key={row.month}>
                  <TableCell className="font-semibold text-xs">{row.month}</TableCell>
                  <TableCell className="text-xs">{row.invoices_count} Invoices</TableCell>
                  <TableCell className="font-bold text-xs">{formatMoney(row.revenue)}</TableCell>
                  <TableCell className="text-right font-mono text-xs font-bold text-emerald-600 dark:text-emerald-400">
                    {formatMoney(row.tax)}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};
