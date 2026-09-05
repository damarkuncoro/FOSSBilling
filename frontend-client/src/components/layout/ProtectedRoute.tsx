import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useClientAuth } from '@/lib/auth';

export const ProtectedRoute: React.FC = () => {
  const { isAuthenticated, isLoading } = useClientAuth();

  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
          <p className="text-sm text-muted-foreground">Loading customer session...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
};
