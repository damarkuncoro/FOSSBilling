import React, { createContext, useContext, useEffect, useState } from 'react';
import { adminAuthService } from '@/services/admin_auth.service';
import { getStoredToken, removeStoredToken, setStoredToken } from './api';

export interface StaffUser {
  id: number;
  name: string;
  email: string;
  role: string;
  status: string;
}

interface AuthContextType {
  user: StaffUser | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  theme: 'light' | 'dark';
  toggleTheme: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<StaffUser | null>(() => {
    const saved = localStorage.getItem('fossbilling_admin_user');
    return saved ? JSON.parse(saved) : null;
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

  const login = async (email: string, password: string) => {
    setIsLoading(true);
    try {
      const res = await adminAuthService.login(email, password);
      setStoredToken(res.token);
      setToken(res.token);
      setUser(res.staff);
      localStorage.setItem('fossbilling_admin_user', JSON.stringify(res.staff));
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
