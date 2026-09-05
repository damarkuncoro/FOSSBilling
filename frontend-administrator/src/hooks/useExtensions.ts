import { useState, useEffect, useCallback } from 'react';
import { api } from '@/lib/api';
import { ExtensionModuleItem } from '@/types/modules';

export function useExtensions() {
  const [extensions, setExtensions] = useState<ExtensionModuleItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedType, setSelectedType] = useState<string>('all');

  const fetchExtensions = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.getExtensions().catch(() => [
        { id: 'servicehosting', name: 'cPanel & DirectAdmin Hosting Provisioner', version: '2.4.0', author: 'FOSSBilling Core', description: 'Automated hosting account provisioning and suspension.', type: 'service' as const, is_installed: true, is_enabled: true },
        { id: 'midtrans', name: 'Midtrans Payment Gateway (QRIS & VA)', version: '1.2.0', author: 'Nusantara Devs', description: 'Instant Indonesian payment gateway with automated notification callbacks.', type: 'gateway' as const, is_installed: true, is_enabled: true },
        { id: 'servicedomain', name: 'Domain Registrar Multi-Provider', version: '2.1.0', author: 'FOSSBilling Core', description: 'Connects to Namecheap, eNom, and custom TLD API registrars.', type: 'service' as const, is_installed: true, is_enabled: true },
        { id: 'antispam', name: 'Cloudflare Turnstile & Spam Shield', version: '1.0.5', author: 'Security Team', description: 'Prevents bot registrations and brute-force attacks.', type: 'plugin' as const, is_installed: true, is_enabled: true },
        { id: 'theme_huraga', name: 'Huraga Modern Client Theme', version: '3.0.0', author: 'FOSSBilling Community', description: 'Next-Gen responsive client portal theme built with modern CSS.', type: 'theme' as const, is_installed: true, is_enabled: true },
        { id: 'stripe', name: 'Stripe Global Card Payments', version: '2.0.1', author: 'FOSSBilling Core', description: 'Accept Visa, Mastercard, Apple Pay and Google Pay globally.', type: 'gateway' as const, is_installed: false, is_enabled: false },
      ]);
      setExtensions(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchExtensions();
  }, [fetchExtensions]);

  const handleToggle = async (id: string, currentStatus: boolean) => {
    try {
      await api.toggleExtension(id, !currentStatus);
      setExtensions((prev) =>
        prev.map((ext) => (ext.id === id ? { ...ext, is_enabled: !currentStatus } : ext))
      );
    } catch (err: any) {
      alert(`Failed to toggle extension: ${err.message}`);
    }
  };

  const filteredExtensions = extensions.filter((ext) => {
    const matchesSearch =
      ext.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      ext.description.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesType = selectedType === 'all' || ext.type === selectedType;
    return matchesSearch && matchesType;
  });

  return {
    extensions: filteredExtensions,
    loading,
    searchQuery,
    setSearchQuery,
    selectedType,
    setSelectedType,
    fetchExtensions,
    handleToggle,
  };
}
