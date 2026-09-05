import React from 'react';
import { Mail, RefreshCw, Send, AlertCircle } from 'lucide-react';
import { useMassMail } from '@/hooks/useMassMail';
import { AddCampaignDialog } from '@/components/massmail/AddCampaignDialog';
import { formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const MassMail: React.FC = () => {
  const {
    campaigns,
    loading,
    openModal,
    setOpenModal,
    form,
    setForm,
    saving,
    sendingId,
    message,
    fetchCampaigns,
    handleCreate,
    handleSend,
  } = useMassMail();

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

          <AddCampaignDialog
            open={openModal}
            onOpenChange={setOpenModal}
            form={form}
            setForm={setForm}
            onSave={handleCreate}
            saving={saving}
          />
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

export default MassMail;
