import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/lib/auth';
import { systemApi } from '@/lib/api/system';
import { setStoredToken } from '@/lib/api/client';

export function useLogin() {
  const [email, setEmail] = useState('admin@fossbilling.org');
  const [password, setPassword] = useState('admin123');
  const [twoFactorCode, setTwoFactorCode] = useState('');
  const [step, setStep] = useState<'login' | '2fa'>('login');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { refreshUser } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res: any = await systemApi.login(email, password);
      if (res.two_factor_required) {
        setStep('2fa');
      } else {
        setStoredToken(res.token);
        await refreshUser();
        navigate('/');
      }
    } catch (err: any) {
      setError(err.message || 'Invalid staff credentials');
    } finally {
      setLoading(false);
    }
  };

  const handleVerify2FA = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res: any = await systemApi.verifyTwoFactor(email, twoFactorCode);
      setStoredToken(res.token);
      await refreshUser();
      navigate('/');
    } catch (err: any) {
      setError(err.message || 'Invalid 2FA code');
    } finally {
      setLoading(false);
    }
  };

  return {
    email,
    setEmail,
    password,
    setPassword,
    twoFactorCode,
    setTwoFactorCode,
    step,
    setStep,
    error,
    loading,
    handleSubmit,
    handleVerify2FA,
  };
}

// Internal helper for useLogin since request is not exported from systemApi
async function request(path: string, opts: any) {
  const { request: baseRequest } = await import('../lib/api/client');
  return baseRequest(path, opts);
}
