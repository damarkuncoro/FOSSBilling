import { request } from '../lib/api/client';

export class NewsRepository {
  list = () => request<any[]>('/guest/news');
  get = (slug: string) => request<any>(`/guest/news/${slug}`);
}

export const newsRepository = new NewsRepository();
