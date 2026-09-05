import React from 'react';
import { BookOpen, Plus, Eye } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatDate } from '@/lib/utils';
import { KnowledgebaseArticle } from '@/types/modules';

interface KnowledgebaseListTabProps {
  articles: KnowledgebaseArticle[];
  loading: boolean;
}

export const KnowledgebaseListTab: React.FC<KnowledgebaseListTabProps> = ({ articles, loading }) => {
  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <p className="text-xs text-muted-foreground">
          Self-help articles, tutorials, and documentation for customer troubleshooting.
        </p>
        <Button size="sm" onClick={() => alert('New Article dialog')} className="gap-1.5 shadow-sm">
          <Plus className="h-4 w-4" />
          Write KB Article
        </Button>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Article Title</TableHead>
            <TableHead>Category</TableHead>
            <TableHead>Views</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Last Updated</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {loading ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                Loading knowledgebase...
              </TableCell>
            </TableRow>
          ) : articles.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                No articles published yet.
              </TableCell>
            </TableRow>
          ) : (
            articles.map((art) => (
              <TableRow key={art.id}>
                <TableCell className="font-semibold text-sm flex items-center gap-2">
                  <BookOpen className="h-4 w-4 text-purple-500" />
                  {art.title}
                </TableCell>
                <TableCell>
                  <Badge variant="outline" className="text-xs">{art.category}</Badge>
                </TableCell>
                <TableCell className="font-mono text-xs">
                  <span className="flex items-center gap-1">
                    <Eye className="h-3 w-3 text-muted-foreground" />
                    {art.views}
                  </span>
                </TableCell>
                <TableCell>
                  <Badge variant={art.published ? 'success' : 'secondary'}>
                    {art.published ? 'Active' : 'Draft'}
                  </Badge>
                </TableCell>
                <TableCell className="text-xs text-muted-foreground">
                  {art.updated_at ? formatDate(art.updated_at) : 'Recently'}
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
};
