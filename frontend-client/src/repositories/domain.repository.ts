import { request } from '../lib/api/client';
import type { DomainSearchResult } from '@/types/api';
import type { DomainRecord } from '@/types/clientModules';

export class DomainRepository {
  checkAvailability = (d: string) => request<DomainSearchResult>(`/guest/domains/check?domain=${encodeURIComponent(d)}`);
  listDomains = () => request<DomainRecord[]>('/client/domains');
  updateNameservers = (id: number, ns: string[]) => request(`/client/domains/${id}/nameservers`, { method: 'PUT', body: JSON.stringify({ nameservers: ns }) });
  toggleAutoRenew = (id: number) => request(`/client/domains/${id}/toggle-autorenew`, { method: 'POST' });
  listDnsRecords = (id: number) => request<any[]>(`/client/domains/${id}/dns`);
  addDnsRecord = (id: number, r: any) => request(`/client/domains/${id}/dns`, { method: 'POST', body: JSON.stringify(r) });
  deleteDnsRecord = (id: number, rid: string) => request(`/client/domains/${id}/dns/${rid}`, { method: 'DELETE' });
}

export const domainRepository = new DomainRepository();
