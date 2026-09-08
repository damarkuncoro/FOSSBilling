import { request } from '../lib/api/client';
import type { DomainSearchResult } from '@/types/api';
import type { DomainRecord } from '@/types/clientModules';

export interface IDomainRepository {
  checkAvailability(domain: string): Promise<DomainSearchResult>;
  listDomains(): Promise<DomainRecord[]>;
  updateNameservers(id: number, nameservers: string[]): Promise<any>;
  toggleAutoRenew(id: number): Promise<any>;
  listDnsRecords(id: number): Promise<any[]>;
  addDnsRecord(id: number, record: any): Promise<any>;
  deleteDnsRecord(id: number, recordId: string): Promise<any>;
}

export class DomainRepository implements IDomainRepository {
  async checkAvailability(domain: string): Promise<DomainSearchResult> {
    return request<DomainSearchResult>(`/guest/domains/check?domain=${encodeURIComponent(domain)}`);
  }

  async listDomains(): Promise<DomainRecord[]> {
    return request<DomainRecord[]>('/client/domains');
  }

  async updateNameservers(id: number, nameservers: string[]): Promise<any> {
    return request(`/client/domains/${id}/nameservers`, {
      method: 'PUT',
      body: JSON.stringify({ nameservers }),
    });
  }

  async toggleAutoRenew(id: number): Promise<any> {
    return request(`/client/domains/${id}/toggle-autorenew`, {
      method: 'POST',
    });
  }

  async listDnsRecords(id: number): Promise<any[]> {
    return request<any[]>(`/client/domains/${id}/dns`);
  }

  async addDnsRecord(id: number, record: any): Promise<any> {
    return request(`/client/domains/${id}/dns`, {
      method: 'POST',
      body: JSON.stringify(record),
    });
  }

  async deleteDnsRecord(id: number, recordId: string): Promise<any> {
    return request(`/client/domains/${id}/dns/${recordId}`, {
      method: 'DELETE',
    });
  }
}

export const domainRepository = new DomainRepository();
