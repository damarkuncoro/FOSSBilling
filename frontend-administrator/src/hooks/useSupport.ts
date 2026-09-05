import { useState, useEffect } from 'react';
import { api } from '@/lib/api';

export function useSupport() {
  const [tickets, setTickets] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedTicket, setSelectedTicket] = useState<any | null>(null);
  const [replyText, setReplyText] = useState('');
  const [replyLoading, setReplyLoading] = useState(false);

  const fetchTickets = async () => {
    setLoading(true);
    try {
      const data = await api.getSupportTickets();
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
      await api.replySupportTicket(selectedTicket.id, replyText);
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
