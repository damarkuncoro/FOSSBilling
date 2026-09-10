import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';

export function usePages() {
  const queryClient = useQueryClient();
  const [openPageModal, setOpenPageModal] = useState(false);
  const [openKBModal, setOpenKBModal] = useState(false);
  const [pageForm, setPageForm] = useState<any>({ title: '', slug: '', content: '', published: true });
  const [kbForm, setKBForm] = useState<any>({ title: '', slug: '', content: '', category: '', published: true });

  const { data: pages = [], isLoading: pLoading, refetch: refetchPages } = useQuery({ queryKey: ['admin', 'pages'], queryFn: () => api.getPages().catch(() => []) });
  const { data: articles = [], isLoading: kLoading, refetch: refetchKB } = useQuery({ queryKey: ['admin', 'kb'], queryFn: () => api.getKnowledgebase().catch(() => []) });

  const pageMutation = useMutation({
    mutationFn: (d: any) => api.savePage(d),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin', 'pages'] }); setOpenPageModal(false); },
  });

  const kbMutation = useMutation({
    mutationFn: (d: any) => api.saveKnowledgebase(d),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin', 'kb'] }); setOpenKBModal(false); },
  });

  const deletePageMutation = useMutation({
    mutationFn: (id: number) => api.deletePage(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'pages'] }),
  });

  const deleteKBMutation = useMutation({
    mutationFn: (id: number) => api.deleteKnowledgebase(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'kb'] }),
  });

  const fetchData = async () => {
    await Promise.all([refetchPages(), refetchKB()]);
  };

  return {
    pages, articles, loading: pLoading || kLoading,
    fetchData,
    openPageModal, setOpenPageModal, openKBModal, setOpenKBModal,
    pageForm, setPageForm, kbForm, setKBForm,
    handleSavePage: (e: any) => { e.preventDefault(); pageMutation.mutate(pageForm); },
    handleSaveKB: (e: any) => { e.preventDefault(); kbMutation.mutate(kbForm); },
    handleDeletePage: (id: number) => { if (confirm('Delete?')) deletePageMutation.mutate(id); },
    handleDeleteKB: (id: number) => { if (confirm('Delete?')) deleteKBMutation.mutate(id); },
  };
}
