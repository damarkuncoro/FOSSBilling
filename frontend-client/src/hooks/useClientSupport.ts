import { useState, useEffect, useCallback } from 'react';
import { supportService } from '@/services/support.service';
import { SupportTicket } from '@/types/api';

export function useClientSupport() {
  const [tickets, setTickets] = useState<SupportTicket[]>([]);
  const [loading, setLoading] = useState(true);
  const [openNewModal, setOpenNewModal] = useState(false);
  const [newTicketForm, setNewTicketForm] = useState({
    subject: '',
    message: '',
    priority: 'medium',
  });
  const [selectedTicket, setSelectedTicket] = useState<SupportTicket | null>(null);
  const [replyContent, setReplyContent] = useState('');
  const [replying, setReplying] = useState(false);

  const fetchTickets = useCallback(async () => {
    setLoading(true);
    try {
      const data = await supportService.listTickets();
      setTickets(data || []);
    } catch (err) {
      console.error('Failed to fetch tickets:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTickets();
  }, [fetchTickets]);

  const handleCreateTicket = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await supportService.openTicket(newTicketForm.subject, newTicketForm.message, newTicketForm.priority);
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
      await supportService.replyTicket(selectedTicket.id, replyContent);
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
      await supportService.closeTicket(id);
      setSelectedTicket(null);
      await fetchTickets();
    } catch (err: any) {
      alert(`Failed to close ticket: ${err.message}`);
    }
  };

  return {
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
  };
}
