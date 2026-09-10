import React, { useState, useEffect } from 'react';
import { Mail, Edit2, Save, ArrowLeft, Info } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { systemApi } from '@/lib/api/system';
import { useToast } from '../hooks/use-toast';

export const EmailTemplates: React.FC = () => {
  const [templates, setTemplates] = useState<any[]>([]);
  const [selectedTemplate, setSelectedTemplate] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    fetchTemplates();
  }, []);

  const fetchTemplates = async () => {
    try {
      const data = await systemApi.listEmailTemplates();
      setTemplates(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (tpl: any) => {
    setSelectedTemplate({ ...tpl });
  };

  const handleSave = async () => {
    if (!selectedTemplate) return;
    setSaving(true);
    try {
      await systemApi.updateEmailTemplate(selectedTemplate);
      toast({ title: 'Success', description: 'Template updated successfully' });
      setSelectedTemplate(null);
      fetchTemplates();
    } catch (err) {
      toast({ title: 'Error', description: 'Failed to update template', variant: 'destructive' });
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-10 w-64" />
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[1, 2, 3, 4, 5, 6].map(i => <Skeleton key={i} className="h-32 rounded-xl" />)}
        </div>
      </div>
    );
  }

  if (selectedTemplate) {
    return (
      <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-300">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => setSelectedTemplate(null)}>
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">Edit Template: {selectedTemplate.code}</h1>
            <p className="text-sm text-muted-foreground">{selectedTemplate.description}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="lg:col-span-2 border-border/60 shadow-sm">
            <CardHeader>
              <CardTitle className="text-base">Template Content</CardTitle>
              <CardDescription>Customize the subject and HTML body of the email.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <label className="text-xs font-bold uppercase text-muted-foreground">Email Subject</label>
                <Input
                  value={selectedTemplate.subject}
                  onChange={e => setSelectedTemplate({ ...selectedTemplate, subject: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <label className="text-xs font-bold uppercase text-muted-foreground">HTML Body</label>
                <Textarea
                  className="min-h-[400px] font-mono text-xs"
                  value={selectedTemplate.content}
                  onChange={e => setSelectedTemplate({ ...selectedTemplate, content: e.target.value })}
                />
              </div>
              <Button className="w-full gap-2" onClick={handleSave} disabled={saving}>
                <Save className="h-4 w-4" />
                {saving ? 'Saving...' : 'Save Template Changes'}
              </Button>
            </CardContent>
          </Card>

          <Card className="border-border/60 shadow-sm">
            <CardHeader>
              <CardTitle className="text-base flex items-center gap-2">
                <Info className="h-4 w-4 text-indigo-500" />
                Variable Documentation
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="p-3 rounded-lg bg-indigo-50 dark:bg-indigo-950/30 border border-indigo-100 dark:border-indigo-900/50">
                <p className="text-xs font-medium text-indigo-900 dark:text-indigo-300">
                  You can use Go templates syntax <code>{"{{.Variable}}"}</code>.
                </p>
              </div>

              <div className="space-y-3">
                 <h4 className="text-xs font-bold uppercase">Common Variables</h4>
                 <ul className="text-xs space-y-2 text-muted-foreground">
                    <li><code>.AppName</code> - Site name</li>
                    <li><code>.FirstName</code> - Client first name</li>
                    <li><code>.Email</code> - Client email</li>
                 </ul>

                 <h4 className="text-xs font-bold uppercase pt-2">Invoice Specific</h4>
                 <ul className="text-xs space-y-2 text-muted-foreground">
                    <li><code>.InvoiceNr</code> - e.g. INV-1001</li>
                    <li><code>.Total</code> - Total amount</li>
                    <li><code>.Currency</code> - e.g. USD</li>
                 </ul>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-300">
      <div>
        <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2.5">
          <Mail className="h-7 w-7 text-primary" /> Email Templates
        </h1>
        <p className="text-sm text-muted-foreground">
          Manage and customize automated system emails sent to your clients and staff.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {templates.map((tpl) => (
          <Card key={tpl.id} className="group hover:border-primary/50 transition-all border-border/60 shadow-sm overflow-hidden">
            <CardHeader className="pb-3">
              <div className="flex justify-between items-start">
                <Badge variant="outline" className="font-mono text-[10px] uppercase">
                  {tpl.code}
                </Badge>
                <Button variant="ghost" size="icon" className="h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity" onClick={() => handleEdit(tpl)}>
                  <Edit2 className="h-3.5 w-3.5" />
                </Button>
              </div>
              <CardTitle className="text-sm font-semibold mt-2 line-clamp-1">{tpl.subject}</CardTitle>
              <CardDescription className="text-xs line-clamp-2 mt-1">{tpl.description}</CardDescription>
            </CardHeader>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default EmailTemplates;
