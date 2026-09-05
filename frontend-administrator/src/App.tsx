import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from '@/lib/auth';
import { ProtectedRoute } from '@/components/layout/ProtectedRoute';
import { AdminLayout } from '@/components/layout/AdminLayout';
import { Login } from '@/pages/Login';
import { Dashboard } from '@/pages/Dashboard';
import { Clients } from '@/pages/Clients';
import { Orders } from '@/pages/Orders';
import { Invoices } from '@/pages/Invoices';
import { Support } from '@/pages/Support';
import { Currencies } from '@/pages/Currencies';
import { News } from '@/pages/News';
import { MassMail } from '@/pages/MassMail';
import { AuditLogs } from '@/pages/AuditLogs';

export const App: React.FC = () => {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public Login Route */}
          <Route path="/login" element={<Login />} />

          {/* Protected Administrator Routes */}
          <Route element={<ProtectedRoute />}>
            <Route element={<AdminLayout />}>
              <Route path="/" element={<Dashboard />} />
              <Route path="/clients" element={<Clients />} />
              <Route path="/orders" element={<Orders />} />
              <Route path="/invoices" element={<Invoices />} />
              <Route path="/support" element={<Support />} />
              <Route path="/currencies" element={<Currencies />} />
              <Route path="/news" element={<News />} />
              <Route path="/mass-mail" element={<MassMail />} />
              <Route path="/audit-logs" element={<AuditLogs />} />
            </Route>
          </Route>

          {/* Catch-all redirect to Dashboard */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
};

export default App;
