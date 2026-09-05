import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ChangePasswordCard } from '../ChangePasswordCard';
import { api } from '@/lib/api';

describe('ChangePasswordCard - TDD Component Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders password input fields', () => {
    render(<ChangePasswordCard />);

    expect(screen.getByText('Security & Password')).toBeDefined();
    expect(screen.getByPlaceholderText('Enter current password...')).toBeDefined();
    expect(screen.getByPlaceholderText('At least 6 characters...')).toBeDefined();
    expect(screen.getByPlaceholderText('Re-type new password...')).toBeDefined();
  });

  it('shows error if new passwords do not match', async () => {
    render(<ChangePasswordCard />);

    fireEvent.change(screen.getByPlaceholderText('Enter current password...'), {
      target: { value: 'CurrentPass123!' },
    });
    fireEvent.change(screen.getByPlaceholderText('At least 6 characters...'), {
      target: { value: 'NewPassword123!' },
    });
    fireEvent.change(screen.getByPlaceholderText('Re-type new password...'), {
      target: { value: 'DifferentPassword123!' },
    });

    fireEvent.click(screen.getByText('Update Password'));

    await waitFor(() => {
      expect(screen.getByText('New passwords do not match.')).toBeDefined();
    });
  });

  it('submits valid password change and shows success notification', async () => {
    const changePassSpy = vi.spyOn(api, 'changePassword').mockResolvedValue({
      success: true,
      message: 'Password changed successfully',
    });

    render(<ChangePasswordCard />);

    fireEvent.change(screen.getByPlaceholderText('Enter current password...'), {
      target: { value: 'OldPassword123!' },
    });
    fireEvent.change(screen.getByPlaceholderText('At least 6 characters...'), {
      target: { value: 'NewStrongPassword123!' },
    });
    fireEvent.change(screen.getByPlaceholderText('Re-type new password...'), {
      target: { value: 'NewStrongPassword123!' },
    });

    fireEvent.click(screen.getByText('Update Password'));

    await waitFor(() => {
      expect(changePassSpy).toHaveBeenCalledWith({
        current_password: 'OldPassword123!',
        new_password: 'NewStrongPassword123!',
      });
      expect(screen.getByText('Your password has been changed successfully!')).toBeDefined();
    });
  });
});
