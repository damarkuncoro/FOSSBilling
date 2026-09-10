import { domainRepository as repo } from '../repositories/domain.repository';
import type { DomainSearchResult } from '@/types/api';
import type { DomainRecord } from '@/types/clientModules';

export class DomainService {
  normalizeDomainName(i: string) {
    let c = i.toLowerCase().replace(/https?:\/\//, '').replace(/\/$/, '').trim();
    if (c && !c.includes('.')) c += '.com';
    return c;
  }

  checkAvailability = async (i: string): Promise<DomainSearchResult> => {
    const c = this.normalizeDomainName(i); if (!c) throw new Error('Empty');
    return repo.checkAvailability(c);
  };

  listClientDomains = () => repo.listDomains();

  updateNameservers = (id: number, ns: string[]) => {
    const v = ns.map(x => x.trim().toLowerCase()).filter(x => x.length > 0 && x.includes('.'));
    if (v.length === 0) throw new Error('Invalid NS');
    return repo.updateNameservers(id, v);
  };

  toggleAutoRenew = (id: number) => repo.toggleAutoRenew(id);
  getDnsRecords = (id: number) => repo.listDnsRecords(id);
  addDnsRecord = (id: number, r: any) => repo.addDnsRecord(id, r);
  deleteDnsRecord = (id: number, rid: string) => repo.deleteDnsRecord(id, rid);
}

export const domainService = new DomainService();
