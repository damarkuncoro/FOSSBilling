import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { DashboardStatsGrid } from '../DashboardStatsGrid';
import { ClientProfile, Order, Invoice, SupportTicket } from '@/types/api';

describe('DashboardStatsGrid - Component Tests', () => {
  const mockUser: ClientProfile = {
    id: 1,
    email: 'client@example.com',
    first_name: 'Budi',
    last_name: 'Santoso',
    currency: 'USD',
  };

  const mockOrders: Order[] = [
    { id: 1, client_id: 1, product_id: 101, title: 'Cloud VPS', status: 'active', price: 15, period: '1M', created_at: '2026-01-01' },
    { id: 2, client_id: 1, product_id: 102, title: 'DirectAdmin', status: 'active', price: 10, period: '1M', created_at: '2026-01-01' },
  ];

  const mockInvoices: Invoice[] = [
    { id: 1, client_id: 1, serie_nr: 'INV001', status: 'unpaid', subtotal: 25, discount: 0, tax: 0, total: 25, currency: 'USD', due_at: '2026-02-01', created_at: '2026-01-01' },
  ];

  const mockTickets: SupportTicket[] = [
    { id: 1, client_id: 1, subject: 'Need Help', status: 'open', priority: 'high', created_at: '2026-01-01', updated_at: '2026-01-01' },
  ];

  it('renders credit balance, active services, unpaid invoices, and open tickets correctly', () => {
    render(
      <DashboardStatsGrid
        user={mockUser}
        balance={250.75}
        orders={mockOrders}
        invoices={mockInvoices}
        tickets={mockTickets}
        unpaidInvoices={mockInvoices}
        activeOrders={mockOrders}
      />
    );

    expect(screen.getByText('Account Credit Balance')).toBeDefined();
    expect(screen.getByText('$250.75')).toBeDefined();
    expect(screen.getByText('Active Services')).toBeDefined();
    expect(screen.getByText('2')).toBeDefined();
    expect(screen.getByText('Total Invoices')).toBeDefined();
    expect(screen.getByText('Support Tickets')).toBeDefined();
    expect(screen.getByText('1 unpaid')).toBeDefined();
  });
});
