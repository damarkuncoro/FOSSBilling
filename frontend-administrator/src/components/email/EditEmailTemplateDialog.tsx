import React from 'react';
import { EmailTemplateItem } from '@/lib/api';
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

interface EditEmailTemplateDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  editingTemplate: EmailTemplateItem | null;
  tplForm: Partial<EmailTemplateItem>;
  setTplForm: React.Dispatch<React.SetStateAction<Partial<EmailTemplateItem>>>;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const EditEmailTemplateDialog: React.FC<EditEmailTemplateDialogProps> = ({
  open,
  onOpenChange,
  editingTemplate,
  tplForm,
  setTplForm,
  onSave,
  saving,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Edit Template: {editingTemplate?.code}</DialogTitle>
          <DialogDescription>
            Modify subject and body text. Click variables below to insert them into your template.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Subject Line</label>
            <Input
              required
              value={tplForm.subject || ''}
              onChange={(e) => setTplForm((prev) => ({ ...prev, subject: e.target.value }))}
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Email Content (Markdown / Twig)</label>
            <Textarea
              rows={9}
              className="font-mono text-xs leading-relaxed"
              value={tplForm.content || ''}
              onChange={(e) => setTplForm((prev) => ({ ...prev, content: e.target.value }))}
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-[11px] font-semibold text-muted-foreground">Available Variables</label>
            <div className="flex flex-wrap gap-1.5">
              {editingTemplate?.variables?.map((v) => (
                <button
                  key={v}
                  type="button"
                  onClick={() =>
                    setTplForm((prev) => ({ ...prev, content: (prev.content || '') + ` {{ ${v} }}` }))
                  }
                  className="text-[10px] font-mono bg-muted hover:bg-primary/20 hover:text-primary px-2 py-0.5 rounded border transition-colors"
                >
                  {`{{ ${v} }}`}
                </button>
              ))}
            </div>
          </div>

          <div className="flex items-center gap-2 pt-2">
            <input
              type="checkbox"
              id="tpl_enabled"
              checked={tplForm.enabled ?? true}
              onChange={(e) => setTplForm((prev) => ({ ...prev, enabled: e.target.checked }))}
              className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
            />
            <label htmlFor="tpl_enabled" className="text-xs font-medium cursor-pointer">
              Enable automated sending for this event
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-3 border-t">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Saving...' : 'Save Template'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
