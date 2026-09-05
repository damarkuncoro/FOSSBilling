import React, { useState, useEffect } from 'react';
import { Shield, X } from 'lucide-react';

export const CookieConsentBanner: React.FC = () => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const consent = localStorage.getItem('fossbilling_cookie_consent');
    if (!consent) {
      setIsVisible(true);
    }
  }, []);

  const handleDecision = (decision: 'accepted' | 'declined') => {
    localStorage.setItem('fossbilling_cookie_consent', decision);
    setIsVisible(false);
  };

  if (!isVisible) return null;

  return (
    <div className="fixed bottom-4 left-4 right-4 sm:left-auto sm:right-6 sm:max-w-md z-50 animate-in slide-in-from-bottom-5 duration-300">
      <div className="bg-gray-900 text-white p-5 rounded-2xl shadow-2xl border border-gray-800 space-y-3">
        <div className="flex items-start gap-3">
          <div className="p-2 bg-indigo-600/20 text-indigo-400 rounded-xl shrink-0">
            <Shield className="w-5 h-5" />
          </div>
          <div className="space-y-1">
            <h4 className="text-xs font-bold uppercase tracking-wider text-indigo-300">Cookie Notice</h4>
            <p className="text-xs text-gray-300 leading-relaxed">
              We use cookies to maintain your login session, cart contents, and enhance platform security.
            </p>
          </div>
          <button onClick={() => handleDecision('declined')} className="text-gray-400 hover:text-gray-200">
            <X className="w-4 h-4" />
          </button>
        </div>

        <div className="flex items-center justify-end gap-2 pt-1">
          <button
            onClick={() => handleDecision('declined')}
            className="px-3 py-1.5 text-xs font-medium text-gray-400 hover:text-white bg-white/5 hover:bg-white/10 rounded-lg transition-colors"
          >
            Essential Only
          </button>
          <button
            onClick={() => handleDecision('accepted')}
            className="px-4 py-1.5 text-xs font-bold text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-sm transition-colors"
          >
            Accept All
          </button>
        </div>
      </div>
    </div>
  );
};
