import { useState } from 'react';
import type { NewsArticle } from '../types/clientModules';

const initialNews: NewsArticle[] = [
  {
    id: 1,
    title: 'Scheduled Maintenance: Singapore Datacenter Core Switch Upgrade',
    slug: 'singapore-core-switch-upgrade',
    category: 'maintenance',
    content: 'We will be performing a scheduled hardware firmware upgrade on core Singapore switches on Sunday at 02:00 UTC. Expected downtime is under 5 minutes.',
    published_at: '2026-09-04T08:00:00Z',
    author: 'Network Ops Team',
  },
  {
    id: 2,
    title: 'New High-Performance AMD EPYC 9004 Gen 5 NVMe Plans Launched',
    slug: 'amd-epyc-gen5-launch',
    category: 'announcement',
    content: 'We are thrilled to announce our new Gen 5 NVMe Cloud Compute instances powered by AMD EPYC processors with 10Gbps unmetered uplink.',
    published_at: '2026-08-28T10:00:00Z',
    author: 'Product Management',
  },
];

export function useNewsArticles() {
  const [articles] = useState<NewsArticle[]>(initialNews);
  const [selectedArticle, setSelectedArticle] = useState<NewsArticle | null>(null);

  return {
    articles,
    selectedArticle,
    setSelectedArticle,
  };
}
