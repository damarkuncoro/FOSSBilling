import { describe, it, expect, vi } from 'vitest';
import { domainService } from '../domain.service';
import { authService } from '../auth.service';
import { invoiceService } from '../invoice.service';
import { supportService } from '../support.service';
import { domainRepository } from '../../repositories/domain.repository';
import { authRepository } from '../../repositories/auth.repository';

describe('DomainService', () => {
  it('normalizes domain input correctly', () => {
    expect(domainService.normalizeDomainName('https://MYCOMPANY.COM/')).toBe('mycompany.com');
    expect(domainService.normalizeDomainName('mystartup')).toBe('mystartup.com');
    expect(domainService.normalizeDomainName('toffin.id')).toBe('toffin.id');
  });

  it('delegates checkAvailability with cleaned domain name', async () => {
    vi.spyOn(domainRepository, 'checkAvailability').mockResolvedValue({
      domain: 'toffin.id', tld: 'id', available: false, price: 18.99, currency: 'USD'
    });
    const res = await domainService.checkAvailability('https://TOFFIN.ID/');
    expect(domainRepository.checkAvailability).toHaveBeenCalledWith('toffin.id');
    expect(res.available).toBe(false);
  });

  it('validates nameservers before updating', async () => {
    await expect(domainService.updateNameservers(1, ['invalid'])).rejects.toThrow('Invalid NS');
    vi.spyOn(domainRepository, 'updateNameservers').mockResolvedValue({ message: 'ok' });
    await domainService.updateNameservers(1, ['ns1.cloudflare.com', 'ns2.cloudflare.com']);
    expect(domainRepository.updateNameservers).toHaveBeenCalledWith(1, ['ns1.cloudflare.com', 'ns2.cloudflare.com']);
  });
});

describe('AuthService', () => {
  it('validates password length before changing password', async () => {
    await expect(authService.changePassword('oldpass', 'short')).rejects.toThrow('Short');
  });
});

describe('InvoiceService', () => {
  it('filters invoices by status', () => {
    const invoices: any[] = [{ id: 1, status: 'paid' }, { id: 2, status: 'unpaid' }];
    expect(invoiceService.filterByStatus(invoices, 'unpaid')).toEqual([{ id: 2, status: 'unpaid' }]);
    expect(invoiceService.filterByStatus(invoices, 'all')).toEqual(invoices);
  });
});

describe('SupportService', () => {
  it('validates subject and content before opening ticket', async () => {
    await expect(supportService.openTicket('', 'message')).rejects.toThrow('Subject is required');
    await expect(supportService.openTicket('subject', '')).rejects.toThrow('Message content is required');
  });
});
