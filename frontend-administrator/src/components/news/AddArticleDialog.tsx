import React from 'react';
import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

interface AddArticleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: { title: string; content: string };
  setForm: React.Dispatch<React.SetStateAction<{ title: string; content: string }>>;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const AddArticleDialog: React.FC<AddArticleDialogProps> = ({
  open,
  onOpenChange,
  form,
  setForm,
  onSave,
  saving,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button size="sm" className="gap-1.5 shadow-sm">
          <Plus className="h-4 w-4" />
          New Announcement
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Publish Announcement</DialogTitle>
          <DialogDescription>Create a public announcement visible on client portal</DialogDescription>
        </DialogHeader>
        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Title</label>
            <Input
              required
              placeholder="e.g. Infrastructure Maintenance Notice"
              value={form.title}
              onChange={(e) => setForm((prev) => ({ ...prev, title: e.target.value }))}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Content (Markdown Supported)</label>
            <Textarea
              required
              rows={6}
              placeholder="Write announcement body..."
              value={form.content}
              onChange={(e) => setForm((prev) => ({ ...prev, content: e.target.value }))}
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Publishing...' : 'Publish'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
