import { useState, useEffect } from 'react';
import { adminSupportService } from '@/services/admin_support.service';
import type { SupportTicket } from '@/types/api';

export function useSupport() {
  const [tickets, setTickets] = useState<SupportTicket[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedTicket, setSelectedTicket] = useState<SupportTicket | null>(null);
  const [replyText, setReplyText] = useState('');
  const [replyLoading, setReplyLoading] = useState(false);

  const fetchTickets = async () => {
    setLoading(true);
    try {
      const data = await adminSupportService.listTickets();
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

  const handleReply = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedTicket || !replyText.trim()) return;

    setReplyLoading(true);
    try {
      await adminSupportService.replyTicket(selectedTicket.id, replyText);
      setReplyText('');
      setSelectedTicket(null);
      await fetchTickets();
    } catch (err) {
      console.error(err);
    } finally {
      setReplyLoading(false);
    }
  };

  return {
    tickets,
    loading,
    selectedTicket,
    setSelectedTicket,
    replyText,
    setReplyText,
    replyLoading,
    fetchTickets,
    handleReply,
  };
}
