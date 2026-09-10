import { newsRepository as repo } from '../repositories/news.repository';

export class NewsService {
  listPublishedNews = () => repo.list();
  getNewsDetail = (s: string) => repo.get(s);
}

export const newsService = new NewsService();
