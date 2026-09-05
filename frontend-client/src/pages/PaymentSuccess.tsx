import React from 'react';
import { CheckCircle2, Package, FileText } from 'lucide-react';
import { Link, useSearchParams } from 'react-router-dom';

export const PaymentSuccess: React.FC = () => {
  const [params] = useSearchParams();
  const invoiceId = params.get('invoice_id') || '1042';

  return (
    <div className="max-w-md mx-auto px-4 py-16 text-center space-y-6">
      <div className="w-16 h-16 bg-emerald-50 text-emerald-600 rounded-full flex items-center justify-center mx-auto ring-8 ring-emerald-50/50">
        <CheckCircle2 className="w-10 h-10" />
      </div>

      <div className="space-y-2">
        <h1 className="text-2xl font-black text-gray-900 tracking-tight">Payment Completed!</h1>
        <p className="text-sm text-gray-500">
          Your payment for invoice #{invoiceId} has been successfully verified. Your service is now being provisioned automatically.
        </p>
      </div>

      <div className="p-4 bg-gray-50 rounded-2xl border border-gray-100 space-y-2 text-xs text-gray-600">
        <div className="flex justify-between">
          <span>Status</span>
          <span className="font-bold text-emerald-600">PAID & VERIFIED</span>
        </div>
        <div className="flex justify-between">
          <span>Provisioning Time</span>
          <span className="font-semibold text-gray-900">&lt; 60 seconds</span>
        </div>
      </div>

      <div className="flex flex-col gap-2 pt-2">
        <Link
          to="/services"
          className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm rounded-xl shadow-md transition-all"
        >
          <Package className="w-4 h-4" /> Go to My Active Services
        </Link>
        <Link
          to="/invoices"
          className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-white border border-gray-200 hover:bg-gray-50 text-gray-700 font-semibold text-sm rounded-xl transition-all"
        >
          <FileText className="w-4 h-4" /> View Invoice Receipt
        </Link>
      </div>
    </div>
  );
};
