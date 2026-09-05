import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import { ClientSidebar } from '../ClientSidebar';
import * as authLib from '@/lib/auth';

describe('ClientSidebar - TDD Component Tests', () => {
  it('renders brand and navigation links for authenticated user', () => {
    vi.spyOn(authLib, 'useClientAuth').mockReturnValue({
      user: {
        id: 1,
        first_name: 'Budi',
        last_name: 'Santoso',
        email: 'budi@example.com',
        currency: 'USD',
        created_at: '2026-01-01',
      },
      balance: 150,
      isAuthenticated: true,
      login: vi.fn(),
      register: vi.fn(),
      logout: vi.fn(),
      theme: 'light',
      toggleTheme: vi.fn(),
      loading: false,
    });

    render(
      <BrowserRouter>
        <ClientSidebar onOpenDeposit={vi.fn()} onItemClick={vi.fn()} />
      </BrowserRouter>
    );

    expect(screen.getByText('FOSSBilling')).toBeDefined();
    expect(screen.getByText('Client Portal')).toBeDefined();
    expect(screen.getByText('Dashboard')).toBeDefined();
    expect(screen.getByText('My Services')).toBeDefined();
    expect(screen.getByText('Invoices & Funds')).toBeDefined();
    expect(screen.getByText('Account Balance')).toBeDefined();
    expect(screen.getByText('$150.00')).toBeDefined();
    expect(screen.getByText('Budi Santoso')).toBeDefined();
  });

  it('triggers onOpenDeposit when Add balance button is clicked', () => {
    const onOpenDepositMock = vi.fn();
    vi.spyOn(authLib, 'useClientAuth').mockReturnValue({
      user: { id: 1, first_name: 'Budi', last_name: 'Santoso', email: 'budi@example.com', created_at: '' },
      balance: 50,
      isAuthenticated: true,
      login: vi.fn(),
      register: vi.fn(),
      logout: vi.fn(),
      theme: 'light',
      toggleTheme: vi.fn(),
      loading: false,
    });

    render(
      <BrowserRouter>
        <ClientSidebar onOpenDeposit={onOpenDepositMock} />
      </BrowserRouter>
    );

    const addBtn = screen.getByText('Add');
    fireEvent.click(addBtn);
    expect(onOpenDepositMock).toHaveBeenCalledTimes(1);
  });

  it('renders Login and Register actions when user is unauthenticated', () => {
    vi.spyOn(authLib, 'useClientAuth').mockReturnValue({
      user: null,
      balance: 0,
      isAuthenticated: false,
      login: vi.fn(),
      register: vi.fn(),
      logout: vi.fn(),
      theme: 'light',
      toggleTheme: vi.fn(),
      loading: false,
    });

    render(
      <BrowserRouter>
        <ClientSidebar />
      </BrowserRouter>
    );

    expect(screen.getByText('Login')).toBeDefined();
    expect(screen.getByText('Register')).toBeDefined();
  });
});
