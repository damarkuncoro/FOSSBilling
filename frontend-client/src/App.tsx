import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
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
import { Knowledgebase } from '@/pages/Knowledgebase';
import { News } from '@/pages/News';
import { Domains } from '@/pages/Domains';
import { Downloads } from '@/pages/Downloads';
import { Licenses } from '@/pages/Licenses';
import { PaymentSuccess } from '@/pages/PaymentSuccess';
import { PaymentFailed } from '@/pages/PaymentFailed';

export const App: React.FC = () => {
  return (
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

              {/* Protected Customer Routes */}
              <Route element={<ProtectedRoute />}>
                <Route path="/dashboard" element={<Dashboard />} />
                <Route path="/services" element={<Services />} />
                <Route path="/licenses" element={<Licenses />} />
                <Route path="/downloads" element={<Downloads />} />
                <Route path="/invoices" element={<Invoices />} />
                <Route path="/support" element={<Support />} />
                <Route path="/settings" element={<Settings />} />
              </Route>
            </Route>

            {/* Fallback */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </CartProvider>
    </ClientAuthProvider>
  );
};

export default App;
