export interface KbArticle {
  id: number;
  category: string;
  title: string;
  slug: string;
  summary: string;
  content: string;
  views: number;
  helpful_count: number;
  updated_at: string;
}

export interface NewsArticle {
  id: number;
  title: string;
  slug: string;
  category: 'announcement' | 'maintenance' | 'feature' | 'security';
  content: string;
  published_at: string;
  author: string;
}

export interface DomainRecord {
  id: number;
  domain_name: string;
  tld: string;
  status: 'active' | 'pending' | 'expired' | 'suspended';
  nameservers: string[];
  epp_code?: string;
  auto_renew: boolean;
  expires_at: string;
}

export interface DownloadItem {
  id: number;
  title: string;
  category: string;
  version: string;
  file_size: string;
  description: string;
  requires_active_service: boolean;
  download_url: string;
  updated_at: string;
}

export interface ClientLicense {
  id: number;
  product_title: string;
  license_key: string;
  status: 'active' | 'suspended' | 'expired';
  licensed_domain?: string;
  licensed_ip?: string;
  max_instances: number;
  expires_at?: string;
}
