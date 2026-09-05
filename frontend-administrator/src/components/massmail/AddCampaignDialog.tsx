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

interface AddCampaignDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: { subject: string; content: string };
  setForm: React.Dispatch<React.SetStateAction<{ subject: string; content: string }>>;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const AddCampaignDialog: React.FC<AddCampaignDialogProps> = ({
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
          New Campaign
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Draft Email Campaign</DialogTitle>
          <DialogDescription>Create a transactional or promotional broadcast message</DialogDescription>
        </DialogHeader>
        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Subject Line</label>
            <Input
              required
              placeholder="e.g. Exclusive Special Discounts for Cloud VPS"
              value={form.subject}
              onChange={(e) => setForm((prev) => ({ ...prev, subject: e.target.value }))}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Email Content Body</label>
            <Textarea
              required
              rows={6}
              placeholder="Write email body text..."
              value={form.content}
              onChange={(e) => setForm((prev) => ({ ...prev, content: e.target.value }))}
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Drafting...' : 'Save Draft'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
