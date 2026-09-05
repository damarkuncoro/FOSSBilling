import React from 'react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

interface EditPageDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: any;
  onChange: (form: any) => void;
  onSubmit: (e: React.FormEvent) => void;
}

export const EditPageDialog: React.FC<EditPageDialogProps> = ({
  open,
  onOpenChange,
  form,
  onChange,
  onSubmit,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xl">
        <DialogHeader>
          <DialogTitle>{form.id ? 'Edit Custom Page' : 'Create New Page'}</DialogTitle>
          <DialogDescription>Markdown and HTML supported static content</DialogDescription>
        </DialogHeader>
        <form onSubmit={onSubmit} className="space-y-4 pt-2">
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="text-xs font-semibold">Page Title</label>
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
            <label className="text-xs font-semibold">Page Content (Markdown / HTML)</label>
            <Textarea
              required
              rows={8}
              value={form.content || ''}
              onChange={(e) => onChange({ ...form, content: e.target.value })}
            />
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Save Page</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
