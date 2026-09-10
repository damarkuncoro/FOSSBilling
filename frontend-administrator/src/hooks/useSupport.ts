import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminSupportService } from '@/services/admin_support.service';
import type { SupportTicket } from '@/types/api';

export function useSupport() {
  const queryClient = useQueryClient();
  const [selectedTicketId, setSelectedTicketId] = useState<number | null>(null);
  const [replyText, setReplyText] = useState('');

  const { data: tickets = [], isLoading: ticketsLoading, refetch: fetchTickets } = useQuery({
    queryKey: ['admin', 'support', 'tickets'],
    queryFn: () => adminSupportService.listTickets(),
  });

  const { data: selectedTicket = null, isLoading: ticketLoading } = useQuery({
    queryKey: ['admin', 'support', 'tickets', selectedTicketId],
    queryFn: () => (selectedTicketId ? adminSupportService.getTicketDetail(selectedTicketId) : null),
    enabled: !!selectedTicketId,
  });

  const replyMutation = useMutation({
    mutationFn: ({ id, message }: { id: number; message: string }) => adminSupportService.replyTicket(id, message),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'support', 'tickets', selectedTicketId] });
      setReplyText('');
    },
  });

  const closeMutation = useMutation({
    mutationFn: (id: number) => adminSupportService.closeTicket(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'support', 'tickets'] });
      setSelectedTicketId(null);
    },
  });

  return {
    tickets,
    loading: ticketsLoading || ticketLoading,
    selectedTicket,
    setSelectedTicket: (t: SupportTicket | null) => setSelectedTicketId(t ? t.id : null),
    replyText,
    setReplyText,
    replyLoading: replyMutation.isPending,
    fetchTickets,
    handleReply: (e: React.FormEvent) => {
      e.preventDefault();
      if (selectedTicketId && replyText.trim()) replyMutation.mutate({ id: selectedTicketId, message: replyText });
    },
    handleClose: (id: number) => closeMutation.mutate(id),
  };
}
