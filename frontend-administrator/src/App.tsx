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
import { Reports } from '@/pages/Reports';
import { Support } from '@/pages/Support';
import { Products } from '@/pages/Products';
import { Licenses } from '@/pages/Licenses';
import { FormBuilder } from '@/pages/FormBuilder';
import { Domains } from '@/pages/Domains';
import { Servers } from '@/pages/Servers';
import { PaymentGateways } from '@/pages/PaymentGateways';
import { Coupons } from '@/pages/Coupons';
import { Currencies } from '@/pages/Currencies';
import { Pages } from '@/pages/Pages';
import { News } from '@/pages/News';
import { MassMail } from '@/pages/MassMail';
import { EmailTemplates } from '@/pages/EmailTemplates';
import { SeoSettings } from '@/pages/SeoSettings';
import { EmbedWidgets } from '@/pages/EmbedWidgets';
import { Company } from '@/pages/Company';
import { Extensions } from '@/pages/Extensions';
import { StaffSecurity } from '@/pages/StaffSecurity';
import { AntiSpam } from '@/pages/AntiSpam';
import { Webhooks } from '@/pages/Webhooks';
import { CookieConsent } from '@/pages/CookieConsent';
import { Redirects } from '@/pages/Redirects';
import { SystemHealth } from '@/pages/SystemHealth';
import { AuditLogs } from '@/pages/AuditLogs';

export const App: React.FC = () => {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />

          <Route element={<ProtectedRoute />}>
            <Route element={<AdminLayout />}>
              <Route path="/" element={<Dashboard />} />
              <Route path="/clients" element={<Clients />} />
              <Route path="/orders" element={<Orders />} />
              <Route path="/invoices" element={<Invoices />} />
              <Route path="/reports" element={<Reports />} />
              <Route path="/support" element={<Support />} />
              <Route path="/products" element={<Products />} />
              <Route path="/licenses" element={<Licenses />} />
              <Route path="/form-builder" element={<FormBuilder />} />
              <Route path="/domains" element={<Domains />} />
              <Route path="/servers" element={<Servers />} />
              <Route path="/gateways" element={<PaymentGateways />} />
              <Route path="/coupons" element={<Coupons />} />
              <Route path="/currencies" element={<Currencies />} />
              <Route path="/pages" element={<Pages />} />
              <Route path="/news" element={<News />} />
              <Route path="/mass-mail" element={<MassMail />} />
              <Route path="/email-templates" element={<EmailTemplates />} />
              <Route path="/seo" element={<SeoSettings />} />
              <Route path="/embed-widgets" element={<EmbedWidgets />} />
              <Route path="/company" element={<Company />} />
              <Route path="/extensions" element={<Extensions />} />
              <Route path="/staff-security" element={<StaffSecurity />} />
              <Route path="/antispam" element={<AntiSpam />} />
              <Route path="/webhooks" element={<Webhooks />} />
              <Route path="/cookie-consent" element={<CookieConsent />} />
              <Route path="/redirects" element={<Redirects />} />
              <Route path="/system" element={<SystemHealth />} />
              <Route path="/audit-logs" element={<AuditLogs />} />
            </Route>
          </Route>

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
};

export default App;
