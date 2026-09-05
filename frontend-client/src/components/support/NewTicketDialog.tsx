import React from 'react';
import { Plus } from 'lucide-react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

interface NewTicketDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: { subject: string; message: string; priority: string };
  onChange: (form: { subject: string; message: string; priority: string }) => void;
  onSubmit: (e: React.FormEvent) => void;
}

export const NewTicketDialog: React.FC<NewTicketDialogProps> = ({
  open,
  onOpenChange,
  form,
  onChange,
  onSubmit,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button size="sm" className="gap-1.5 shadow-sm">
          <Plus className="h-4 w-4" />
          Open New Ticket
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Open Support Inquiry</DialogTitle>
          <DialogDescription>Submit your request to our technical team</DialogDescription>
        </DialogHeader>
        <form onSubmit={onSubmit} className="space-y-4 pt-2">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Inquiry Subject</label>
            <Input
              required
              placeholder="e.g. Assistance with DNS records & SSL"
              value={form.subject}
              onChange={(e) => onChange({ ...form, subject: e.target.value })}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Priority Level</label>
            <select
              value={form.priority}
              onChange={(e) => onChange({ ...form, priority: e.target.value })}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus:outline-none focus:ring-1 focus:ring-ring"
            >
              <option value="low">Low - General Question</option>
              <option value="medium">Medium - Standard Request</option>
              <option value="high">High - Production Degraded</option>
              <option value="urgent">Urgent - Outage Critical</option>
            </select>
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Detailed Message</label>
            <Textarea
              required
              rows={5}
              placeholder="Describe your issue or question..."
              value={form.message}
              onChange={(e) => onChange({ ...form, message: e.target.value })}
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Submit Ticket</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
