import React from 'react';
import { RefreshCw } from 'lucide-react';
import { useClientSupport } from '@/hooks/useClientSupport';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { SupportTicketsTable } from '@/components/support/SupportTicketsTable';
import { NewTicketDialog } from '@/components/support/NewTicketDialog';
import { TicketThreadDialog } from '@/components/support/TicketThreadDialog';

export const Support: React.FC = () => {
  const {
    tickets,
    loading,
    openNewModal,
    setOpenNewModal,
    newTicketForm,
    setNewTicketForm,
    selectedTicket,
    setSelectedTicket,
    replyContent,
    setReplyContent,
    replying,
    fetchTickets,
    handleCreateTicket,
    handleReply,
    handleClose,
  } = useClientSupport();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Customer Support & Helpdesk</h1>
          <p className="text-sm text-muted-foreground">
            Get technical assistance from our 24/7 cloud engineering and billing support team.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchTickets} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          <NewTicketDialog
            open={openNewModal}
            onOpenChange={setOpenNewModal}
            form={newTicketForm}
            onChange={setNewTicketForm}
            onSubmit={handleCreateTicket}
          />
        </div>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">
            Support Tickets ({tickets.length})
          </CardTitle>
          <CardDescription>
            Track status updates, staff replies, and resolution threads
          </CardDescription>
        </CardHeader>
        <CardContent>
          <SupportTicketsTable
            tickets={tickets}
            loading={loading}
            onSelectTicket={setSelectedTicket}
          />
        </CardContent>
      </Card>

      <TicketThreadDialog
        ticket={selectedTicket}
        onClose={() => setSelectedTicket(null)}
        replyContent={replyContent}
        onReplyChange={setReplyContent}
        onSendReply={handleReply}
        onCloseTicket={handleClose}
        replying={replying}
      />
    </div>
  );
};
