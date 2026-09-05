export interface UrlRedirect {
  id: number;
  source_path: string;
  target_url: string;
  status_code: 301 | 302;
  is_active: boolean;
  hits_count: number;
  created_at: string;
  updated_at: string;
}
