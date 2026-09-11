import React from 'react';
import { renderHook } from '@testing-library/react-hooks';
import { useAdminNotifications } from '../useAdminNotifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import * as api from '../../lib/api/client';

// Mocking the request function
jest.mock('../../lib/api/client', () => ({
  request: jest.fn(),
}));

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
);

describe('Notifications Stability Regression', () => {
  it('should handle null response from API without crashing', async () => {
    // Skenario: API mengembalikan null (seperti saat error 404/500 terjadi)
    (api.request as jest.Mock).mockResolvedValue(null);

    const { result, waitForNextUpdate } = renderHook(() => useAdminNotifications(), { wrapper });

    await waitForNextUpdate();

    // Verifikasi: notifications harus menjadi array kosong [], bukan null
    expect(result.current.notifications).toBeInstanceOf(Array);
    expect(result.current.notifications).toHaveLength(0);
    expect(result.current.unreadCount).toBe(0);
  });
});
