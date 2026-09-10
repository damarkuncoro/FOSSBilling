import { useState, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';

export function useKnowledgebase() {
  const [search, setSearch] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [activeSlug, setActiveSlug] = useState<string | null>(null);

  const { data: categories = [] } = useQuery({ queryKey: ['public', 'kb', 'categories'], queryFn: () => api.getKBCategories().catch(() => []) });
  const { data: articles = [], isLoading } = useQuery({ queryKey: ['public', 'kb', 'articles'], queryFn: () => api.getKBArticles().catch(() => []) });

  const { data: activeArticle = null } = useQuery({
    queryKey: ['public', 'kb', 'article', activeSlug],
    queryFn: () => activeSlug ? api.getKBArticle(activeSlug) : null,
    enabled: !!activeSlug,
  });

  const filtered = useMemo(() => articles.filter((a: any) => {
    const mC = selectedCategory === 'all' || a.category === selectedCategory;
    const mS = a.title.toLowerCase().includes(search.toLowerCase());
    return mC && mS;
  }), [articles, search, selectedCategory]);

  return {
    articles: filtered, categories: ['all', ...categories.map((c: any) => c.title)],
    loading: isLoading, search, setSearch, selectedCategory, setSelectedCategory,
    activeArticle, setActiveArticle: (a: any) => setActiveSlug(a ? a.slug : null),
  };
}
