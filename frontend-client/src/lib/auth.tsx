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
}

interface ClientAuthContextType {
  user: ClientUser | null;
  token: string | null;
  balance: number;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
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
    const saved = localStorage.getItem('fossbilling_client_user');
    return saved ? JSON.parse(saved) : null;
  });
  const [token, setToken] = useState<string | null>(getStoredClientToken);
  const [balance, setBalance] = useState<number>(0);
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
      refreshProfile().finally(() => setIsLoading(false));
    } else {
      setIsLoading(false);
    }
  }, [token]);

  const login = async (email: string, password: string) => {
    setIsLoading(true);
    try {
      const res = await api.login(email, password);
      setStoredClientToken(res.token);
      setToken(res.token);
      setUser(res.client);
      localStorage.setItem('fossbilling_client_user', JSON.stringify(res.client));
      await refreshProfile();
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
      setStoredClientToken(res.token);
      setToken(res.token);
      setUser(res.client);
      localStorage.setItem('fossbilling_client_user', JSON.stringify(res.client));
      await refreshProfile();
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
        isLoading,
        login,
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

export const useClientAuth = () => {
  const context = useContext(ClientAuthContext);
  if (!context) {
    throw new Error('useClientAuth must be used within a ClientAuthProvider');
  }
  return context;
};
