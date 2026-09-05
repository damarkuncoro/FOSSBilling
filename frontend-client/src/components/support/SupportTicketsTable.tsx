import React from 'react';
import { MessageSquare } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatDate } from '@/lib/utils';
import { SupportTicket } from '@/types/api';

interface SupportTicketsTableProps {
  tickets: SupportTicket[];
  loading: boolean;
  onSelectTicket: (ticket: SupportTicket) => void;
}

export const SupportTicketsTable: React.FC<SupportTicketsTableProps> = ({
  tickets,
  loading,
  onSelectTicket,
}) => {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-16">ID</TableHead>
          <TableHead>Subject</TableHead>
          <TableHead>Priority</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Submitted Date</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {loading ? (
          <TableRow>
            <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
              Loading tickets...
            </TableCell>
          </TableRow>
        ) : tickets.length === 0 ? (
          <TableRow>
            <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
              You have no active support tickets.
            </TableCell>
          </TableRow>
        ) : (
          tickets.map((t) => (
            <TableRow key={t.id}>
              <TableCell className="font-mono text-xs font-semibold">#{t.id}</TableCell>
              <TableCell>
                <div className="flex items-center gap-2 font-medium">
                  <MessageSquare className="h-4 w-4 text-primary shrink-0" />
                  <span>{t.subject}</span>
                </div>
              </TableCell>
              <TableCell>
                <Badge
                  variant={
                    t.priority === 'urgent'
                      ? 'destructive'
                      : t.priority === 'high'
                      ? 'warning'
                      : 'secondary'
                  }
                  className="text-[10px] uppercase font-bold"
                >
                  {t.priority || 'medium'}
                </Badge>
              </TableCell>
              <TableCell>
                <Badge
                  variant={
                    t.status === 'open'
                      ? 'info'
                      : t.status === 'closed'
                      ? 'secondary'
                      : 'warning'
                  }
                >
                  {t.status.toUpperCase()}
                </Badge>
              </TableCell>
              <TableCell className="text-xs text-muted-foreground">
                {formatDate(t.created_at)}
              </TableCell>
              <TableCell className="text-right">
                <Button
                  size="sm"
                  variant="outline"
                  className="h-7 text-xs"
                  onClick={() => onSelectTicket(t)}
                >
                  View Thread
                </Button>
              </TableCell>
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  );
};
