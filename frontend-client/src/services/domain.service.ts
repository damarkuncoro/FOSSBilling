import { DomainRepository, domainRepository, IDomainRepository } from '../repositories/domain.repository';
import type { DomainSearchResult } from '@/types/api';
import type { DomainRecord } from '@/types/clientModules';

export class DomainService {
  constructor(private repo: IDomainRepository = domainRepository) {}

  /**
   * Normalizes domain name input by stripping protocol, trailing slashes,
   * converting to lowercase, and defaulting to .com if no TLD is present.
   */
  normalizeDomainName(input: string): string {
    let clean = input.toLowerCase().replace(/https?:\/\//, '').replace(/\/$/, '').trim();
    if (clean && !clean.includes('.')) {
      clean += '.com';
    }
    return clean;
  }

  /**
   * Checks real-time registry availability for a domain.
   */
  async checkAvailability(domainInput: string): Promise<DomainSearchResult> {
    const clean = this.normalizeDomainName(domainInput);
    if (!clean) {
      throw new Error('Domain name cannot be empty');
    }
    return this.repo.checkAvailability(clean);
  }

  /**
   * Lists all domains registered by the authenticated client.
   */
  async listClientDomains(): Promise<DomainRecord[]> {
    return this.repo.listDomains();
  }

  /**
   * Updates DNS nameservers for a client domain after validating input.
   */
  async updateNameservers(id: number, nameservers: string[]): Promise<any> {
    const validNS = nameservers
      .map((ns) => ns.trim().toLowerCase())
      .filter((ns) => ns.length > 0 && ns.includes('.'));

    if (validNS.length === 0) {
      throw new Error('At least one valid nameserver (e.g. ns1.example.com) is required');
    }

    return this.repo.updateNameservers(id, validNS);
  }

  /**
   * Toggles the auto-renewal status of a domain.
   */
  async toggleAutoRenew(id: number): Promise<any> {
    return this.repo.toggleAutoRenew(id);
  }
}

export const domainService = new DomainService();
