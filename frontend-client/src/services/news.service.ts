import { NewsRepository, newsRepository, INewsRepository } from '../repositories/news.repository';

export class NewsService {
  constructor(private repo: INewsRepository = newsRepository) {}

  async listPublishedNews(): Promise<any[]> {
    return this.repo.listNews();
  }

  async getNewsArticle(slug: string): Promise<any> {
    if (!slug) {
      throw new Error('News article slug is required');
    }
    return this.repo.getNewsBySlug(slug);
  }
}

export const newsService = new NewsService();
