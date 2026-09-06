import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import React from 'react';
import { AddCouponDialog } from '../AddCouponDialog';

describe('AddCouponDialog - Component Tests', () => {
  it('renders trigger button and modal form inputs when open', () => {
    const onOpenChange = vi.fn();
    const setForm = vi.fn();
    const onGenerateCode = vi.fn();
    const onSave = vi.fn((e) => e.preventDefault());

    render(
      <AddCouponDialog
        open={true}
        onOpenChange={onOpenChange}
        form={{
          code: 'PROMO2026',
          type: 'percentage',
          value: 20,
          max_uses: 50,
          is_active: true,
        }}
        setForm={setForm}
        onGenerateCode={onGenerateCode}
        onSave={onSave}
        saving={false}
      />
    );

    expect(screen.getByText('Create Promo Coupon')).toBeDefined();
    expect(screen.getByDisplayValue('PROMO2026')).toBeDefined();
    expect(screen.getByText('Generate Code')).toBeDefined();

    const generateBtn = screen.getByText('Generate Code');
    fireEvent.click(generateBtn);
    expect(onGenerateCode).toHaveBeenCalled();
  });

  it('triggers form submit on clicking submit button', () => {
    const onOpenChange = vi.fn();
    const setForm = vi.fn();
    const onGenerateCode = vi.fn();
    const onSave = vi.fn((e) => e.preventDefault());

    render(
      <AddCouponDialog
        open={true}
        onOpenChange={onOpenChange}
        form={{
          code: 'SALE50',
          type: 'fixed',
          value: 50,
          max_uses: 10,
          is_active: true,
        }}
        setForm={setForm}
        onGenerateCode={onGenerateCode}
        onSave={onSave}
        saving={false}
      />
    );

    const submitBtn = screen.getByText('Create Coupon');
    fireEvent.click(submitBtn);
    expect(onSave).toHaveBeenCalled();
  });
});
