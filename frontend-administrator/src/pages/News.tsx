import React from 'react';
import { Newspaper, RefreshCw, Trash2, Globe } from 'lucide-react';
import { useNews } from '@/hooks/useNews';
import { AddArticleDialog } from '@/components/news/AddArticleDialog';
import { formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';

export const News: React.FC = () => {
  const {
    articles,
    loading,
    openModal,
    setOpenModal,
    form,
    setForm,
    saving,
    fetchNews,
    handleCreate,
    handleDelete,
  } = useNews();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">News & Announcements CMS</h1>
          <p className="text-sm text-muted-foreground">
            Publish company updates, maintenance schedules, and customer announcements.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchNews} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          <AddArticleDialog
            open={openModal}
            onOpenChange={setOpenModal}
            form={form}
            setForm={setForm}
            onSave={handleCreate}
            saving={saving}
          />
        </div>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Published Articles ({articles.length})</CardTitle>
          <CardDescription>SEO-friendly slugs generated automatically with public guest endpoint access</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-16">ID</TableHead>
                <TableHead>Title</TableHead>
                <TableHead>Slug URL</TableHead>
                <TableHead>Published Date</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                    Loading articles...
                  </TableCell>
                </TableRow>
              ) : articles.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                    No articles published yet.
                  </TableCell>
                </TableRow>
              ) : (
                articles.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell className="font-mono text-xs font-semibold">#{item.id}</TableCell>
                    <TableCell className="font-semibold text-sm">
                      <div className="flex items-center gap-2">
                        <Newspaper className="h-4 w-4 text-primary shrink-0" />
                        <span>{item.title}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <span className="font-mono text-xs text-muted-foreground flex items-center gap-1">
                        <Globe className="h-3 w-3 text-muted-foreground" />
                        /news/{item.slug}
                      </span>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {formatDate(item.created_at)}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        size="icon"
                        variant="ghost"
                        className="h-7 w-7 text-muted-foreground hover:text-destructive"
                        onClick={() => handleDelete(item.id)}
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default News;
