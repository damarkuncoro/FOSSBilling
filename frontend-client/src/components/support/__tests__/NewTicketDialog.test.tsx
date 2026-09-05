import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { NewTicketDialog } from '../NewTicketDialog';

describe('NewTicketDialog - TDD Component Tests', () => {
  const mockForm = {
    subject: 'Cannot connect to cPanel',
    priority: 'high',
    message: 'Port 2083 is returning connection timed out.',
  };

  it('renders form inputs and priority options when open', () => {
    render(
      <NewTicketDialog
        open={true}
        onOpenChange={vi.fn()}
        form={mockForm}
        onChange={vi.fn()}
        onSubmit={vi.fn()}
      />
    );

    expect(screen.getByText('Open Support Inquiry')).toBeDefined();
    expect(screen.getByDisplayValue('Cannot connect to cPanel')).toBeDefined();
    expect(screen.getByDisplayValue('Port 2083 is returning connection timed out.')).toBeDefined();
    expect(screen.getByText('High - Production Degraded')).toBeDefined();
  });

  it('triggers onSubmit when Submit Ticket button is clicked', () => {
    const onSubmitMock = vi.fn((e) => e.preventDefault());

    render(
      <NewTicketDialog
        open={true}
        onOpenChange={vi.fn()}
        form={mockForm}
        onChange={vi.fn()}
        onSubmit={onSubmitMock}
      />
    );

    const submitBtn = screen.getByText('Submit Ticket');
    fireEvent.click(submitBtn);

    expect(onSubmitMock).toHaveBeenCalledTimes(1);
  });
});
