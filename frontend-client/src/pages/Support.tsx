import React, { useEffect, useState } from 'react';
import {
  LifeBuoy,
  Plus,
  Send,
  MessageSquare,
  CheckCircle,
  RefreshCw,
  XCircle,
} from 'lucide-react';
import { api } from '@/lib/api';
import { formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';

export const Support: React.FC = () => {
  const [tickets, setTickets] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [openNewModal, setOpenNewModal] = useState(false);
  const [newTicketForm, setNewTicketForm] = useState({ subject: '', message: '', priority: 'medium' });
  const [selectedTicket, setSelectedTicket] = useState<any | null>(null);
  const [replyContent, setReplyContent] = useState('');
  const [replying, setReplying] = useState(false);

  const fetchTickets = async () => {
    setLoading(true);
    try {
      const data = await api.getTickets();
      setTickets(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTickets();
  }, []);

  const handleCreateTicket = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.openTicket(newTicketForm);
      setOpenNewModal(false);
      setNewTicketForm({ subject: '', message: '', priority: 'medium' });
      await fetchTickets();
    } catch (err: any) {
      alert(`Failed to open ticket: ${err.message}`);
    }
  };

  const handleReply = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedTicket || !replyContent.trim()) return;

    setReplying(true);
    try {
      await api.replyTicket(selectedTicket.id, replyContent);
      setReplyContent('');
      setSelectedTicket(null);
      await fetchTickets();
    } catch (err: any) {
      alert(`Failed to reply: ${err.message}`);
    } finally {
      setReplying(false);
    }
  };

  const handleClose = async (id: number) => {
    if (!confirm('Are you sure you want to close this ticket?')) return;
    try {
      await api.closeTicket(id);
      setSelectedTicket(null);
      await fetchTickets();
    } catch (err: any) {
      alert(`Failed to close ticket: ${err.message}`);
    }
  };

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

          <Dialog open={openNewModal} onOpenChange={setOpenNewModal}>
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
              <form onSubmit={handleCreateTicket} className="space-y-4 pt-2">
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">Inquiry Subject</label>
                  <Input
                    required
                    placeholder="e.g. Assistance with DNS records & SSL"
                    value={newTicketForm.subject}
                    onChange={(e) => setNewTicketForm({ ...newTicketForm, subject: e.target.value })}
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">Priority Level</label>
                  <select
                    value={newTicketForm.priority}
                    onChange={(e) => setNewTicketForm({ ...newTicketForm, priority: e.target.value })}
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
                    value={newTicketForm.message}
                    onChange={(e) => setNewTicketForm({ ...newTicketForm, message: e.target.value })}
                  />
                </div>
                <div className="flex justify-end gap-2 pt-2">
                  <Button type="button" variant="outline" onClick={() => setOpenNewModal(false)}>
                    Cancel
                  </Button>
                  <Button type="submit">Submit Ticket</Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Support Tickets ({tickets.length})</CardTitle>
          <CardDescription>Track status updates, staff replies, and resolution threads</CardDescription>
        </CardHeader>
        <CardContent>
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
                    <TableCell className="text-xs text-muted-foreground">{formatDate(t.created_at)}</TableCell>
                    <TableCell className="text-right">
                      <Button
                        size="sm"
                        variant="outline"
                        className="h-7 text-xs"
                        onClick={() => setSelectedTicket(t)}
                      >
                        View Thread
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Ticket Conversation Thread Modal */}
      {selectedTicket && (
        <Dialog open={!!selectedTicket} onOpenChange={() => setSelectedTicket(null)}>
          <DialogContent className="max-w-xl">
            <DialogHeader>
              <DialogTitle>Ticket #{selectedTicket.id}: {selectedTicket.subject}</DialogTitle>
              <DialogDescription>
                Priority: {selectedTicket.priority} • Status: {selectedTicket.status}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 my-2">
              <div className="p-4 rounded-xl bg-muted/40 border text-sm space-y-1">
                <p className="text-xs font-bold text-muted-foreground uppercase">Initial Request:</p>
                <p className="whitespace-pre-wrap">{selectedTicket.content || selectedTicket.subject}</p>
              </div>

              {selectedTicket.status !== 'closed' ? (
                <form onSubmit={handleReply} className="space-y-3">
                  <label className="text-xs font-semibold">Post a Reply:</label>
                  <Textarea
                    required
                    placeholder="Type additional information or reply to staff..."
                    value={replyContent}
                    onChange={(e) => setReplyContent(e.target.value)}
                    rows={4}
                  />
                  <div className="flex justify-between items-center">
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="text-xs text-destructive gap-1"
                      onClick={() => handleClose(selectedTicket.id)}
                    >
                      <XCircle className="h-3.5 w-3.5" />
                      Close Ticket
                    </Button>
                    <div className="flex gap-2">
                      <Button type="button" variant="outline" onClick={() => setSelectedTicket(null)}>
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
      )}
    </div>
  );
};
