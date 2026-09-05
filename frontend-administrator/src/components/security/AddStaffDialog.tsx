import React from 'react';
import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

interface AddStaffDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  staffForm: {
    name: string;
    email: string;
    role: 'superadmin' | 'admin' | 'support' | 'billing';
    password: string;
  };
  setStaffForm: React.Dispatch<
    React.SetStateAction<{
      name: string;
      email: string;
      role: 'superadmin' | 'admin' | 'support' | 'billing';
      password: string;
    }>
  >;
  onSave: (e: React.FormEvent) => void;
  saving: boolean;
}

export const AddStaffDialog: React.FC<AddStaffDialogProps> = ({
  open,
  onOpenChange,
  staffForm,
  setStaffForm,
  onSave,
  saving,
}) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button size="sm" className="gap-1.5 shadow-sm">
          <Plus className="h-4 w-4" />
          Add Staff Member
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Add New Staff Member</DialogTitle>
          <DialogDescription>Create an administrative login with assigned permissions</DialogDescription>
        </DialogHeader>

        <form onSubmit={onSave} className="space-y-4 pt-2">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Full Name</label>
            <Input
              required
              placeholder="e.g. Jane Doe"
              value={staffForm.name}
              onChange={(e) => setStaffForm((prev) => ({ ...prev, name: e.target.value }))}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Email Address</label>
            <Input
              type="email"
              required
              placeholder="jane@company.com"
              value={staffForm.email}
              onChange={(e) => setStaffForm((prev) => ({ ...prev, email: e.target.value }))}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Assigned Role</label>
            <select
              className="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
              value={staffForm.role}
              onChange={(e) => setStaffForm((prev) => ({ ...prev, role: e.target.value as any }))}
            >
              <option value="admin">Full Administrator</option>
              <option value="support">Support Agent (Tickets & Knowledgebase)</option>
              <option value="billing">Billing Manager (Invoices & Orders)</option>
              <option value="superadmin">Super Administrator (Full System Control)</option>
            </select>
          </div>
          <div className="space-y-1.5">
            <label className="text-xs font-semibold">Temporary Password</label>
            <Input
              type="password"
              required
              placeholder="••••••••••••"
              value={staffForm.password}
              onChange={(e) => setStaffForm((prev) => ({ ...prev, password: e.target.value }))}
            />
          </div>

          <div className="flex justify-end gap-2 pt-3 border-t">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? 'Adding...' : 'Create Staff Member'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
