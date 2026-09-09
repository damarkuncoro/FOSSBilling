import React, { createContext, useContext, useEffect, useState } from 'react';
import { api, getStoredClientToken, removeStoredClientToken, setStoredClientToken } from './api';

export interface ClientUser {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  currency?: string;
  company?: string;
  country?: string;
  status?: string;
  two_factor_enabled?: boolean;
}

interface ClientAuthContextType {
  user: ClientUser | null;
  token: string | null;
  balance: number;
  isAuthenticated: boolean;
  isImpersonated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  completeLogin: (token: string, user?: ClientUser) => Promise<void>;
  refreshUser: () => Promise<void>;
  register: (dto: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    currency?: string;
  }) => Promise<void>;
  logout: () => void;
  refreshProfile: () => Promise<void>;
  theme: 'light' | 'dark';
  toggleTheme: () => void;
}

const ClientAuthContext = createContext<ClientAuthContextType | undefined>(undefined);

export const ClientAuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<ClientUser | null>(() => {
    try {
      const saved = localStorage.getItem('fossbilling_client_user');
      if (!saved) return null;
      const parsed = JSON.parse(saved);
      return parsed && typeof parsed === 'object' ? parsed : null;
    } catch {
      return null;
    }
  });
  const [token, setToken] = useState<string | null>(getStoredClientToken);
  const [balance, setBalance] = useState<number>(0);
  const [isImpersonated, setIsImpersonated] = useState(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [theme, setTheme] = useState<'light' | 'dark'>(() => {
    return (localStorage.getItem('fossbilling_client_theme') as 'light' | 'dark') || 'dark';
  });

  useEffect(() => {
    const root = window.document.documentElement;
    if (theme === 'dark') {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }
    localStorage.setItem('fossbilling_client_theme', theme);
  }, [theme]);

  const refreshProfile = async () => {
    if (!token) return;
    try {
      const data = await api.getProfile();
      setUser(data.client);
      setBalance(data.balance || 0);
      localStorage.setItem('fossbilling_client_user', JSON.stringify(data.client));
    } catch {
      // ignore
    }
  };

  useEffect(() => {
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        setIsImpersonated(!!payload.is_impersonated);
      } catch {
        setIsImpersonated(false);
      }
      refreshProfile().finally(() => setIsLoading(false));
    } else {
      setIsImpersonated(false);
      setIsLoading(false);
    }
  }, [token]);

  const completeLogin = async (newToken: string, userData?: ClientUser) => {
    setStoredClientToken(newToken);
    setToken(newToken);
    if (userData) {
      setUser(userData);
      localStorage.setItem('fossbilling_client_user', JSON.stringify(userData));
    }
    // Fetch latest profile and balance
    try {
      const data = await api.getProfile();
      setUser(data.client);
      setBalance(data.balance || 0);
      localStorage.setItem('fossbilling_client_user', JSON.stringify(data.client));
    } catch (err) {
      console.error('Failed to fetch profile after login', err);
    }
  };

  const login = async (email: string, password: string) => {
    setIsLoading(true);
    try {
      const res = await api.login(email, password);
      await completeLogin(res.token, res.client);
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (dto: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    currency?: string;
  }) => {
    setIsLoading(true);
    try {
      const res = await api.register(dto);
      await completeLogin(res.token, res.client);
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    removeStoredClientToken();
    setUser(null);
    setToken(null);
    setBalance(0);
  };

  const toggleTheme = () => {
    setTheme((prev) => (prev === 'light' ? 'dark' : 'light'));
  };

  return (
    <ClientAuthContext.Provider
      value={{
        user,
        token,
        balance,
        isAuthenticated: !!token && !!user,
        isImpersonated,
        isLoading,
        login,
        completeLogin,
        refreshUser: refreshProfile,
        register,
        logout,
        refreshProfile,
        theme,
        toggleTheme,
      }}
    >
      {children}
    </ClientAuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(ClientAuthContext);
  if (!context) {
    throw new Error('useAuth must be used within a ClientAuthProvider');
  }
  return context;
};

export const useClientAuth = useAuth;
