import React, { useState, useEffect } from 'react';
import { Plus, Trash2, RefreshCw, AlertCircle } from 'lucide-react';
import { domainService } from '../../services/domain.service';
import { Button } from '@/components/ui/button';

interface ZoneEditorProps {
  domainId: number;
}

export const ZoneEditor: React.FC<ZoneEditorProps> = ({ domainId }) => {
  const [records, setRecords] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [newRecord, setNewRecord] = useState({
    type: 'A',
    name: '',
    content: '',
    ttl: 3600,
  });

  const fetchRecords = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await domainService.getDnsRecords(domainId);
      setRecords(data || []);
    } catch (err: any) {
      setError(err.message || 'Failed to load DNS records. Ensure domain is pointed to our nameservers.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRecords();
  }, [domainId]);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await domainService.addDnsRecord(domainId, newRecord);
      setNewRecord({ type: 'A', name: '', content: '', ttl: 3600 });
      fetchRecords();
    } catch (err: any) {
      alert(`Failed to add record: ${err.message}`);
    }
  };

  const handleDelete = async (recordId: string) => {
    if (!confirm('Are you sure you want to delete this record?')) return;
    try {
      await domainService.deleteDnsRecord(domainId, recordId);
      fetchRecords();
    } catch (err: any) {
      alert(`Failed to delete: ${err.message}`);
    }
  };

  return (
    <div className="space-y-4 pt-4">
      {error && (
        <div className="p-3 bg-amber-50 border border-amber-200 rounded-lg flex gap-2 text-xs text-amber-800">
          <AlertCircle className="h-4 w-4 shrink-0" />
          <p>{error}</p>
        </div>
      )}

      {/* Add Record Form */}
      <form onSubmit={handleAdd} className="grid grid-cols-12 gap-2 bg-muted/30 p-3 rounded-xl border border-dashed">
        <div className="col-span-2">
          <select
            value={newRecord.type}
            onChange={(e) => setNewRecord({ ...newRecord, type: e.target.value })}
            className="w-full h-9 rounded-lg border bg-background px-2 text-xs focus:ring-1 focus:ring-primary outline-none"
          >
            <option>A</option>
            <option>AAAA</option>
            <option>CNAME</option>
            <option>MX</option>
            <option>TXT</option>
            <option>SRV</option>
          </select>
        </div>
        <div className="col-span-3">
          <input
            type="text"
            placeholder="Name (@, www)"
            required
            value={newRecord.name}
            onChange={(e) => setNewRecord({ ...newRecord, name: e.target.value })}
            className="w-full h-9 rounded-lg border bg-background px-3 text-xs outline-none focus:ring-1 focus:ring-primary"
          />
        </div>
        <div className="col-span-5">
          <input
            type="text"
            placeholder="Content (IP, Target)"
            required
            value={newRecord.content}
            onChange={(e) => setNewRecord({ ...newRecord, content: e.target.value })}
            className="w-full h-9 rounded-lg border bg-background px-3 text-xs outline-none focus:ring-1 focus:ring-primary"
          />
        </div>
        <div className="col-span-2">
          <Button type="submit" size="sm" className="w-full gap-1 h-9">
            <Plus className="h-3.5 w-3.5" /> Add
          </Button>
        </div>
      </form>

      {/* Records Table */}
      <div className="border rounded-xl overflow-hidden shadow-sm">
        <table className="w-full text-left text-xs">
          <thead className="bg-muted/50 border-b">
            <tr>
              <th className="px-3 py-2 font-bold text-muted-foreground uppercase tracking-wider">Type</th>
              <th className="px-3 py-2 font-bold text-muted-foreground uppercase tracking-wider">Name</th>
              <th className="px-3 py-2 font-bold text-muted-foreground uppercase tracking-wider">Content</th>
              <th className="px-3 py-2 text-right"></th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {loading ? (
              <tr>
                <td colSpan={4} className="px-3 py-6 text-center text-muted-foreground italic">
                   <RefreshCw className="h-4 w-4 animate-spin mx-auto mb-1" /> Loading zone...
                </td>
              </tr>
            ) : records.length === 0 ? (
              <tr>
                <td colSpan={4} className="px-3 py-6 text-center text-muted-foreground">No records found.</td>
              </tr>
            ) : (
              records.map((r) => (
                <tr key={r.id} className="hover:bg-muted/10">
                  <td className="px-3 py-2 font-bold text-primary">{r.type}</td>
                  <td className="px-3 py-2 font-mono text-foreground">{r.name}</td>
                  <td className="px-3 py-2 text-muted-foreground truncate max-w-[200px]">{r.content}</td>
                  <td className="px-3 py-2 text-right">
                    <button
                      onClick={() => handleDelete(r.id)}
                      className="p-1.5 text-muted-foreground hover:text-destructive transition-colors"
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
