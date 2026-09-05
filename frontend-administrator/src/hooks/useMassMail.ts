import { useState, useEffect } from 'react';
import { api } from '@/lib/api';

export function useMassMail() {
  const [campaigns, setCampaigns] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState({ subject: '', content: '' });
  const [saving, setSaving] = useState(false);
  const [sendingId, setSendingId] = useState<number | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  const fetchCampaigns = async () => {
    setLoading(true);
    try {
      const data = await api.getMassMailCampaigns();
      setCampaigns(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCampaigns();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await api.createMassMailCampaign(form);
      setOpenModal(false);
      setForm({ subject: '', content: '' });
      await fetchCampaigns();
    } finally {
      setSaving(false);
    }
  };

  const handleSend = async (id: number) => {
    setSendingId(id);
    setMessage(null);
    try {
      await api.sendMassMailCampaign(id);
      setMessage(`Campaign #${id} has been broadcasted to all active clients!`);
      await fetchCampaigns();
    } catch (err: any) {
      setMessage(`Error sending campaign: ${err.message}`);
    } finally {
      setSendingId(null);
    }
  };

  return {
    campaigns,
    loading,
    openModal,
    setOpenModal,
    form,
    setForm,
    saving,
    sendingId,
    message,
    setMessage,
    fetchCampaigns,
    handleCreate,
    handleSend,
  };
}
