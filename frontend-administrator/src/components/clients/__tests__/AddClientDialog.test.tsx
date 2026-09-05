import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { AddClientDialog } from '../AddClientDialog';

describe('AddClientDialog - TDD Component Tests', () => {
  it('renders dialog form when open is true', () => {
    render(
      <AddClientDialog
        open={true}
        onOpenChange={vi.fn()}
        onSubmit={vi.fn()}
      />
    );

    expect(screen.getByText('Add New Client')).toBeDefined();
    expect(screen.getByPlaceholderText('e.g. Budi')).toBeDefined();
    expect(screen.getByPlaceholderText('e.g. Santoso')).toBeDefined();
    expect(screen.getByPlaceholderText('client@example.com')).toBeDefined();
  });

  it('submits valid form data and resets form', async () => {
    const onSubmitMock = vi.fn().mockResolvedValue(true);
    const onOpenChangeMock = vi.fn();

    render(
      <AddClientDialog
        open={true}
        onOpenChange={onOpenChangeMock}
        onSubmit={onSubmitMock}
      />
    );

    fireEvent.change(screen.getByPlaceholderText('e.g. Budi'), { target: { value: 'Jane' } });
    fireEvent.change(screen.getByPlaceholderText('e.g. Santoso'), { target: { value: 'Doe' } });
    fireEvent.change(screen.getByPlaceholderText('client@example.com'), { target: { value: 'jane@example.com' } });

    fireEvent.click(screen.getByText('Save Client'));

    await waitFor(() => {
      expect(onSubmitMock).toHaveBeenCalledWith(
        expect.objectContaining({
          first_name: 'Jane',
          last_name: 'Doe',
          email: 'jane@example.com',
          status: 'active',
        })
      );
      expect(onOpenChangeMock).toHaveBeenCalledWith(false);
    });
  });
});
