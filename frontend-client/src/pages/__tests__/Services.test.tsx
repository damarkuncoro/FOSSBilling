import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Services } from '../Services';
import * as servicesHook from '@/hooks/useClientServices';

// Mock the hook
vi.mock('@/hooks/useClientServices');

describe('Services Page - TDD Component Tests', () => {
  it('renders loading state correctly', () => {
    vi.mocked(servicesHook.useClientServices).mockReturnValue({
      orders: [],
      loading: true,
      downloadLink: null,
      downloadModal: false,
      setDownloadModal: vi.fn(),
      fetchServices: vi.fn(),
      handleGetDownload: vi.fn(),
      handleSyncStatus: vi.fn(),
      handleChangePassword: vi.fn(),
    });

    render(<Services />);

    expect(screen.getByText(/Loading your cloud services/i)).toBeDefined();
  });

  it('renders empty state when no orders are found', () => {
    vi.mocked(servicesHook.useClientServices).mockReturnValue({
      orders: [],
      loading: false,
      downloadLink: null,
      downloadModal: false,
      setDownloadModal: vi.fn(),
      fetchServices: vi.fn(),
      handleGetDownload: vi.fn(),
      handleSyncStatus: vi.fn(),
      handleChangePassword: vi.fn(),
    });

    render(<Services />);

    expect(screen.getByText(/You don't have any active subscriptions yet/i)).toBeDefined();
  });

  it('renders list of services correctly', () => {
    const mockOrders = [
      {
        id: 101,
        title: 'Premium Hosting',
        status: 'active',
        period: '1M',
        price: 150000,
        currency: 'USD',
        product_id: 1,
        config: {}
      },
      {
        id: 102,
        title: 'Cloud VPS',
        status: 'suspended',
        period: '1Y',
        price: 2500000,
        currency: 'USD',
        product_id: 2,
        config: {}
      }
    ];

    vi.mocked(servicesHook.useClientServices).mockReturnValue({
      orders: mockOrders as any,
      loading: false,
      downloadLink: null,
      downloadModal: false,
      setDownloadModal: vi.fn(),
      fetchServices: vi.fn(),
      handleGetDownload: vi.fn(),
      handleSyncStatus: vi.fn(),
      handleChangePassword: vi.fn(),
    });

    render(<Services />);

    expect(screen.getByText('Premium Hosting')).toBeDefined();
    expect(screen.getByText('Cloud VPS')).toBeDefined();
    expect(screen.getByText('ACTIVE')).toBeDefined();
    expect(screen.getByText('SUSPENDED')).toBeDefined();
  });

  it('calls fetchServices when Refresh button is clicked', () => {
    const fetchServicesMock = vi.fn();
    vi.mocked(servicesHook.useClientServices).mockReturnValue({
      orders: [],
      loading: false,
      downloadLink: null,
      downloadModal: false,
      setDownloadModal: vi.fn(),
      fetchServices: fetchServicesMock,
      handleGetDownload: vi.fn(),
      handleSyncStatus: vi.fn(),
      handleChangePassword: vi.fn(),
    });

    render(<Services />);

    const refreshBtn = screen.getByRole('button', { name: /refresh/i });
    fireEvent.click(refreshBtn);

    expect(fetchServicesMock).toHaveBeenCalledTimes(1);
  });

  it('calls handleSyncStatus when Sync button is clicked', () => {
    const handleSyncStatusMock = vi.fn();
    vi.mocked(servicesHook.useClientServices).mockReturnValue({
      orders: [{ id: 101, title: 'Service 1', status: 'active', price: 100, currency: 'USD' }] as any,
      loading: false,
      downloadLink: null,
      downloadModal: false,
      setDownloadModal: vi.fn(),
      fetchServices: vi.fn(),
      handleGetDownload: vi.fn(),
      handleSyncStatus: handleSyncStatusMock,
      handleChangePassword: vi.fn(),
    });

    render(<Services />);

    const syncBtn = screen.getByTitle('Sync Status');
    fireEvent.click(syncBtn);

    expect(handleSyncStatusMock).toHaveBeenCalledWith(101);
  });

  it('calls handleGetDownload when Download button is clicked', () => {
    const handleGetDownloadMock = vi.fn();
    vi.mocked(servicesHook.useClientServices).mockReturnValue({
      orders: [{ id: 101, title: 'Service 1', status: 'active', price: 100, currency: 'USD' }] as any,
      loading: false,
      downloadLink: null,
      downloadModal: false,
      setDownloadModal: vi.fn(),
      fetchServices: vi.fn(),
      handleGetDownload: handleGetDownloadMock,
      handleSyncStatus: vi.fn(),
      handleChangePassword: vi.fn(),
    });

    render(<Services />);

    const downloadBtn = screen.getByText(/Get Download \/ License/i);
    fireEvent.click(downloadBtn);

    expect(handleGetDownloadMock).toHaveBeenCalledWith(101);
  });

  it('opens password modal when Key icon is clicked and submits change', async () => {
    const handleChangePasswordMock = vi.fn();
    vi.mocked(servicesHook.useClientServices).mockReturnValue({
      orders: [{
        id: 101,
        title: 'Hosting',
        status: 'active',
        price: 100,
        currency: 'USD',
        config: { account_details: JSON.stringify({ username: 'budi', password: 'old-password' }) }
      }] as any,
      loading: false,
      downloadLink: null,
      downloadModal: false,
      setDownloadModal: vi.fn(),
      fetchServices: vi.fn(),
      handleGetDownload: vi.fn(),
      handleSyncStatus: vi.fn(),
      handleChangePassword: handleChangePasswordMock,
    });

    render(<Services />);

    // 1. Click the small Key icon button next to password
    const passwordTrigger = screen.getByTitle('Change Password');
    fireEvent.click(passwordTrigger);

    // 2. Look for modal title (using findBy to wait for portal)
    expect(await screen.findByText(/Change Service Password/i)).toBeDefined();

    // 3. Type new password
    const passwordInput = screen.getByLabelText(/New Password/i);
    fireEvent.change(passwordInput, { target: { value: 'new-secret-123' } });

    // 4. Click confirmation button in modal
    const confirmBtn = screen.getByRole('button', { name: /Change Password/i });
    fireEvent.click(confirmBtn);

    expect(handleChangePasswordMock).toHaveBeenCalledWith(101, 'new-secret-123');
  });
});
