import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { AddCurrencyDialog } from '../AddCurrencyDialog';

describe('AddCurrencyDialog - TDD Component Tests', () => {
  const mockForm = {
    code: 'EUR',
    title: 'Euro',
    conversion_rate: 0.92,
    format: '€ {{price}}',
  };

  it('renders currency fields when dialog is open', () => {
    render(
      <AddCurrencyDialog
        open={true}
        onOpenChange={vi.fn()}
        form={mockForm}
        setForm={vi.fn()}
        onSave={vi.fn()}
        saving={false}
      />
    );

    expect(screen.getByText('Add New Currency')).toBeDefined();
    expect(screen.getByDisplayValue('EUR')).toBeDefined();
    expect(screen.getByDisplayValue('Euro')).toBeDefined();
    expect(screen.getByDisplayValue('0.92')).toBeDefined();
    expect(screen.getByDisplayValue('€ {{price}}')).toBeDefined();
  });

  it('calls onSave when Save Currency form is submitted', () => {
    const onSaveMock = vi.fn((e) => e.preventDefault());

    render(
      <AddCurrencyDialog
        open={true}
        onOpenChange={vi.fn()}
        form={mockForm}
        setForm={vi.fn()}
        onSave={onSaveMock}
        saving={false}
      />
    );

    const submitBtn = screen.getByText('Save Currency');
    fireEvent.click(submitBtn);

    expect(onSaveMock).toHaveBeenCalledTimes(1);
  });
});
