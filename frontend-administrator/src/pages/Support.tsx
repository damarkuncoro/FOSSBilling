import React from 'react';
import { RefreshCw, MessageSquare } from 'lucide-react';
import { useSupport } from '@/hooks/useSupport';
import { TicketReplyDialog } from '@/components/support/TicketReplyDialog';
import { formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const Support: React.FC = () => {
  const {
    tickets,
    loading,
    selectedTicket,
    setSelectedTicket,
    replyText,
    setReplyText,
    replyLoading,
    fetchTickets,
    handleReply,
  } = useSupport();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Customer Support Helpdesk</h1>
          <p className="text-sm text-muted-foreground">
            Manage customer inquiries, reply to support tickets, and resolve issues.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchTickets} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Active Support Inquiries ({tickets.length})</CardTitle>
          <CardDescription>Multi-department ticketing with conversation threading and status workflows</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-16">ID</TableHead>
                <TableHead>Subject</TableHead>
                <TableHead>Client ID</TableHead>
                <TableHead>Priority</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Opened At</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                    Loading tickets...
                  </TableCell>
                </TableRow>
              ) : tickets.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                    No support tickets found.
                  </TableCell>
                </TableRow>
              ) : (
                tickets.map((ticket) => (
                  <TableRow key={ticket.id}>
                    <TableCell className="font-mono text-xs font-semibold">#{ticket.id}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <MessageSquare className="h-4 w-4 text-primary shrink-0" />
                        <span className="font-medium text-sm">{ticket.subject}</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-xs font-mono">Client #{ticket.client_id}</TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          ticket.priority === 'urgent'
                            ? 'destructive'
                            : ticket.priority === 'high'
                            ? 'warning'
                            : 'secondary'
                        }
                        className="text-[10px] uppercase font-bold"
                      >
                        {ticket.priority || 'medium'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          ticket.status === 'open'
                            ? 'info'
                            : ticket.status === 'closed'
                            ? 'secondary'
                            : 'warning'
                        }
                      >
                        {ticket.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {formatDate(ticket.created_at)}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        size="sm"
                        variant="outline"
                        className="h-7 text-xs gap-1"
                        onClick={() => setSelectedTicket(ticket)}
                      >
                        Inspect & Reply
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <TicketReplyDialog
        selectedTicket={selectedTicket}
        onClose={() => setSelectedTicket(null)}
        replyText={replyText}
        setReplyText={setReplyText}
        onReply={handleReply}
        replyLoading={replyLoading}
      />
    </div>
  );
};

export default Support;
