import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { guestApi } from '@/lib/api/guest';
import { Button } from '@/components/ui/button';
import { Home, ArrowLeft } from 'lucide-react';

export const NotFound: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    const checkRedirects = async () => {
      try {
        const res = await guestApi.lookupRedirect(location.pathname);
        if (res && res.target) {
          // If it's an external URL
          if (res.target.startsWith('http')) {
            window.location.href = res.target;
          } else {
            navigate(res.target, { replace: true });
          }
          return;
        }
      } catch (err) {
        // No redirect found or error
      } finally {
        setChecking(false);
      }
    };

    checkRedirects();
  }, [location.pathname, navigate]);

  if (checking) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] space-y-4">
        <div className="h-8 w-8 border-4 border-primary border-t-transparent rounded-full animate-spin" />
        <p className="text-muted-foreground animate-pulse">Routing your request...</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-[70vh] text-center px-6">
      <div className="space-y-2">
        <h1 className="text-8xl font-black text-primary/10">404</h1>
        <h2 className="text-2xl font-bold tracking-tight">Oops! Page not found.</h2>
        <p className="text-muted-foreground max-w-[400px] mx-auto">
          The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.
        </p>
      </div>

      <div className="flex flex-col sm:flex-row gap-3 mt-8">
        <Button variant="outline" onClick={() => navigate(-1)} className="gap-2">
          <ArrowLeft className="h-4 w-4" /> Go Back
        </Button>
        <Button onClick={() => navigate('/')} className="gap-2">
          <Home className="h-4 w-4" /> Back to Home
        </Button>
      </div>
    </div>
  );
};
