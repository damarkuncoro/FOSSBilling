import { describe, it, expect, vi } from 'vitest';
import { DomainService } from '../domain.service';
import { AuthService } from '../auth.service';
import { InvoiceService } from '../invoice.service';
import { OrderService } from '../order.service';
import { LicenseService } from '../license.service';
import { SupportService } from '../support.service';
import { IDomainRepository } from '../../repositories/domain.repository';
import { IAuthRepository } from '../../repositories/auth.repository';
import { IInvoiceRepository } from '../../repositories/invoice.repository';
import { IOrderRepository } from '../../repositories/order.repository';
import { ILicenseRepository } from '../../repositories/license.repository';
import { ISupportRepository } from '../../repositories/support.repository';

describe('DomainService', () => {
  it('normalizes domain input correctly', () => {
    const mockRepo: IDomainRepository = {
      checkAvailability: vi.fn(),
      listDomains: vi.fn(),
      updateNameservers: vi.fn(),
      toggleAutoRenew: vi.fn(),
    };
    const service = new DomainService(mockRepo);

    expect(service.normalizeDomainName('https://MYCOMPANY.COM/')).toBe('mycompany.com');
    expect(service.normalizeDomainName('mystartup')).toBe('mystartup.com');
    expect(service.normalizeDomainName('toffin.id')).toBe('toffin.id');
  });

  it('delegates checkAvailability with cleaned domain name', async () => {
    const mockRepo: IDomainRepository = {
      checkAvailability: vi.fn().mockResolvedValue({
        domain: 'toffin.id',
        tld: 'id',
        available: false,
        price: 18.99,
        currency: 'USD',
      }),
      listDomains: vi.fn(),
      updateNameservers: vi.fn(),
      toggleAutoRenew: vi.fn(),
    };
    const service = new DomainService(mockRepo);

    const res = await service.checkAvailability('https://TOFFIN.ID/');
    expect(mockRepo.checkAvailability).toHaveBeenCalledWith('toffin.id');
    expect(res.available).toBe(false);
    expect(res.tld).toBe('id');
  });

  it('validates nameservers before updating', async () => {
    const mockRepo: IDomainRepository = {
      checkAvailability: vi.fn(),
      listDomains: vi.fn(),
      updateNameservers: vi.fn().mockResolvedValue({ message: 'ok' }),
      toggleAutoRenew: vi.fn(),
    };
    const service = new DomainService(mockRepo);

    await expect(service.updateNameservers(1, ['invalid_no_dot'])).rejects.toThrow(
      'At least one valid nameserver'
    );

    await service.updateNameservers(1, ['ns1.cloudflare.com', 'ns2.cloudflare.com']);
    expect(mockRepo.updateNameservers).toHaveBeenCalledWith(1, ['ns1.cloudflare.com', 'ns2.cloudflare.com']);
  });
});

describe('AuthService', () => {
  it('validates password length before changing password', async () => {
    const mockRepo: IAuthRepository = {
      login: vi.fn(),
      register: vi.fn(),
      getProfile: vi.fn(),
      updateProfile: vi.fn(),
      changePassword: vi.fn(),
    };
    const service = new AuthService(mockRepo);

    await expect(service.changePassword('oldpass', 'short')).rejects.toThrow('at least 8 characters');
  });
});

describe('InvoiceService', () => {
  it('filters invoices by status', () => {
    const mockRepo: IInvoiceRepository = {
      listInvoices: vi.fn(),
      getInvoice: vi.fn(),
      payWithBalance: vi.fn(),
      depositFunds: vi.fn(),
    };
    const service = new InvoiceService(mockRepo);

    const invoices: any[] = [
      { id: 1, status: 'paid' },
      { id: 2, status: 'unpaid' },
    ];

    expect(service.filterByStatus(invoices, 'unpaid')).toEqual([{ id: 2, status: 'unpaid' }]);
    expect(service.filterByStatus(invoices, 'all')).toEqual(invoices);
  });
});

describe('SupportService', () => {
  it('validates subject and content before opening ticket', async () => {
    const mockRepo: ISupportRepository = {
      listTickets: vi.fn(),
      getTicket: vi.fn(),
      openTicket: vi.fn(),
      replyTicket: vi.fn(),
      closeTicket: vi.fn(),
    };
    const service = new SupportService(mockRepo);

    await expect(service.openTicket('', 'message')).rejects.toThrow('Subject is required');
    await expect(service.openTicket('subject', '')).rejects.toThrow('Message content is required');
  });
});
