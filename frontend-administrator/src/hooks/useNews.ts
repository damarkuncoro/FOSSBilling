import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminSystemService } from '@/services/admin_system.service';

export function useNews() {
  const queryClient = useQueryClient();
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState({ title: '', content: '' });

  const { data: articles = [], isLoading: loading, refetch: fetchNews } = useQuery({
    queryKey: ['admin', 'news'],
    queryFn: () => adminSystemService.listNewsArticles(),
  });

  const createMutation = useMutation({
    mutationFn: () => adminSystemService.createNewsArticle(form.title, form.content),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'news'] });
      setOpenModal(false);
      setForm({ title: '', content: '' });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminSystemService.deleteNewsArticle(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'news'] }),
  });

  return {
    articles,
    loading,
    openModal,
    setOpenModal,
    form,
    setForm,
    saving: createMutation.isPending || deleteMutation.isPending,
    fetchNews,
    handleCreate: (e: React.FormEvent) => {
      e.preventDefault();
      createMutation.mutate();
    },
    handleDelete: (id: number) => {
      if (confirm('Delete announcement?')) deleteMutation.mutate(id);
    },
  };
}
