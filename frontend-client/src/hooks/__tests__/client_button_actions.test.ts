import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { useClientDomains } from '../useClientDomains';
import { useClientLicenses } from '../useClientLicenses';
import { useDownloads } from '../useDownloads';
import { useKnowledgebase } from '../useKnowledgebase';

describe('Client Interactive Button Actions (TDD)', () => {
  it('Domain Manager buttons: Search availability, Update Nameservers, Toggle Auto-renew', () => {
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
    expect(result.current.domains.length).toBe(1);
    expect(result.current.domains[0].domain_name).toBe('mycompanycloud.com');

    // Reset filter
    act(() => {
      result.current.setSearch('');
    });

    // 3. Toggle Auto Renew Button
    const targetDomain = result.current.domains[0];
    const initialRenew = targetDomain.auto_renew;
    act(() => {
      result.current.toggleAutoRenew(targetDomain.id);
    });
    const updatedDomain = result.current.domains.find((d) => d.id === targetDomain.id);
    expect(updatedDomain?.auto_renew).toBe(!initialRenew);

    // 4. Update Nameservers Button
    act(() => {
      result.current.updateNameservers(targetDomain.id, ['ns3.customdns.org', 'ns4.customdns.org']);
    });
    const nsDomain = result.current.domains.find((d) => d.id === targetDomain.id);
    expect(nsDomain?.nameservers).toEqual(['ns3.customdns.org', 'ns4.customdns.org']);
  });

  it('Client Licenses buttons: Copy Key and Reset IP/Host validation lock', () => {
    const { result } = renderHook(() => useClientLicenses());
    const license = result.current.licenses[0];
    expect(license).toBeDefined();

    // Mock clipboard
    Object.assign(navigator, {
      clipboard: {
        writeText: vi.fn().mockImplementation(() => Promise.resolve()),
      },
    });

    // 1. Copy Key Button
    act(() => {
      result.current.copyKey(license.license_key);
    });
    expect(result.current.copiedKey).toBe(license.license_key);

    // 2. Reset Lock Button
    act(() => {
      result.current.resetLock(license.id);
    });
    const updated = result.current.licenses.find((l) => l.id === license.id);
    expect(updated?.licensed_domain).toBe('');
    expect(updated?.licensed_ip).toBe('');
  });

  it('Downloads buttons: Trigger download link', () => {
    const { result } = renderHook(() => useDownloads());
    const download = result.current.downloads[0];

    // Mock window.open
    const openSpy = vi.spyOn(window, 'open').mockImplementation(() => null);

    act(() => {
      result.current.triggerDownload(download);
    });

    expect(result.current.downloadingId).toBe(download.id);
    openSpy.mockRestore();
  });

  it('Knowledgebase buttons: Search filtering and category navigation', () => {
    const { result } = renderHook(() => useKnowledgebase());

    // 1. Search filter input/button
    act(() => {
      result.current.setSearch('cPanel');
    });
    expect(result.current.articles.length).toBe(1);
    expect(result.current.articles[0].title).toContain('cPanel');

    // 2. Select Category Filter button
    act(() => {
      result.current.setSearch('');
      result.current.setSelectedCategory('Security & SSL');
    });
    expect(result.current.articles.length).toBe(1);
    expect(result.current.articles[0].category).toBe('Security & SSL');

    // 3. Open Article Details button
    act(() => {
      result.current.setActiveArticle(result.current.articles[0]);
    });
    expect(result.current.activeArticle?.slug).toBe('install-ssl-certificate');
  });
});
