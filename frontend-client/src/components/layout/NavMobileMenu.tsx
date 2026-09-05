import React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';

interface NavMobileMenuProps {
  open: boolean;
  onClose: () => void;
  isAuthenticated: boolean;
  onLogout: () => void;
}

export const NavMobileMenu: React.FC<NavMobileMenuProps> = ({
  open,
  onClose,
  isAuthenticated,
  onLogout,
}) => {
  if (!open) return null;

  return (
    <div className="md:hidden border-b bg-card p-4 space-y-2">
      <Link to="/" onClick={onClose} className="block text-sm font-medium py-1 text-gray-800 dark:text-gray-200">
        Store
      </Link>
      <Link to="/domains" onClick={onClose} className="block text-sm font-medium py-1 text-gray-800 dark:text-gray-200">
        Domains
      </Link>
      <Link to="/kb" onClick={onClose} className="block text-sm font-medium py-1 text-gray-800 dark:text-gray-200">
        Knowledgebase
      </Link>
      <Link to="/news" onClick={onClose} className="block text-sm font-medium py-1 text-gray-800 dark:text-gray-200">
        News & Updates
      </Link>
      {isAuthenticated ? (
        <>
          <div className="pt-2 border-t border-gray-200 dark:border-gray-800 space-y-1">
            <Link to="/dashboard" onClick={onClose} className="block text-sm font-medium py-1">
              Dashboard
            </Link>
            <Link to="/services" onClick={onClose} className="block text-sm font-medium py-1">
              My Services
            </Link>
            <Link to="/licenses" onClick={onClose} className="block text-sm font-medium py-1">
              Software Licenses
            </Link>
            <Link to="/downloads" onClick={onClose} className="block text-sm font-medium py-1">
              Downloads
            </Link>
            <Link to="/invoices" onClick={onClose} className="block text-sm font-medium py-1">
              Invoices
            </Link>
            <Link to="/support" onClick={onClose} className="block text-sm font-medium py-1">
              Support Tickets
            </Link>
            <Link to="/settings" onClick={onClose} className="block text-sm font-medium py-1">
              Account Settings
            </Link>
          </div>
          <Button
            variant="destructive"
            size="sm"
            onClick={() => {
              onLogout();
              onClose();
            }}
            className="w-full mt-2"
          >
            Logout
          </Button>
        </>
      ) : (
        <div className="flex gap-2 pt-2 border-t border-gray-200 dark:border-gray-800">
          <Link to="/login" className="flex-1" onClick={onClose}>
            <Button variant="outline" className="w-full">
              Sign In
            </Button>
          </Link>
          <Link to="/register" className="flex-1" onClick={onClose}>
            <Button className="w-full">Register</Button>
          </Link>
        </div>
      )}
    </div>
  );
};
