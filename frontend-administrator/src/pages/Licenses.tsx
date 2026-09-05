import React from 'react';
import { Key, Plus, Search } from 'lucide-react';
import { useLicenses } from '../hooks/useLicenses';
import { LicenseTable } from '../components/licenses/LicenseTable';
import { AddLicenseDialog } from '../components/licenses/AddLicenseDialog';
import { LicenseDetailsModal } from '../components/licenses/LicenseDetailsModal';

export const Licenses: React.FC = () => {
  const {
    licenses,
    logs,
    search,
    setSearch,
    statusFilter,
    setStatusFilter,
    isAddOpen,
    setIsAddOpen,
    selectedLicense,
    setSelectedLicense,
    createLicense,
    toggleStatus,
    reIssueKey,
    resetLock,
    generateKey,
  } = useLicenses();

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
            <Key className="w-7 h-7 text-indigo-600" /> Software Licenses
          </h1>
          <p className="text-sm text-gray-500 mt-1">
            Generate and manage software license keys, domain/IP locks, and validation logs.
          </p>
        </div>
        <button
          onClick={() => setIsAddOpen(true)}
          className="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-sm rounded-xl shadow-sm transition-all"
        >
          <Plus className="w-4 h-4" /> Issue License Key
        </button>
      </div>

      <div className="flex flex-col sm:flex-row items-center gap-3">
        <div className="relative flex-1 w-full">
          <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Search by license key, client name, or domain..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-white border border-gray-200 rounded-xl text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
          />
        </div>
        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value)}
          className="px-3.5 py-2.5 bg-white border border-gray-200 rounded-xl text-sm text-gray-700 focus:border-indigo-500 outline-none w-full sm:w-auto"
        >
          <option value="all">All Statuses</option>
          <option value="active">Active</option>
          <option value="suspended">Suspended</option>
          <option value="expired">Expired</option>
        </select>
      </div>

      <LicenseTable
        licenses={licenses}
        onSelect={setSelectedLicense}
        onToggleStatus={toggleStatus}
        onReissue={reIssueKey}
        onResetLock={resetLock}
      />

      <AddLicenseDialog
        isOpen={isAddOpen}
        onClose={() => setIsAddOpen(false)}
        onSubmit={createLicense}
        onGenerateKey={generateKey}
      />

      <LicenseDetailsModal
        license={selectedLicense}
        logs={logs}
        onClose={() => setSelectedLicense(null)}
      />
    </div>
  );
};
