import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Shield, Mail, Lock, ArrowRight, AlertCircle } from 'lucide-react';
import { useClientAuth } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

export const Register: React.FC = () => {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [currency, setCurrency] = useState('USD');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { register } = useClientAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const referrerId = localStorage.getItem('fb_referrer_id');
      await register({
        email,
        password,
        first_name: firstName,
        last_name: lastName,
        currency,
        referrer_id: referrerId ? parseInt(referrerId) : undefined,
      });
      localStorage.removeItem('fb_referrer_id');
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.message || 'Registration failed');
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
          <h1 className="text-2xl font-bold tracking-tight">Create Customer Account</h1>
          <p className="text-sm text-muted-foreground">Register to order cloud hosting and digital licenses.</p>
        </div>

        <Card className="border-border/60 shadow-xl">
          <CardHeader>
            <CardTitle className="text-lg">Register</CardTitle>
            <CardDescription>Fill in your personal details to get started</CardDescription>
          </CardHeader>
          <CardContent>
            {error && (
              <div className="mb-4 p-3 rounded-lg bg-destructive/15 border border-destructive/20 text-destructive text-sm flex items-center gap-2">
                <AlertCircle className="h-4 w-4 shrink-0" />
                <span>{error}</span>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold text-muted-foreground">First Name</label>
                  <Input
                    required
                    value={firstName}
                    onChange={(e) => setFirstName(e.target.value)}
                    placeholder="Budi"
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold text-muted-foreground">Last Name</label>
                  <Input
                    required
                    value={lastName}
                    onChange={(e) => setLastName(e.target.value)}
                    placeholder="Santoso"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground">Email Address</label>
                <div className="relative">
                  <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <Input
                    type="email"
                    required
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="budi@example.com"
                    className="pl-9"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground flex justify-between items-center">
                   <span>Password (Min. 8 chars)</span>
                   {password.length > 0 && (
                     <span className={`text-[10px] uppercase font-bold ${
                       password.length < 8 ? 'text-rose-500' :
                       password.length < 12 ? 'text-amber-500' : 'text-emerald-500'
                     }`}>
                        {password.length < 8 ? 'Weak' : password.length < 12 ? 'Medium' : 'Strong'}
                     </span>
                   )}
                </label>
                <div className="relative">
                  <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <Input
                    type="password"
                    required
                    minLength={8}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="••••••••"
                    className="pl-9"
                  />
                </div>
                {password.length > 0 && (
                  <div className="flex gap-1 mt-1">
                    <div className={`h-1 flex-1 rounded-full transition-colors ${password.length >= 1 ? (password.length < 8 ? 'bg-rose-500' : 'bg-emerald-500') : 'bg-muted'}`} />
                    <div className={`h-1 flex-1 rounded-full transition-colors ${password.length >= 8 ? (password.length < 12 ? 'bg-amber-500' : 'bg-emerald-500') : 'bg-muted'}`} />
                    <div className={`h-1 flex-1 rounded-full transition-colors ${password.length >= 12 ? 'bg-emerald-500' : 'bg-muted'}`} />
                  </div>
                )}
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground">Preferred Currency</label>
                <select
                  value={currency}
                  onChange={(e) => setCurrency(e.target.value)}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus:outline-none focus:ring-1 focus:ring-ring"
                >
                  <option value="USD">USD - US Dollar ($)</option>
                  <option value="IDR">IDR - Indonesian Rupiah (Rp)</option>
                  <option value="EUR">EUR - Euro (€)</option>
                  <option value="SGD">SGD - Singapore Dollar (S$)</option>
                </select>
              </div>

              <Button type="submit" className="w-full gap-2 font-semibold shadow-md shadow-primary/20" disabled={loading}>
                {loading ? 'Creating Account...' : 'Complete Registration'}
                {!loading && <ArrowRight className="h-4 w-4" />}
              </Button>

              <div className="relative my-4">
                <div className="absolute inset-0 flex items-center"><span className="w-full border-t border-muted" /></div>
                <div className="relative flex justify-center text-[10px] uppercase"><span className="bg-card px-2 text-muted-foreground font-bold">Or register with</span></div>
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
                Sign Up with Google
              </Button>
            </form>

            <div className="mt-6 text-center text-xs text-muted-foreground">
              Already have an account?{' '}
              <Link to="/login" className="text-primary font-semibold hover:underline">
                Sign In
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
