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

export const App: React.FC = () => {
  return (
    <ClientAuthProvider>
      <CartProvider>
        <BrowserRouter>
          <Routes>
            <Route element={<ClientLayout />}>
              {/* Public Storefront & Cart Routes */}
              <Route path="/" element={<Storefront />} />
              <Route path="/cart" element={<Cart />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />

              {/* Protected Customer Routes */}
              <Route element={<ProtectedRoute />}>
                <Route path="/dashboard" element={<Dashboard />} />
                <Route path="/services" element={<Services />} />
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
