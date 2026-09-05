import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { CustomPageItem, KnowledgebaseArticle } from '@/types/modules';

export function usePages() {
  const [pages, setPages] = useState<CustomPageItem[]>([]);
  const [articles, setArticles] = useState<KnowledgebaseArticle[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedPage, setSelectedPage] = useState<CustomPageItem | null>(null);
  const [openPageModal, setOpenPageModal] = useState(false);
  const [pageForm, setPageForm] = useState<Partial<CustomPageItem>>({
    title: '',
    slug: '',
    content: '',
    published: true,
  });

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const [pagesData, kbData] = await Promise.all([
        api.getPages().catch(() => [
          { id: 1, title: 'Terms of Service', slug: 'terms-of-service', content: 'Standard hosting terms...', published: true, updated_at: '2026-09-01T00:00:00Z' },
          { id: 2, title: 'Privacy Policy', slug: 'privacy-policy', content: 'Data protection policies...', published: true, updated_at: '2026-09-01T00:00:00Z' },
          { id: 3, title: 'Service Level Agreement (SLA)', slug: 'sla', content: '99.9% uptime guarantee rules...', published: true, updated_at: '2026-09-01T00:00:00Z' },
        ]),
        api.getKnowledgebase().catch(() => [
          { id: 1, category: 'Hosting & cPanel', title: 'How to point DNS A Records to your VPS', slug: 'how-to-point-dns', content: 'Guide for nameserver mapping...', views: 420, published: true, updated_at: '2026-09-02T00:00:00Z' },
          { id: 2, category: 'Billing & Payments', title: 'Automatic invoice settlement with QRIS', slug: 'qris-payment-guide', content: 'Instant QRIS payment walkthrough...', views: 180, published: true, updated_at: '2026-09-03T00:00:00Z' },
        ]),
      ]);
      setPages(pagesData || []);
      setArticles(kbData || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleSavePage = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.savePage(pageForm);
      setOpenPageModal(false);
      await fetchData();
    } catch (err: any) {
      alert(`Save failed: ${err.message}`);
    }
  };

  const handleDeletePage = async (id: number) => {
    if (!confirm('Are you sure you want to delete this custom page?')) return;
    try {
      await api.deletePage(id);
      await fetchData();
    } catch (err: any) {
      alert(`Delete failed: ${err.message}`);
    }
  };

  return {
    pages,
    articles,
    loading,
    selectedPage,
    setSelectedPage,
    openPageModal,
    setOpenPageModal,
    pageForm,
    setPageForm,
    fetchData,
    handleSavePage,
    handleDeletePage,
  };
}
