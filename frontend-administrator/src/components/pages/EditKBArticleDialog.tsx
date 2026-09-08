import React from 'react';
import { KnowledgebaseArticle } from '@/types/modules';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

interface EditKBArticleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: Partial<KnowledgebaseArticle>;
  onChange: (form: Partial<KnowledgebaseArticle>) => void;
  onSubmit: (e: React.FormEvent) => void;
}

export const EditKBArticleDialog: React.FC<EditKBArticleDialogProps> = ({
  open,
  onOpenChange,
  form,
  onChange,
  onSubmit,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>{form.id ? 'Edit Article' : 'New KB Article'}</DialogTitle>
          <DialogDescription>
            Write tutorials or troubleshooting guides for your clients.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 pt-2">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Article Title</label>
              <Input
                required
                value={form.title || ''}
                onChange={(e) => onChange({ ...form, title: e.target.value })}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">URL Slug</label>
              <Input
                required
                value={form.slug || ''}
                onChange={(e) => onChange({ ...form, slug: e.target.value })}
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Category</label>
            <Input
              required
              placeholder="e.g. Hosting, Billing"
              value={form.category || ''}
              onChange={(e) => onChange({ ...form, category: e.target.value })}
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Content (Markdown / HTML)</label>
            <Textarea
              required
              rows={8}
              value={form.content || ''}
              onChange={(e) => onChange({ ...form, content: e.target.value })}
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="kb_published"
              checked={form.published ?? true}
              onChange={(e) => onChange({ ...form, published: e.target.checked })}
              className="h-4 w-4 rounded border-gray-300 text-primary"
            />
            <label htmlFor="kb_published" className="text-xs font-medium cursor-pointer">
              Publish this article immediately
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-3 border-t">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">
              {form.id ? 'Update Article' : 'Publish Article'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
