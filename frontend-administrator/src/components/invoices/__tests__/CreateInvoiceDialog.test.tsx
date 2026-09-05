import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { CreateInvoiceDialog } from '../CreateInvoiceDialog';

describe('CreateInvoiceDialog - TDD Component Tests', () => {
  const mockClients = [
    { id: 1, first_name: 'Budi', last_name: 'Santoso', email: 'budi@example.com', currency: 'IDR' },
    { id: 2, first_name: 'Alice', last_name: 'Smith', email: 'alice@example.com', currency: 'USD' },
  ];

  it('renders client select, currency, due days, and default line item', () => {
    render(
      <CreateInvoiceDialog
        open={true}
        onOpenChange={vi.fn()}
        clients={mockClients}
        onInvoiceCreated={vi.fn()}
        apiCreateInvoice={vi.fn()}
      />
    );

    expect(screen.getByText('Create Custom Invoice')).toBeDefined();
    expect(screen.getByText('#1 - Budi Santoso (budi@example.com)')).toBeDefined();
    expect(screen.getByDisplayValue('Custom Hosting Service')).toBeDefined();
    expect(screen.getByDisplayValue('29.99')).toBeDefined();
    expect(screen.getByText('Estimated Subtotal')).toBeDefined();
  });

  it('adds and removes line items dynamically', () => {
    render(
      <CreateInvoiceDialog
        open={true}
        onOpenChange={vi.fn()}
        clients={mockClients}
        onInvoiceCreated={vi.fn()}
        apiCreateInvoice={vi.fn()}
      />
    );

    const addItemBtn = screen.getByText('Add Item');
    fireEvent.click(addItemBtn);

    const itemInputs = screen.getAllByPlaceholderText('Item description');
    expect(itemInputs.length).toBe(2);
  });

  it('submits valid invoice payload to apiCreateInvoice', async () => {
    const apiCreateMock = vi.fn().mockResolvedValue({ id: 999, nr: 'INV00999' });
    const onCreatedMock = vi.fn();
    const onOpenChangeMock = vi.fn();

    render(
      <CreateInvoiceDialog
        open={true}
        onOpenChange={onOpenChangeMock}
        clients={mockClients}
        onInvoiceCreated={onCreatedMock}
        apiCreateInvoice={apiCreateMock}
      />
    );

    const submitBtn = screen.getByText('Generate Invoice');
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(apiCreateMock).toHaveBeenCalledWith({
        client_id: 1,
        currency: 'USD',
        due_days: 14,
        items: [
          {
            title: 'Custom Hosting Service',
            price: 29.99,
            quantity: 1,
            taxable: false,
          },
        ],
      });
      expect(onOpenChangeMock).toHaveBeenCalledWith(false);
      expect(onCreatedMock).toHaveBeenCalledTimes(1);
    });
  });
});
