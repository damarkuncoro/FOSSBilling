import React from 'react';
import { ShieldCheck, Lock, Mail, ArrowRight, AlertCircle } from 'lucide-react';
import { useLogin } from '@/hooks/useLogin';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

export const Login: React.FC = () => {
  const {
    email,
    setEmail,
    password,
    setPassword,
    error,
    loading,
    handleSubmit,
  } = useLogin();

  return (
    <div className="min-h-screen w-full flex items-center justify-center p-4 bg-gradient-to-br from-background via-muted/30 to-background">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center space-y-2">
          <div className="inline-flex h-14 w-14 rounded-2xl bg-gradient-to-tr from-primary to-indigo-500 items-center justify-center text-white shadow-xl shadow-primary/30 mb-2">
            <ShieldCheck className="h-8 w-8" />
          </div>
          <h1 className="text-2xl font-bold tracking-tight">FOSSBilling Admin</h1>
          <p className="text-sm text-muted-foreground flex items-center justify-center gap-1.5">
            Powered by Golang High-Performance API
            <Badge variant="success" className="text-[10px]">v1.0</Badge>
          </p>
        </div>

        <Card className="border-border/60 shadow-xl backdrop-blur-xl bg-card/80">
          <CardHeader className="space-y-1">
            <CardTitle className="text-xl">Sign in to control panel</CardTitle>
            <CardDescription>Enter your administrator email and password</CardDescription>
          </CardHeader>
          <CardContent>
            {error && (
              <div className="mb-4 p-3 rounded-lg bg-destructive/15 border border-destructive/20 text-destructive text-sm flex items-center gap-2">
                <AlertCircle className="h-4 w-4 shrink-0" />
                <span>{error}</span>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground">Staff Email</label>
                <div className="relative">
                  <Mail className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    type="email"
                    required
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="admin@fossbilling.org"
                    className="pl-9"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground">Password</label>
                <div className="relative">
                  <Lock className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    type="password"
                    required
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="••••••••"
                    className="pl-9"
                  />
                </div>
              </div>

              <Button type="submit" className="w-full gap-2 font-semibold shadow-lg shadow-primary/25" disabled={loading}>
                {loading ? 'Authenticating...' : 'Sign In'}
                {!loading && <ArrowRight className="h-4 w-4" />}
              </Button>
            </form>

            <div className="mt-6 p-3 rounded-lg bg-muted/50 border text-xs text-muted-foreground space-y-1">
              <p className="font-semibold text-foreground">💡 Default Demo Credentials:</p>
              <p>Email: <code className="text-primary">admin@fossbilling.org</code></p>
              <p>Password: <code className="text-primary">admin123</code></p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Login;
