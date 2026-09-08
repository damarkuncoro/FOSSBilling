import { useState, useEffect, useMemo } from 'react';
import { api } from '@/lib/api';
import type { KbArticle } from '../types/clientModules';

const DEFAULT_SAMPLE_ARTICLES: KbArticle[] = [
  {
    id: 1,
    title: 'How to Access Your cPanel Hosting Control Panel',
    slug: 'access-cpanel-hosting',
    category: 'Web Hosting',
    views: 1240,
    helpful_count: 42,
    updated_at: '2026-08-10',
    summary: 'Step by step guide to logging into cPanel using SSO or direct credentials.',
    content: 'Full step-by-step instructions on accessing cPanel directly or through one-click SSO.',
  },
  {
    id: 2,
    title: 'How to Install Free SSL Certificate (Let\'s Encrypt)',
    slug: 'install-ssl-certificate',
    category: 'Security & SSL',
    views: 890,
    helpful_count: 28,
    updated_at: '2026-08-12',
    summary: 'Enable one-click automated SSL issuance for your domain.',
    content: 'Navigate to Security settings in your cPanel dashboard and issue free SSL certificates.',
  },
];

const DEFAULT_SAMPLE_CATEGORIES = [
  { id: 1, title: 'Web Hosting' },
  { id: 2, title: 'Security & SSL' },
  { id: 3, title: 'Domains & DNS' },
];

export function useKnowledgebase(
  initialArticles: KbArticle[] = DEFAULT_SAMPLE_ARTICLES,
  initialCategories: any[] = DEFAULT_SAMPLE_CATEGORIES
) {
  const [articles, setArticles] = useState<KbArticle[]>(initialArticles);
  const [categories, setCategories] = useState<any[]>(initialCategories);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [activeArticle, setActiveArticle] = useState<KbArticle | null>(null);

  useEffect(() => {
    let isMounted = true;
    const load = async () => {
      try {
        const [cats, arts] = await Promise.all([
          api.getKBCategories(),
          api.getKBArticles(),
        ]);
        if (isMounted) {
          if (cats && cats.length > 0) setCategories(cats);
          if (arts && arts.length > 0) setArticles(arts);
        }
      } catch {
        // Retain fallback data gracefully
      }
    };
    load();
    return () => { isMounted = false; };
  }, []);

  const categoryOptions = useMemo(() => {
    return ['all', ...categories.map(c => c.title)];
  }, [categories]);

  const filteredArticles = useMemo(() => {
    return articles.filter((art) => {
      const matchesCategory = selectedCategory === 'all' || art.category === selectedCategory;
      const matchesSearch =
        art.title.toLowerCase().includes(search.toLowerCase()) ||
        (art.summary?.toLowerCase().includes(search.toLowerCase()));
      return matchesCategory && matchesSearch;
    });
  }, [articles, search, selectedCategory]);

  const handleSetActiveArticle = (art: KbArticle | null) => {
    setActiveArticle(art);
    if (art?.slug) {
      api.getKBArticle(art.slug)
        .then((detail) => {
          setActiveArticle((prev) => (prev ? { ...prev, ...detail } : prev));
        })
        .catch(() => {});
    }
  };

  return {
    articles: filteredArticles,
    categories: categoryOptions,
    loading,
    search,
    setSearch,
    selectedCategory,
    setSelectedCategory,
    activeArticle,
    setActiveArticle: handleSetActiveArticle,
  };
}
