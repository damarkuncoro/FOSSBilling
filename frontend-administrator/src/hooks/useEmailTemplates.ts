import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';

export function useEmailTemplates() {
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<'templates' | 'smtp'>('templates');
  const [editModalOpen, setEditModalOpen] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState<any>(null);
  const [tplForm, setTplForm] = useState<any>({});
  const [testEmail, setTestEmail] = useState('');

  const { data: templates = [], isLoading: tLoading } = useQuery({ queryKey: ['admin', 'email', 'templates'], queryFn: () => api.getEmailTemplates().catch(() => []) });
  const { data: mailConfig = {}, isLoading: cLoading } = useQuery({ queryKey: ['admin', 'email', 'config'], queryFn: () => api.getMailConfig().catch(() => ({})) });

  const tplMutation = useMutation({
    mutationFn: (d: any) => api.updateEmailTemplate(d),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin', 'email', 'templates'] }); setEditModalOpen(false); },
  });

  const configMutation = useMutation({
    mutationFn: (d: any) => api.updateMailConfig(d),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'email', 'config'] }),
  });

  const testMutation = useMutation({
    mutationFn: (email: string) => api.sendTestEmail(email),
  });

  return {
    activeTab, setActiveTab, templates, mailConfig, loading: tLoading || cLoading,
    editModalOpen, setEditModalOpen, tplForm, setTplForm, savingTpl: tplMutation.isPending,
    testEmail, setTestEmail, sendingTest: testMutation.isPending, testStatus: testMutation.data,
    savingConfig: configMutation.isPending,
    openEdit: (t: any) => { setEditingTemplate(t); setTplForm(t); setEditModalOpen(true); },
    handleSaveTemplate: (e: any) => { e.preventDefault(); tplMutation.mutate(tplForm); },
    handleSaveMailConfig: (e: any) => { e.preventDefault(); configMutation.mutate(mailConfig); },
    handleSendTestEmail: () => testMutation.mutate(testEmail),
  };
}
