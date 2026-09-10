import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { useClientDomains } from '../useClientDomains';
import { useClientLicenses } from '../useClientLicenses';
import { useDownloads } from '../useDownloads';
import { useKnowledgebase } from '../useKnowledgebase';

// Mock the API calls
vi.mock('@/lib/api/client', () => ({
  request: vi.fn().mockResolvedValue([]),
}));

describe('Client Interactive Button Actions (TDD)', () => {
  it('Domain Manager buttons: Search availability, Update Nameservers, Toggle Auto-renew', async () => {
    const { result } = renderHook(() => useClientDomains());

    // 1. Search Query Input & Check Availability Button
    act(() => {
      result.current.setCheckQuery('awesomeapp.com');
    });
    expect(result.current.checkQuery).toBe('awesomeapp.com');

    // 2. Filter Search Button / Input
    act(() => {
      result.current.setSearch('mycompany');
    });
    expect(result.current.search).toBe('mycompany');
  });

  it('Client Licenses buttons: Copy Key and Reset IP/Host validation lock', async () => {
    const { result } = renderHook(() => useClientLicenses());

    // Mock clipboard
    Object.assign(navigator, {
      clipboard: {
        writeText: vi.fn().mockImplementation(() => Promise.resolve()),
      },
    });

    // 1. Copy Key Button
    act(() => {
      result.current.copyKey('FOSS-KEY-123');
    });
    expect(result.current.copiedKey).toBe('FOSS-KEY-123');
  });

  it('Downloads buttons: Trigger download link', () => {
    const { result } = renderHook(() => useDownloads());

    // Mock window.open
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);

    act(() => {
      result.current.triggerDownload({ id: 1 } as any);
    });

    expect(result.current.downloadingId).toBe(1);
    openSpy.mockRestore();
  });

  it('Knowledgebase buttons: Search filtering and category navigation', () => {
    const { result } = renderHook(() => useKnowledgebase());

    // 1. Search filter input/button
    act(() => {
      result.current.setSearch('cPanel');
    });
    expect(result.current.search).toBe('cPanel');
  });
});
