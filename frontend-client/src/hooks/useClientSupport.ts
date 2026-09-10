import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { supportService } from '@/services/support.service';
import type { SupportTicket } from '@/types/api';

export function useClientSupport() {
  const queryClient = useQueryClient();
  const [openNewModal, setOpenNewModal] = useState(false);
  const [newTicketForm, setNewTicketForm] = useState({ subject: '', message: '', priority: 'medium' });
  const [selectedTicketId, setSelectedTicketId] = useState<number | null>(null);
  const [replyContent, setReplyContent] = useState('');

  const { data: tickets = [], isLoading: tl, refetch } = useQuery({ queryKey: ['client', 'support', 'tickets'], queryFn: () => supportService.listTickets() });
  const { data: selectedTicket = null, isLoading: tkl } = useQuery({ queryKey: ['client', 'support', 'tickets', selectedTicketId], queryFn: () => (selectedTicketId ? supportService.getTicketDetail(selectedTicketId) : null), enabled: !!selectedTicketId });

  const cM = useMutation({ mutationFn: () => supportService.openTicket(newTicketForm.subject, newTicketForm.message, newTicketForm.priority), onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['client', 'support', 'tickets'] }); setOpenNewModal(false); setNewTicketForm({ subject: '', message: '', priority: 'medium' }); } });
  const rM = useMutation({ mutationFn: ({ id, message }: { id: number; message: string }) => supportService.replyTicket(id, message), onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['client', 'support', 'tickets', selectedTicketId] }); setReplyContent(''); } });
  const clM = useMutation({ mutationFn: (id: number) => supportService.closeTicket(id), onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['client', 'support', 'tickets'] }); setSelectedTicketId(null); } });

  return {
    tickets, loading: tl || tkl, openNewModal, setOpenNewModal, newTicketForm, setNewTicketForm, selectedTicket, setSelectedTicket: (t: any) => setSelectedTicketId(t ? t.id : null), replyContent, setReplyContent,
    replying: rM.isPending, fetchTickets: () => refetch(),
    handleCreateTicket: (e: any) => { e.preventDefault(); cM.mutate(); },
    handleReply: (e: any) => { e.preventDefault(); if (selectedTicketId && replyContent.trim()) rM.mutate({ id: selectedTicketId, message: replyContent }); },
    handleClose: (id: number) => { if (confirm('Close?')) clM.mutate(id); },
  };
}
