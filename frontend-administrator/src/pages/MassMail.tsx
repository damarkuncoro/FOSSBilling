import React, { useEffect, useState } from 'react';
import { Mail, RefreshCw, Plus, Send, CheckCircle, AlertCircle } from 'lucide-react';
import { api } from '@/lib/api';
import { formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';

export const MassMail: React.FC = () => {
  const [campaigns, setCampaigns] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState({ subject: '', content: '' });
  const [saving, setSaving] = useState(false);
  const [sendingId, setSendingId] = useState<number | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  const fetchCampaigns = async () => {
    setLoading(true);
    try {
      const data = await api.getMassMailCampaigns();
      setCampaigns(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCampaigns();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await api.createMassMailCampaign(form);
      setOpenModal(false);
      setForm({ subject: '', content: '' });
      await fetchCampaigns();
    } catch (err) {
      console.error(err);
    } finally {
      setSaving(false);
    }
  };

  const handleSend = async (id: number) => {
    setSendingId(id);
    setMessage(null);
    try {
      await api.sendMassMailCampaign(id);
      setMessage(`Campaign #${id} has been broadcasted to all active clients!`);
      await fetchCampaigns();
    } catch (err: any) {
      setMessage(`Error sending campaign: ${err.message}`);
    } finally {
      setSendingId(null);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Mass Mailer & Campaigns</h1>
          <p className="text-sm text-muted-foreground">
            Create promotional newsletters and broadcast email campaigns to registered customers.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchCampaigns} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          <Dialog open={openModal} onOpenChange={setOpenModal}>
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
              <form onSubmit={handleCreate} className="space-y-4 pt-2">
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">Subject Line</label>
                  <Input
                    required
                    placeholder="e.g. Exclusive Special Discounts for Cloud VPS"
                    value={form.subject}
                    onChange={(e) => setForm({ ...form, subject: e.target.value })}
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold">Email Content Body</label>
                  <Textarea
                    required
                    rows={6}
                    placeholder="Write email body text..."
                    value={form.content}
                    onChange={(e) => setForm({ ...form, content: e.target.value })}
                  />
                </div>
                <div className="flex justify-end gap-2 pt-2">
                  <Button type="button" variant="outline" onClick={() => setOpenModal(false)}>
                    Cancel
                  </Button>
                  <Button type="submit" disabled={saving}>
                    {saving ? 'Drafting...' : 'Save Draft'}
                  </Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {message && (
        <div className="p-4 rounded-xl bg-primary/10 border border-primary/20 text-primary text-sm flex items-center gap-2">
          <AlertCircle className="h-4 w-4 shrink-0" />
          <span>{message}</span>
        </div>
      )}

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Campaigns History ({campaigns.length})</CardTitle>
          <CardDescription>Track send status, recipient delivery counts, and dispatch triggers</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-16">ID</TableHead>
                <TableHead>Subject</TableHead>
                <TableHead>Target</TableHead>
                <TableHead>Sent Count</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Created</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                    Loading campaigns...
                  </TableCell>
                </TableRow>
              ) : campaigns.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                    No email campaigns found.
                  </TableCell>
                </TableRow>
              ) : (
                campaigns.map((camp) => (
                  <TableRow key={camp.id}>
                    <TableCell className="font-mono text-xs font-semibold">#{camp.id}</TableCell>
                    <TableCell className="font-medium text-sm">
                      <div className="flex items-center gap-2">
                        <Mail className="h-4 w-4 text-primary shrink-0" />
                        <span>{camp.subject}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className="text-[10px]">
                        {camp.target_group || 'All Clients'}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-semibold text-xs">{camp.sent_count || 0}</TableCell>
                    <TableCell>
                      <Badge variant={camp.status === 'completed' ? 'success' : 'warning'}>
                        {camp.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {formatDate(camp.created_at)}
                    </TableCell>
                    <TableCell className="text-right">
                      {camp.status !== 'completed' && (
                        <Button
                          size="sm"
                          className="h-7 text-xs gap-1.5"
                          disabled={sendingId === camp.id}
                          onClick={() => handleSend(camp.id)}
                        >
                          <Send className="h-3 w-3" />
                          {sendingId === camp.id ? 'Sending...' : 'Broadcast'}
                        </Button>
                      )}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};
