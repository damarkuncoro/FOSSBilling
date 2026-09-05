import React from 'react';
import { LifeBuoy, Send } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

interface TicketReplyDialogProps {
  selectedTicket: any | null;
  onClose: () => void;
  replyText: string;
  setReplyText: (val: string) => void;
  onReply: (e: React.FormEvent) => void;
  replyLoading: boolean;
}

export const TicketReplyDialog: React.FC<TicketReplyDialogProps> = ({
  selectedTicket,
  onClose,
  replyText,
  setReplyText,
  onReply,
  replyLoading,
}) => {
  if (!selectedTicket) return null;

  return (
    <Dialog open={!!selectedTicket} onOpenChange={onClose}>
      <DialogContent className="max-w-xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <LifeBuoy className="h-5 w-5 text-primary" />
            <span>Ticket #{selectedTicket.id}: {selectedTicket.subject}</span>
          </DialogTitle>
          <DialogDescription>
            Client #{selectedTicket.client_id} • Priority: {selectedTicket.priority} • Status: {selectedTicket.status}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 my-2">
          <div className="p-4 rounded-lg bg-muted/40 border text-sm space-y-2">
            <p className="font-semibold text-xs text-muted-foreground uppercase">Initial Message:</p>
            <p className="whitespace-pre-wrap">{selectedTicket.content || selectedTicket.subject}</p>
          </div>

          <form onSubmit={onReply} className="space-y-3">
            <label className="text-xs font-semibold">Staff Response:</label>
            <Textarea
              required
              placeholder="Type your official response to the customer..."
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
              rows={4}
            />

            <div className="flex justify-end gap-2">
              <Button type="button" variant="outline" onClick={onClose}>
                Cancel
              </Button>
              <Button type="submit" disabled={replyLoading} className="gap-2">
                <Send className="h-4 w-4" />
                {replyLoading ? 'Sending...' : 'Send Reply'}
              </Button>
            </div>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  );
};
