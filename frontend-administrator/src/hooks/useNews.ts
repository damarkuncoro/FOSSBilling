import { useState, useEffect } from 'react';
import { adminSystemService } from '@/services/admin_system.service';

export function useNews() {
  const [articles, setArticles] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState({ title: '', content: '' });
  const [saving, setSaving] = useState(false);

  const fetchNews = async () => {
    setLoading(true);
    try {
      const data = await adminSystemService.listNewsArticles();
      setArticles(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNews();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await adminSystemService.createNewsArticle(form.title, form.content);
      setOpenModal(false);
      setForm({ title: '', content: '' });
      await fetchNews();
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this announcement?')) return;
    try {
      await adminSystemService.deleteNewsArticle(id);
      await fetchNews();
    } catch (err) {
      console.error(err);
    }
  };

  return {
    articles,
    loading,
    openModal,
    setOpenModal,
    form,
    setForm,
    saving,
    fetchNews,
    handleCreate,
    handleDelete,
  };
}
