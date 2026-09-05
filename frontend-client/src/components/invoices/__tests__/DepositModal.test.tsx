import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { DepositModal } from '../DepositModal';
import { api } from '@/lib/api';

describe('DepositModal - TDD Component Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders quick amounts and allows selecting preset amount', () => {
    render(<DepositModal open={true} onOpenChange={vi.fn()} currency="USD" />);

    expect(screen.getByText('Add Funds to Account Balance')).toBeDefined();
    expect(screen.getByText('$10')).toBeDefined();
    expect(screen.getByText('$25')).toBeDefined();
    expect(screen.getByText('$50')).toBeDefined();
    expect(screen.getByText('$100')).toBeDefined();
    expect(screen.getByText('$250')).toBeDefined();

    const btn100 = screen.getByText('$100');
    fireEvent.click(btn100);

    const input = screen.getByDisplayValue('100');
    expect(input).toBeDefined();
  });

  it('submits deposit request to API and displays invoice confirmation', async () => {
    const depositSpy = vi.spyOn(api, 'depositFunds').mockResolvedValue({
      invoice_id: 888,
      nr: 'INV00888',
      total: 100,
      currency: 'USD',
      status: 'unpaid',
      message: 'Deposit invoice generated successfully',
    });

    const successCallback = vi.fn();
    render(
      <DepositModal
        open={true}
        onOpenChange={vi.fn()}
        currency="USD"
        onDepositSuccess={successCallback}
      />
    );

    const submitBtn = screen.getByText('Generate Top-Up Invoice');
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(depositSpy).toHaveBeenCalledWith(50, 'USD');
      expect(screen.getByText('Deposit Invoice Generated!')).toBeDefined();
      expect(screen.getByText('#INV00888')).toBeDefined();
      expect(successCallback).toHaveBeenCalledTimes(1);
    });
  });
});
