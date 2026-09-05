import React from 'react';
import { XCircle, RefreshCw, LifeBuoy } from 'lucide-react';
import { Link } from 'react-router-dom';

export const PaymentFailed: React.FC = () => {
  return (
    <div className="max-w-md mx-auto px-4 py-16 text-center space-y-6">
      <div className="w-16 h-16 bg-rose-50 text-rose-600 rounded-full flex items-center justify-center mx-auto ring-8 ring-rose-50/50">
        <XCircle className="w-10 h-10" />
      </div>

      <div className="space-y-2">
        <h1 className="text-2xl font-black text-gray-900 tracking-tight">Payment Cancelled or Failed</h1>
        <p className="text-sm text-gray-500">
          We could not process your transaction. No funds were charged. You can retry paying with another method.
        </p>
      </div>

      <div className="flex flex-col gap-2 pt-2">
        <Link
          to="/invoices"
          className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold text-sm rounded-xl shadow-md transition-all"
        >
          <RefreshCw className="w-4 h-4" /> Retry Invoice Payment
        </Link>
        <Link
          to="/support"
          className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-white border border-gray-200 hover:bg-gray-50 text-gray-700 font-semibold text-sm rounded-xl transition-all"
        >
          <LifeBuoy className="w-4 h-4" /> Contact Billing Support
        </Link>
      </div>
    </div>
  );
};
