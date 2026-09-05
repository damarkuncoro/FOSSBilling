import React from 'react';
import { Send, XCircle } from 'lucide-react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { SupportTicket } from '@/types/api';

interface TicketThreadDialogProps {
  ticket: SupportTicket | null;
  onClose: () => void;
  replyContent: string;
  onReplyChange: (content: string) => void;
  onSendReply: (e: React.FormEvent) => void;
  onCloseTicket: (id: number) => void;
  replying: boolean;
}

export const TicketThreadDialog: React.FC<TicketThreadDialogProps> = ({
  ticket,
  onClose,
  replyContent,
  onReplyChange,
  onSendReply,
  onCloseTicket,
  replying,
}) => {
  if (!ticket) return null;

  return (
    <Dialog open={!!ticket} onOpenChange={onClose}>
      <DialogContent className="max-w-xl">
        <DialogHeader>
          <DialogTitle>
            Ticket #{ticket.id}: {ticket.subject}
          </DialogTitle>
          <DialogDescription>
            Priority: {ticket.priority} • Status: {ticket.status}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 my-2">
          <div className="p-4 rounded-xl bg-muted/40 border text-sm space-y-1">
            <p className="text-xs font-bold text-muted-foreground uppercase">Initial Request:</p>
            <p className="whitespace-pre-wrap">{ticket.content || ticket.subject}</p>
          </div>

          {ticket.status !== 'closed' ? (
            <form onSubmit={onSendReply} className="space-y-3">
              <label className="text-xs font-semibold">Post a Reply:</label>
              <Textarea
                required
                placeholder="Type additional information or reply to staff..."
                value={replyContent}
                onChange={(e) => onReplyChange(e.target.value)}
                rows={4}
              />
              <div className="flex justify-between items-center">
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="text-xs text-destructive gap-1"
                  onClick={() => onCloseTicket(ticket.id)}
                >
                  <XCircle className="h-3.5 w-3.5" />
                  Close Ticket
                </Button>
                <div className="flex gap-2">
                  <Button type="button" variant="outline" onClick={onClose}>
                    Cancel
                  </Button>
                  <Button type="submit" disabled={replying} className="gap-1.5">
                    <Send className="h-4 w-4" />
                    {replying ? 'Sending...' : 'Send Reply'}
                  </Button>
                </div>
              </div>
            </form>
          ) : (
            <div className="p-3 rounded-lg bg-muted text-center text-xs text-muted-foreground">
              This ticket has been marked as resolved and closed.
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
};
