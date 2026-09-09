import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { Shield, Lock, Mail, ArrowRight, AlertCircle, KeyRound } from 'lucide-react';
import { useClientAuth } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { guestApi } from '@/lib/api/guest';

export const Login: React.FC = () => {
  const [email, setEmail] = useState('client@fossbilling.org');
  const [password, setPassword] = useState('client123');
  const [twoFactorCode, setTwoFactorCode] = useState('');
  const [step, setStep] = useState<'login' | '2fa'>('login');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { completeLogin } = useClientAuth();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  useEffect(() => {
    const impersonateToken = searchParams.get('token');
    if (impersonateToken) {
      handleTokenLogin(impersonateToken);
    }
  }, [searchParams]);

  const handleTokenLogin = async (token: string) => {
    setLoading(true);
    try {
      // Fetch profile using the impersonate token to confirm it's valid and get user info
      // Since we don't have a direct 'me' with token override easily here without setting it,
      // we'll just use completeLogin with placeholder and then it will fetch profile.
      // But completeLogin in our system usually expects both.

      // Let's assume we can call an API with this token to get the user object
      const res = await fetch('/api/v1/client/profile', {
         headers: { 'Authorization': `Bearer ${token}` }
      });
      const json = await res.json();
      if (json.success) {
        await completeLogin(token, json.data);
        navigate('/dashboard');
      } else {
        setError('Impersonation token is invalid or expired.');
      }
    } catch (err) {
      setError('Failed to authenticate with impersonation token.');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res: any = await guestApi.login(email, password);
      if (res.two_factor_required) {
        setStep('2fa');
      } else {
        await completeLogin(res.token, res.client);
        navigate('/dashboard');
      }
    } catch (err: any) {
      setError(err.message || 'Invalid customer email or password');
    } finally {
      setLoading(false);
    }
  };

  const handleVerify2FA = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res: any = await guestApi.verifyTwoFactor(email, twoFactorCode);
      await completeLogin(res.token, res.client);
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.message || 'Invalid 2FA code');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[70vh] flex items-center justify-center py-12">
      <div className="w-full max-w-md space-y-6">
        <div className="text-center space-y-2">
          <div className="inline-flex h-12 w-12 rounded-2xl bg-gradient-to-tr from-primary to-indigo-500 items-center justify-center text-white shadow-lg shadow-primary/30 mb-1">
            <Shield className="h-6 w-6" />
          </div>
          <h1 className="text-2xl font-bold tracking-tight">Customer Portal Login</h1>
          <p className="text-sm text-muted-foreground">Access your active cloud hosting, invoices, and support tickets.</p>
        </div>

        <Card className="border-border/60 shadow-xl">
          <CardHeader>
            <CardTitle className="text-lg">
              {step === 'login' ? 'Sign In' : 'Verification Required'}
            </CardTitle>
            <CardDescription>
              {step === 'login'
                ? 'Enter your account credentials to continue'
                : 'Enter the 6-digit code from your authenticator app'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {error && (
              <div className="mb-4 p-3 rounded-lg bg-destructive/15 border border-destructive/20 text-destructive text-sm flex items-center gap-2">
                <AlertCircle className="h-4 w-4 shrink-0" />
                <span>{error}</span>
              </div>
            )}

            {step === 'login' ? (
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold text-muted-foreground">Email Address</label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                    <Input
                      type="email"
                      required
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      placeholder="you@example.com"
                      className="pl-9"
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-semibold text-muted-foreground">Password</label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
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

                <Button type="submit" className="w-full gap-2 font-semibold shadow-md shadow-primary/20" disabled={loading}>
                  {loading ? 'Authenticating...' : 'Sign In'}
                  {!loading && <ArrowRight className="h-4 w-4" />}
                </Button>
              </form>
            ) : (
              <form onSubmit={handleVerify2FA} className="space-y-4">
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold text-muted-foreground text-center block">Authenticator Code</label>
                  <div className="relative">
                    <KeyRound className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                    <Input
                      type="text"
                      required
                      autoFocus
                      value={twoFactorCode}
                      onChange={(e) => setTwoFactorCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                      placeholder="000000"
                      className="pl-9 font-mono text-center tracking-[0.5em] text-lg h-12"
                    />
                  </div>
                </div>

                <Button type="submit" className="w-full gap-2 font-semibold shadow-md shadow-primary/20" disabled={loading || twoFactorCode.length !== 6}>
                  {loading ? 'Verifying...' : 'Verify & Sign In'}
                </Button>

                <button
                  type="button"
                  onClick={() => setStep('login')}
                  className="w-full text-xs text-muted-foreground hover:text-foreground font-medium"
                >
                  ← Back to Login
                </button>
              </form>
            )}

            {step === 'login' && (
              <div className="mt-6 text-center text-xs text-muted-foreground">
                Don't have an account yet?{' '}
                <Link to="/register" className="text-primary font-semibold hover:underline">
                  Create Account
                </Link>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
