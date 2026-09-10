import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import { ClientLayout } from '../ClientLayout';
import * as authLib from '@/lib/auth';
import { CartProvider } from '@/lib/cart';

// Mock all complex sub-components and providers BEFORE anything else
vi.mock('../ClientSidebar', () => ({
  ClientSidebar: () => <div data-testid="client-sidebar">Sidebar</div>,
}));

vi.mock('../../invoices/DepositModal', () => ({
  DepositModal: () => <div data-testid="deposit-modal">Modal</div>,
}));

vi.mock('../../common/LanguageSwitcher', () => ({
  LanguageSwitcher: () => <div data-testid="language-switcher">Lang</div>,
}));

vi.mock('../common/CookieConsentBanner', () => ({
  CookieConsentBanner: () => <div data-testid="cookie-banner">Cookie</div>,
}));

vi.mock('../Footer', () => ({
  Footer: () => <footer data-testid="footer">Footer</footer>,
}));

vi.mock('../NavUserMenu', () => ({
  NavUserMenu: () => <div data-testid="user-menu">UserMenu</div>,
}));

describe('ClientLayout - TDD UX Tests', () => {
  const renderWithProviders = (ui: React.ReactElement) => {
    return render(
      <CartProvider>
        <BrowserRouter>
          {ui}
        </BrowserRouter>
      </CartProvider>
    );
  };

  it('does NOT render sidebar when user is unauthenticated', () => {
    vi.spyOn(authLib, 'useClientAuth').mockReturnValue({
      user: null,
      token: null,
      balance: 0,
      isAuthenticated: false,
      isImpersonated: false,
      isLoading: false,
      login: vi.fn(),
      completeLogin: vi.fn(),
      refreshUser: vi.fn(),
      register: vi.fn(),
      logout: vi.fn(),
      refreshProfile: vi.fn(),
      theme: 'light',
      toggleTheme: vi.fn(),
    });

    renderWithProviders(<ClientLayout />);

    const sidebar = screen.queryByTestId('client-sidebar');
    expect(sidebar).toBeNull();

    // Public links should be visible in navbar
    expect(screen.getByText('Store')).toBeDefined();
  });

  it('renders sidebar when user is authenticated', () => {
    vi.spyOn(authLib, 'useClientAuth').mockReturnValue({
      user: { id: 1, first_name: 'Budi', last_name: 'Santoso', email: 'budi@test.com' },
      token: 'valid-token',
      balance: 100,
      isAuthenticated: true,
      isImpersonated: false,
      isLoading: false,
      login: vi.fn(),
      completeLogin: vi.fn(),
      refreshUser: vi.fn(),
      register: vi.fn(),
      logout: vi.fn(),
      refreshProfile: vi.fn(),
      theme: 'light',
      toggleTheme: vi.fn(),
    });

    renderWithProviders(<ClientLayout />);

    const sidebar = screen.getAllByTestId('client-sidebar')[0];
    expect(sidebar).toBeDefined();

    // Store link should NOT be visible when authenticated based on our logic
    const storeLink = screen.queryByText('Store');
    expect(storeLink).toBeNull();
  });
});
