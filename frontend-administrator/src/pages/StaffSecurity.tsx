import React from 'react';
import {
  Users,
  RefreshCw,
  Trash2,
  Lock,
} from 'lucide-react';
import { useStaffSecurity } from '@/hooks/useStaffSecurity';
import { AddStaffDialog } from '@/components/security/AddStaffDialog';
import { SecurityPolicyTab } from '@/components/security/SecurityPolicyTab';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';

export const StaffSecurity: React.FC = () => {
  const {
    activeTab,
    setActiveTab,
    staffList,
    securitySettings,
    setSecuritySettings,
    addStaffOpen,
    setAddStaffOpen,
    staffForm,
    setStaffForm,
    savingStaff,
    savingSecurity,
    blacklistText,
    setBlacklistText,
    saveSuccess,
    fetchData,
    handleCreateStaff,
    handleDeleteStaff,
    handleSaveSecurity,
  } = useStaffSecurity();

  const getRoleBadge = (role: string) => {
    switch (role) {
      case 'superadmin':
        return <Badge variant="default" className="text-[10px] bg-indigo-600">Super Admin</Badge>;
      case 'admin':
        return <Badge variant="secondary" className="text-[10px]">Administrator</Badge>;
      case 'support':
        return <Badge variant="outline" className="text-[10px] border-sky-500 text-sky-500">Support Staff</Badge>;
      case 'billing':
        return <Badge variant="outline" className="text-[10px] border-emerald-500 text-emerald-500">Billing Team</Badge>;
      default:
        return <Badge variant="outline" className="text-[10px]">{role}</Badge>;
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Staff Roles & Security Settings</h1>
          <p className="text-sm text-muted-foreground">
            Manage administrative access permissions, bot prevention (CAPTCHA), and brute-force defenses.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchData} className="gap-2">
            <RefreshCw className="h-4 w-4" />
            Refresh
          </Button>

          {activeTab === 'staff' && (
            <AddStaffDialog
              open={addStaffOpen}
              onOpenChange={setAddStaffOpen}
              staffForm={staffForm}
              setStaffForm={setStaffForm}
              onSave={handleCreateStaff}
              saving={savingStaff}
            />
          )}
        </div>
      </div>

      <div className="flex border-b">
        <button
          onClick={() => setActiveTab('staff')}
          className={`pb-3 px-4 text-sm font-semibold border-b-2 transition-all flex items-center gap-2 ${
            activeTab === 'staff'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Users className="h-4 w-4" />
          Staff Members ({staffList.length})
        </button>
        <button
          onClick={() => setActiveTab('security')}
          className={`pb-3 px-4 text-sm font-semibold border-b-2 transition-all flex items-center gap-2 ${
            activeTab === 'security'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Lock className="h-4 w-4" />
          Anti-Spam & Security
        </button>
      </div>

      {activeTab === 'staff' && (
        <Card className="border-border/60 shadow-sm">
          <CardHeader>
            <CardTitle className="text-base font-semibold">Administrative Team</CardTitle>
            <CardDescription>Users authorized to access the FOSSBilling management dashboard</CardDescription>
          </CardHeader>
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Staff Member</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Role & Permissions</TableHead>
                  <TableHead>Last Activity</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {staffList.map((staff) => (
                  <TableRow key={staff.id} className="hover:bg-muted/40 transition-colors">
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <Avatar className="h-8 w-8 border">
                          <AvatarFallback className="bg-primary/10 text-primary text-xs font-bold">
                            {staff.name.slice(0, 2).toUpperCase()}
                          </AvatarFallback>
                        </Avatar>
                        <span className="font-semibold text-sm">{staff.name}</span>
                      </div>
                    </TableCell>
                    <TableCell className="font-mono text-xs text-muted-foreground">{staff.email}</TableCell>
                    <TableCell>{getRoleBadge(staff.role)}</TableCell>
                    <TableCell className="text-xs text-muted-foreground">{staff.last_login}</TableCell>
                    <TableCell>
                      <Badge variant="success" className="text-[10px]">Active</Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      {staff.role !== 'superadmin' && (
                        <Button
                          size="icon"
                          variant="ghost"
                          className="h-7 w-7 text-muted-foreground hover:text-destructive"
                          onClick={() => handleDeleteStaff(staff.id)}
                        >
                          <Trash2 className="h-3.5 w-3.5" />
                        </Button>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}

      {activeTab === 'security' && (
        <SecurityPolicyTab
          securitySettings={securitySettings}
          setSecuritySettings={setSecuritySettings}
          blacklistText={blacklistText}
          setBlacklistText={setBlacklistText}
          onSave={handleSaveSecurity}
          saving={savingSecurity}
          saveSuccess={saveSuccess}
        />
      )}
    </div>
  );
};

export default StaffSecurity;
