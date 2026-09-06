import { useState, useEffect } from 'react';
import type { NewsArticle } from '../types/clientModules';
import { newsService } from '../services/news.service';

export function useNewsArticles(initial: NewsArticle[] = []) {
  const [articles, setArticles] = useState<NewsArticle[]>(initial);
  const [loading, setLoading] = useState(false);
  const [selectedArticle, setSelectedArticle] = useState<NewsArticle | null>(null);

  const fetchNews = async () => {
    try {
      setLoading(true);
      const res = await newsService.listPublishedNews();
      if (res && Array.isArray(res) && res.length > 0) {
        const mapped: NewsArticle[] = res.map((item: any) => ({
          id: item.id,
          title: item.title,
          slug: item.slug || `news-${item.id}`,
          category: (item.category || 'announcement') as any,
          content: item.content,
          published_at: item.created_at || item.published_at || new Date().toISOString(),
          author: item.author || 'FOSSBilling Staff',
        }));
        setArticles(mapped);
      }
    } catch {
      // Retain initial state if offline or mocked
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNews();
  }, []);

  return {
    articles,
    loading,
    selectedArticle,
    setSelectedArticle,
    refreshNews: fetchNews,
  };
}

