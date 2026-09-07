import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { OrderSummaryCard } from '../OrderSummaryCard';
import React from 'react';

describe('OrderSummaryCard', () => {
  const defaultProps = {
    subtotal: 100,
    discount: 0,
    promoCode: null,
    tax: 11,
    total: 111,
    checkoutLoading: false,
    onCheckout: vi.fn(),
  };

  it('renders subtotal, tax, and total correctly', () => {
    render(<OrderSummaryCard {...defaultProps} />);

    expect(screen.getByText(/\$100\.00/i)).toBeDefined();
    expect(screen.getByText(/\$11\.00/i)).toBeDefined();
    expect(screen.getByText(/\$111\.00/i)).toBeDefined();
  });

  it('renders discount and promo code when applicable', () => {
    const props = {
      ...defaultProps,
      discount: 20,
      promoCode: 'SAVE20',
      total: 91,
    };
    render(<OrderSummaryCard {...props} />);

    expect(screen.getByText(/Promo Discount \(SAVE20\)/i)).toBeDefined();
    expect(screen.getByText(/-\$20\.00/i)).toBeDefined();
    expect(screen.getByText(/\$91\.00/i)).toBeDefined();
  });

  it('prevents NaN display by using formatMoney utility', () => {
    const props = {
      ...defaultProps,
      subtotal: undefined as any,
      tax: NaN,
      total: null as any,
    };
    render(<OrderSummaryCard {...props} />);

    // Based on our fixed formatMoney, these should show $0.00 instead of NaN
    const zeroDisplays = screen.getAllByText(/\$0\.00/i);
    expect(zeroDisplays.length).toBeGreaterThanOrEqual(3);
  });

  it('shows loading state on checkout button', () => {
    render(<OrderSummaryCard {...defaultProps} checkoutLoading={true} />);

    expect(screen.getByText(/Generating Invoice.../i)).toBeDefined();
    expect(screen.getByRole('button')).toHaveProperty('disabled', true);
  });
});
