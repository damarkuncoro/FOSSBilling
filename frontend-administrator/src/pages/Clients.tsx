import React, { useState } from 'react';
import { Search, RefreshCw, Mail, MapPin, UserPlus, Edit3, Trash2, FileSpreadsheet } from 'lucide-react';
import { useClients } from '@/hooks/useClients';
import { formatDate } from '@/lib/utils';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { AddClientDialog } from '@/components/clients/AddClientDialog';
import { EditClientDialog } from '@/components/clients/EditClientDialog';

export const Clients: React.FC = () => {
  const { loading, saving, search, setSearch, filtered, fetchClients, createClient, updateClient, deleteClient } = useClients();
  const [addOpen, setAddOpen] = useState(false);
  const [selectedClient, setSelectedClient] = useState<any | null>(null);

  const handleDelete = async (client: any) => {
    if (window.confirm(`Are you sure you want to delete client ${client.first_name} ${client.last_name} (#${client.id})?`)) {
      await deleteClient(client.id);
    }
  };

  const handleExportCSV = () => {
    if (!filtered.length) return;
    const headers = ['ID', 'First Name', 'Last Name', 'Email', 'Company', 'Country', 'Currency', 'Status', 'Registered'];
    const rows = filtered.map((c) => [c.id, c.first_name, c.last_name, c.email, c.company || '', c.country || '', c.currency || 'USD', c.status, c.created_at]);
    const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map((e) => e.join(','))].join('\n');
    const link = document.createElement('a');
    link.setAttribute('href', encodeURI(csvContent));
    link.setAttribute('download', `FOSSBilling-Clients-${new Date().toISOString().slice(0, 10)}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Client Management</h1>
          <p className="text-sm text-muted-foreground">
            Manage registered clients, inspect customer profiles, and export records.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={handleExportCSV} disabled={filtered.length === 0} className="gap-2">
            <FileSpreadsheet className="h-4 w-4" /> Export CSV
          </Button>
          <Button variant="outline" size="sm" onClick={fetchClients} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} /> Refresh
          </Button>
          <Button size="sm" onClick={() => setAddOpen(true)} className="gap-2">
            <UserPlus className="h-4 w-4" /> Add Client
          </Button>
        </div>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <CardTitle className="text-base font-semibold">Registered Clients ({filtered.length})</CardTitle>
              <CardDescription>Overview of all active and suspended client accounts</CardDescription>
            </div>
            <div className="relative w-full sm:w-72">
              <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by name or email..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9 h-9"
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-16">ID</TableHead>
                <TableHead>Customer</TableHead>
                <TableHead>Contact & Country</TableHead>
                <TableHead>Currency</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Registered</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">Loading clients...</TableCell>
                </TableRow>
              ) : filtered.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">No clients found.</TableCell>
                </TableRow>
              ) : (
                filtered.map((client) => (
                  <TableRow key={client.id}>
                    <TableCell className="font-mono text-xs font-semibold">#{client.id}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <div className="h-8 w-8 rounded-full bg-primary/10 text-primary font-bold flex items-center justify-center text-xs">
                          {client.first_name?.[0] || 'C'}
                        </div>
                        <div>
                          <p className="font-semibold text-sm">{client.first_name} {client.last_name}</p>
                          <p className="text-xs text-muted-foreground flex items-center gap-1">
                            <Mail className="h-3 w-3" /> {client.email}
                          </p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <span className="text-xs text-muted-foreground flex items-center gap-1">
                        <MapPin className="h-3 w-3" /> {client.country || 'Global'}
                      </span>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className="font-mono text-xs font-bold">{client.currency || 'USD'}</Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant={client.status === 'active' ? 'success' : 'destructive'}>{client.status || 'active'}</Badge>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">{formatDate(client.created_at)}</TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-1">
                        <Button variant="ghost" size="sm" className="h-8 w-8 p-0" onClick={() => setSelectedClient(client)} title="Edit Client">
                          <Edit3 className="h-3.5 w-3.5" />
                        </Button>
                        <Button variant="ghost" size="sm" className="h-8 w-8 p-0 text-destructive hover:bg-destructive/10" onClick={() => handleDelete(client)} title="Delete Client">
                          <Trash2 className="h-3.5 w-3.5" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <AddClientDialog open={addOpen} onOpenChange={setAddOpen} onSubmit={createClient} loading={saving} />
      <EditClientDialog
        client={selectedClient}
        open={!!selectedClient}
        onOpenChange={(open) => !open && setSelectedClient(null)}
        onSubmit={updateClient}
        loading={saving}
      />
    </div>
  );
};

export default Clients;
