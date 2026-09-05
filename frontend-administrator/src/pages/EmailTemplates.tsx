import React from 'react';
import {
  FileText,
  Server,
  RefreshCw,
  Edit,
  Search,
} from 'lucide-react';
import { useEmailTemplates } from '@/hooks/useEmailTemplates';
import { EditEmailTemplateDialog } from '@/components/email/EditEmailTemplateDialog';
import { SmtpSettingsCard } from '@/components/email/SmtpSettingsCard';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';

export const EmailTemplates: React.FC = () => {
  const {
    activeTab,
    setActiveTab,
    templates,
    mailConfig,
    setMailConfig,
    searchQuery,
    setSearchQuery,
    editModalOpen,
    setEditModalOpen,
    editingTemplate,
    tplForm,
    setTplForm,
    savingTpl,
    testEmail,
    setTestEmail,
    sendingTest,
    testStatus,
    savingConfig,
    fetchData,
    openEdit,
    handleSaveTemplate,
    handleSaveMailConfig,
    handleSendTestEmail,
    filteredTemplates,
  } = useEmailTemplates();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Email Templates & SMTP Settings</h1>
          <p className="text-sm text-muted-foreground">
            Customize automated system notification emails and configure outgoing SMTP delivery.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchData} className="gap-2">
          <RefreshCw className="h-4 w-4" />
          Refresh
        </Button>
      </div>

      <div className="flex border-b">
        <button
          onClick={() => setActiveTab('templates')}
          className={`pb-3 px-4 text-sm font-semibold border-b-2 transition-all flex items-center gap-2 ${
            activeTab === 'templates'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <FileText className="h-4 w-4" />
          System Templates ({templates.length})
        </button>
        <button
          onClick={() => setActiveTab('smtp')}
          className={`pb-3 px-4 text-sm font-semibold border-b-2 transition-all flex items-center gap-2 ${
            activeTab === 'smtp'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Server className="h-4 w-4" />
          SMTP Server & Delivery
        </button>
      </div>

      {activeTab === 'templates' && (
        <div className="space-y-4">
          <div className="relative w-full sm:w-80">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search templates, subjects..."
              className="pl-9 h-9"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>

          <Card className="border-border/60 shadow-sm">
            <CardHeader>
              <CardTitle className="text-base font-semibold">Automated Email Templates</CardTitle>
              <CardDescription>Twig markdown templates dispatched on system billing and support events</CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Template Code</TableHead>
                    <TableHead>Email Subject Line</TableHead>
                    <TableHead>Category</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Action</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredTemplates.map((tpl) => (
                    <TableRow key={tpl.id} className="hover:bg-muted/40 transition-colors">
                      <TableCell>
                        <div className="font-mono text-xs font-bold text-primary">{tpl.code}</div>
                        <p className="text-[11px] text-muted-foreground truncate max-w-[220px]">
                          {tpl.description}
                        </p>
                      </TableCell>
                      <TableCell className="text-sm font-medium">{tpl.subject}</TableCell>
                      <TableCell>
                        <Badge variant="secondary" className="capitalize text-[10px]">
                          {tpl.category}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        {tpl.enabled ? (
                          <Badge variant="success" className="text-[10px]">Active</Badge>
                        ) : (
                          <Badge variant="outline" className="text-[10px]">Disabled</Badge>
                        )}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="outline"
                          size="sm"
                          className="h-7 text-xs gap-1.5"
                          onClick={() => openEdit(tpl)}
                        >
                          <Edit className="h-3.5 w-3.5" />
                          Edit Template
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>
      )}

      {activeTab === 'smtp' && (
        <SmtpSettingsCard
          mailConfig={mailConfig}
          setMailConfig={setMailConfig}
          onSaveConfig={handleSaveMailConfig}
          savingConfig={savingConfig}
          testEmail={testEmail}
          setTestEmail={setTestEmail}
          onSendTest={handleSendTestEmail}
          sendingTest={sendingTest}
          testStatus={testStatus}
        />
      )}

      <EditEmailTemplateDialog
        open={editModalOpen}
        onOpenChange={setEditModalOpen}
        editingTemplate={editingTemplate}
        tplForm={tplForm}
        setTplForm={setTplForm}
        onSave={handleSaveTemplate}
        saving={savingTpl}
      />
    </div>
  );
};

export default EmailTemplates;
