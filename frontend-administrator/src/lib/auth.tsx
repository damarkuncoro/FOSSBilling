import React, { createContext, useContext, useEffect, useState } from 'react';
import { adminAuthService } from '@/services/admin_auth.service';
import { getStoredToken, removeStoredToken, setStoredToken } from './api';

export interface StaffUser {
  id: number;
  name: string;
  email: string;
  role: string;
  status: string;
  two_factor_enabled?: boolean;
}

interface AuthContextType {
  user: StaffUser | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  completeLogin: (token: string, staff?: StaffUser) => void;
  refreshUser: () => Promise<void>;
  logout: () => void;
  theme: 'light' | 'dark';
  toggleTheme: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<StaffUser | null>(() => {
    try {
      const saved = localStorage.getItem('fossbilling_admin_user');
      if (!saved) return null;
      const parsed = JSON.parse(saved);
      return parsed && typeof parsed === 'object' ? parsed : null;
    } catch {
      return null;
    }
  });
  const [token, setToken] = useState<string | null>(getStoredToken);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [theme, setTheme] = useState<'light' | 'dark'>(() => {
    return (localStorage.getItem('fossbilling_theme') as 'light' | 'dark') || 'dark';
  });

  useEffect(() => {
    const root = window.document.documentElement;
    if (theme === 'dark') {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }
    localStorage.setItem('fossbilling_theme', theme);
  }, [theme]);

  useEffect(() => {
    // Initial token verification
    const currentToken = getStoredToken();
    if (currentToken && !user) {
      // Fallback staff session state
      const defaultUser: StaffUser = {
        id: 1,
        name: 'Administrator',
        email: 'admin@fossbilling.org',
        role: 'superadmin',
        status: 'active',
      };
      setUser(defaultUser);
      localStorage.setItem('fossbilling_admin_user', JSON.stringify(defaultUser));
    }
    setIsLoading(false);
  }, []);

  const refreshUser = async () => {
    // In a real app we'd fetch profile from API
    // For now we just keep the existing state
  };

  const completeLogin = (newToken: string, staffData?: StaffUser) => {
    setStoredToken(newToken);
    setToken(newToken);
    if (staffData) {
      setUser(staffData);
      localStorage.setItem('fossbilling_admin_user', JSON.stringify(staffData));
    }
  };

  const login = async (email: string, password: string) => {
    setIsLoading(true);
    try {
      const res = await adminAuthService.login(email, password);
      completeLogin(res.token, res.staff);
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    removeStoredToken();
    setUser(null);
    setToken(null);
  };

  const toggleTheme = () => {
    setTheme((prev) => (prev === 'light' ? 'dark' : 'light'));
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated: !!token && !!user,
        isLoading,
        login,
        completeLogin,
        refreshUser,
        logout,
        theme,
        toggleTheme,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
