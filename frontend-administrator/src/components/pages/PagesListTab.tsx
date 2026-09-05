import React from 'react';
import { FileText, Plus, Trash2, Edit } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { formatDate } from '@/lib/utils';
import { CustomPageItem } from '@/types/modules';

interface PagesListTabProps {
  pages: CustomPageItem[];
  loading: boolean;
  onEdit: (page: CustomPageItem) => void;
  onDelete: (id: number) => void;
  onNew: () => void;
}

export const PagesListTab: React.FC<PagesListTabProps> = ({
  pages,
  loading,
  onEdit,
  onDelete,
  onNew,
}) => {
  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <p className="text-xs text-muted-foreground">
          Static custom pages for Terms, Privacy Policies, SLAs, and Company Information.
        </p>
        <Button size="sm" onClick={onNew} className="gap-1.5 shadow-sm">
          <Plus className="h-4 w-4" />
          Create Page
        </Button>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Page Title</TableHead>
            <TableHead>URL Slug</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Last Updated</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {loading ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                Loading custom pages...
              </TableCell>
            </TableRow>
          ) : pages.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                No custom pages created.
              </TableCell>
            </TableRow>
          ) : (
            pages.map((p) => (
              <TableRow key={p.id}>
                <TableCell className="font-semibold text-sm flex items-center gap-2">
                  <FileText className="h-4 w-4 text-primary" />
                  {p.title}
                </TableCell>
                <TableCell className="font-mono text-xs text-muted-foreground">/{p.slug}</TableCell>
                <TableCell>
                  <Badge variant={p.published ? 'success' : 'secondary'}>
                    {p.published ? 'Published' : 'Draft'}
                  </Badge>
                </TableCell>
                <TableCell className="text-xs text-muted-foreground">
                  {p.updated_at ? formatDate(p.updated_at) : 'Recently'}
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex items-center justify-end gap-1">
                    <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => onEdit(p)}>
                      <Edit className="h-3.5 w-3.5 text-muted-foreground" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-7 w-7 text-muted-foreground hover:text-destructive"
                      onClick={() => onDelete(p.id)}
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
};
