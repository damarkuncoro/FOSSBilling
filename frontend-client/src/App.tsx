import React from 'react';
import { BrowserRouter, Routes, Route, Navigate, useParams, useNavigate } from 'react-router-dom';
import { ClientAuthProvider } from '@/lib/auth';
import { CartProvider } from '@/lib/cart';
import { ClientLayout } from '@/components/layout/ClientLayout';
import { ProtectedRoute } from '@/components/layout/ProtectedRoute';
import { Storefront } from '@/pages/Storefront';
import { Login } from '@/pages/Login';
import { Register } from '@/pages/Register';
import { Dashboard } from '@/pages/Dashboard';
import { Cart } from '@/pages/Cart';
import { Services } from '@/pages/Services';
import { Invoices } from '@/pages/Invoices';
import { Support } from '@/pages/Support';
import { Settings } from '@/pages/Settings';
import { Affiliate } from '@/pages/Affiliate';
import { Knowledgebase } from '@/pages/Knowledgebase';
import { News } from '@/pages/News';
import { Domains } from '@/pages/Domains';
import { Downloads } from '@/pages/Downloads';
import { Licenses } from '@/pages/Licenses';
import { PaymentSuccess } from '@/pages/PaymentSuccess';
import { PaymentFailed } from '@/pages/PaymentFailed';
import { NotFound } from '@/pages/NotFound';

const ReferralRedirect: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  React.useEffect(() => {
    if (id) {
      localStorage.setItem('fb_referrer_id', id);
    }
    navigate('/register', { replace: true });
  }, [id, navigate]);
  return null;
};

import { I18nProvider } from '@/lib/i18n';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';

const ReferrerTracker: React.FC = () => {
  React.useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const ref = params.get('ref');
    if (ref) {
      localStorage.setItem('fb_referrer_id', ref);
    }
  }, []);
  return null;
};

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 1,
    },
  },
});

export const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <I18nProvider>
        <ReferrerTracker />
        <ClientAuthProvider>
          <CartProvider>
            <BrowserRouter>
              <Routes>
                <Route element={<ClientLayout />}>
                  {/* Public Storefront, Knowledgebase, News & Cart */}
                  <Route path="/" element={<Storefront />} />
                  <Route path="/kb" element={<Knowledgebase />} />
                  <Route path="/news" element={<News />} />
                  <Route path="/domains" element={<Domains />} />
                  <Route path="/cart" element={<Cart />} />
                  <Route path="/payment/success" element={<PaymentSuccess />} />
                  <Route path="/payment/failed" element={<PaymentFailed />} />
                  <Route path="/login" element={<Login />} />
                  <Route path="/register" element={<Register />} />
                  <Route path="/r/:id" element={<ReferralRedirect />} />

                  {/* Protected Customer Routes */}
                  <Route element={<ProtectedRoute />}>
                    <Route path="/dashboard" element={<Dashboard />} />
                    <Route path="/services" element={<Services />} />
                    <Route path="/licenses" element={<Licenses />} />
                    <Route path="/downloads" element={<Downloads />} />
                    <Route path="/invoices" element={<Invoices />} />
                    <Route path="/support" element={<Support />} />
                    <Route path="/affiliate" element={<Affiliate />} />
                    <Route path="/settings" element={<Settings />} />
                  </Route>
                </Route>

                {/* Fallback */}
                <Route path="*" element={<NotFound />} />
              </Routes>
            </BrowserRouter>
          </CartProvider>
        </ClientAuthProvider>
      </I18nProvider>
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
};

export default App;
