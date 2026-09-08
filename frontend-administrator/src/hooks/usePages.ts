import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { CustomPageItem, KnowledgebaseArticle } from '@/types/modules';

export function usePages() {
  const [pages, setPages] = useState<CustomPageItem[]>([]);
  const [articles, setArticles] = useState<KnowledgebaseArticle[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedPage, setSelectedPage] = useState<CustomPageItem | null>(null);
  const [openPageModal, setOpenPageModal] = useState(false);
  const [openKBModal, setOpenKBModal] = useState(false);
  const [pageForm, setPageForm] = useState<Partial<CustomPageItem>>({
    title: '',
    slug: '',
    content: '',
    published: true,
  });
  const [kbForm, setKBForm] = useState<Partial<KnowledgebaseArticle>>({
    title: '',
    slug: '',
    content: '',
    category: '',
    published: true,
  });

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const [pagesData, kbData] = await Promise.all([
        api.getPages().catch(() => []),
        api.getKnowledgebase().catch(() => []),
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

  const handleSaveKB = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.saveKnowledgebase(kbForm);
      setOpenKBModal(false);
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

  const handleDeleteKB = async (id: number) => {
    if (!confirm('Are you sure you want to delete this KB article?')) return;
    try {
      await api.deleteKnowledgebase(id);
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
    openKBModal,
    setOpenKBModal,
    pageForm,
    setPageForm,
    kbForm,
    setKBForm,
    fetchData,
    handleSavePage,
    handleSaveKB,
    handleDeletePage,
    handleDeleteKB,
  };
}
