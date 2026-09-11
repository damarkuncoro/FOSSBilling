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

                <div className="relative my-4">
                  <div className="absolute inset-0 flex items-center"><span className="w-full border-t border-muted" /></div>
                  <div className="relative flex justify-center text-[10px] uppercase"><span className="bg-card px-2 text-muted-foreground font-bold">Or continue with</span></div>
                </div>

                <Button
                  type="button"
                  variant="outline"
                  className="w-full font-bold"
                  onClick={() => window.location.href = '/api/v1/guest/auth/google'}
                >
                  <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24">
                    <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
                    <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
                    <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
                    <path d="M12 5.38c1.62 0 3.06.56 4.21 1.66l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
                  </svg>
                  Login with Google
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
