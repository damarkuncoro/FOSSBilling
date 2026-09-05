import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { InvoicesTable } from '../InvoicesTable';
import type { Invoice, ClientProfile } from '@/types/api';

describe('InvoicesTable - Download Invoice PDF (TDD)', () => {
  const mockUser: ClientProfile = {
    id: 1,
    first_name: 'John',
    last_name: 'Doe',
    email: 'client@fossbilling.org',
    currency: 'USD',
    created_at: '2026-01-01T00:00:00Z',
  };

  const mockInvoices: Invoice[] = [
    {
      id: 101,
      serie_nr: 'INV-2026-0101',
      client_id: 1,
      status: 'paid',
      subtotal: 100,
      discount: 0,
      tax: 11,
      total: 111,
      currency: 'USD',
      created_at: '2026-09-01T10:00:00Z',
      due_at: '2026-09-15T10:00:00Z',
    },
    {
      id: 102,
      serie_nr: 'INV-2026-0102',
      client_id: 1,
      status: 'unpaid',
      subtotal: 50,
      discount: 0,
      tax: 5.5,
      total: 55.5,
      currency: 'USD',
      created_at: '2026-09-02T10:00:00Z',
      due_at: '2026-09-16T10:00:00Z',
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.setItem('auth_token', 'mock-jwt-token');
  });

  it('renders invoice list with status badges and PDF download buttons', () => {
    render(
      <InvoicesTable
        invoices={mockInvoices}
        loading={false}
        user={mockUser}
        onPayModal={vi.fn()}
      />
    );

    expect(screen.getByText('#INV-2026-0101')).toBeDefined();
    expect(screen.getByText('#INV-2026-0102')).toBeDefined();
    expect(screen.getByText('PAID')).toBeDefined();
    expect(screen.getByText('UNPAID')).toBeDefined();
    expect(screen.getAllByText('PDF').length).toBe(2);
  });

  it('triggers invoice PDF download on clicking the PDF button', async () => {
    const mockHtml = '<html><body><h1>FOSSBilling Invoice #INV-2026-0101</h1></body></html>';
    const blob = new Blob([mockHtml], { type: 'text/html' });

    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      blob: () => Promise.resolve(blob),
    } as unknown as Response);

    const mockCreateObjectURL = vi.fn(() => 'blob:mock-url');
    const mockRevokeObjectURL = vi.fn();
    globalThis.URL.createObjectURL = mockCreateObjectURL;
    globalThis.URL.revokeObjectURL = mockRevokeObjectURL;

    render(
      <InvoicesTable
        invoices={mockInvoices}
        loading={false}
        user={mockUser}
        onPayModal={vi.fn()}
      />
    );

    const pdfButtons = screen.getAllByText('PDF');
    fireEvent.click(pdfButtons[0]);

    await waitFor(() => {
      expect(fetchSpy).toHaveBeenCalledWith('/api/v1/client/invoices/101/pdf', {
        headers: { Authorization: 'Bearer mock-jwt-token' },
      });
      expect(mockCreateObjectURL).toHaveBeenCalled();
    });
  });

  it('triggers onPayModal when Pay Now button is clicked for unpaid invoice', () => {
    const onPayModalMock = vi.fn();

    render(
      <InvoicesTable
        invoices={mockInvoices}
        loading={false}
        user={mockUser}
        onPayModal={onPayModalMock}
      />
    );

    const payButton = screen.getByText('Pay Now');
    fireEvent.click(payButton);

    expect(onPayModalMock).toHaveBeenCalledWith(mockInvoices[1]);
  });
});
