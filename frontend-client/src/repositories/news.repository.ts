import { request } from '../lib/api/client';

export interface INewsRepository {
  listNews(): Promise<any[]>;
  getNewsBySlug(slug: string): Promise<any>;
}

export class NewsRepository implements INewsRepository {
  async listNews(): Promise<any[]> {
    return request<any[]>('/guest/news');
  }

  async getNewsBySlug(slug: string): Promise<any> {
    return request<any>(`/guest/news/${slug}`);
  }
}

export const newsRepository = new NewsRepository();
