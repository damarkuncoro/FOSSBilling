import React from 'react';
import { Server, RefreshCw, Activity, CheckCircle2, Zap } from 'lucide-react';
import { useServers } from '@/hooks/useServers';
import { AddServerDialog } from '@/components/servers/AddServerDialog';
import { ServersTable } from '@/components/servers/ServersTable';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export const Servers: React.FC = () => {
  const {
    servers,
    loading,
    openModal,
    setOpenModal,
    testingId,
    testResult,
    setTestResult,
    saving,
    form,
    setForm,
    fetchServers,
    handleCreate,
    handleTestConnection,
    handleDelete,
  } = useServers();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Hosting Server Managers</h1>
          <p className="text-sm text-muted-foreground">
            Connect cPanel/WHM, HestiaCP, CWP, and DirectAdmin servers for automated account provisioning.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={fetchServers} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>

          <AddServerDialog
            open={openModal}
            onOpenChange={setOpenModal}
            form={form}
            setForm={setForm}
            onSave={handleCreate}
            saving={saving}
          />
        </div>
      </div>

      {testResult && (
        <div
          className={`p-4 rounded-xl border flex items-center justify-between text-sm ${
            testResult.success
              ? 'bg-emerald-500/10 border-emerald-500/20 text-emerald-600 dark:text-emerald-400'
              : 'bg-destructive/10 border-destructive/20 text-destructive'
          }`}
        >
          <div className="flex items-center gap-2">
            <CheckCircle2 className="h-4 w-4 shrink-0" />
            <span>{testResult.message}</span>
          </div>
          <Button variant="ghost" size="sm" className="h-7 text-xs" onClick={() => setTestResult(null)}>
            Dismiss
          </Button>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="border-border/60">
          <CardContent className="p-5 flex items-center gap-4">
            <div className="h-10 w-10 rounded-xl bg-sky-500/10 text-sky-500 flex items-center justify-center">
              <Server className="h-5 w-5" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground font-medium">Total Servers</p>
              <p className="text-2xl font-bold">{servers.length}</p>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/60">
          <CardContent className="p-5 flex items-center gap-4">
            <div className="h-10 w-10 rounded-xl bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
              <Activity className="h-5 w-5" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground font-medium">Online Health</p>
              <p className="text-2xl font-bold">100%</p>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/60">
          <CardContent className="p-5 flex items-center gap-4">
            <div className="h-10 w-10 rounded-xl bg-indigo-500/10 text-indigo-500 flex items-center justify-center">
              <Zap className="h-5 w-5" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground font-medium">Allocated Accounts</p>
              <p className="text-2xl font-bold">
                {servers.reduce((sum, s) => sum + s.active_accounts, 0)} /{' '}
                {servers.reduce((sum, s) => sum + s.max_accounts, 0)}
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card className="border-border/60 shadow-sm">
        <CardHeader>
          <CardTitle className="text-base font-semibold">Registered Server Nodes</CardTitle>
          <CardDescription>Live servers configured for instant automatic hosting account creation</CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          <ServersTable
            servers={servers}
            loading={loading}
            testingId={testingId}
            onTestConnection={handleTestConnection}
            onDelete={handleDelete}
          />
        </CardContent>
      </Card>
    </div>
  );
};

export default Servers;
