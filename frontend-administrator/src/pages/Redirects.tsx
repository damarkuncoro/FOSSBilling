import React from 'react';
import { ArrowRightLeft, Plus, Search } from 'lucide-react';
import { useRedirects } from '../hooks/useRedirects';
import { RedirectTable } from '../components/redirects/RedirectTable';
import { AddRedirectDialog } from '../components/redirects/AddRedirectDialog';

export const Redirects: React.FC = () => {
  const {
    redirects,
    search,
    setSearch,
    isAddOpen,
    setIsAddOpen,
    createRedirect,
    toggleRedirect,
    deleteRedirect,
  } = useRedirects();

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
            <ArrowRightLeft className="w-7 h-7 text-indigo-600" /> URL Redirects (301/302)
          </h1>
          <p className="text-sm text-gray-500 mt-1">
            Manage permanent and temporary HTTP redirection rules and track incoming hits.
          </p>
        </div>
        <button
          onClick={() => setIsAddOpen(true)}
          className="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-sm rounded-xl shadow-sm transition-all"
        >
          <Plus className="w-4 h-4" /> Add Redirect
        </button>
      </div>

      <div className="relative w-full max-w-md">
        <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          placeholder="Search redirects by source or target..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full pl-10 pr-4 py-2.5 bg-white border border-gray-200 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
        />
      </div>

      <RedirectTable
        redirects={redirects}
        onToggle={toggleRedirect}
        onDelete={deleteRedirect}
      />

      <AddRedirectDialog
        isOpen={isAddOpen}
        onClose={() => setIsAddOpen(false)}
        onSubmit={createRedirect}
      />
    </div>
  );
};
