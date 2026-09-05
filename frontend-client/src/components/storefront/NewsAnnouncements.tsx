import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface NewsAnnouncementsProps {
  news: any[];
}

export const NewsAnnouncements: React.FC<NewsAnnouncementsProps> = ({ news }) => {
  if (!news || news.length === 0) return null;

  return (
    <section className="space-y-4 pt-8 border-t">
      <h3 className="text-xl font-bold tracking-tight">System Announcements & News</h3>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {news.slice(0, 2).map((item) => (
          <Card key={item.id} className="border-border/60">
            <CardHeader className="pb-2">
              <CardTitle className="text-base">{item.title}</CardTitle>
              <CardDescription className="text-xs">
                Published on {new Date(item.created_at).toLocaleDateString()}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground line-clamp-2">{item.content}</p>
            </CardContent>
          </Card>
        ))}
      </div>
    </section>
  );
};
