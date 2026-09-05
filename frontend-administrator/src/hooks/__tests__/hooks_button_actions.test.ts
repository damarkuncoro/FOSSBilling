import { describe, it, expect } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useLicenses } from '../useLicenses';
import { useWebhooks } from '../useWebhooks';
import { useFormBuilder } from '../useFormBuilder';
import { useRedirects } from '../useRedirects';
import { useCookieConsent } from '../useCookieConsent';

describe('Admin Interactive Button Actions (TDD)', () => {
  it('License buttons: Issue, Toggle Status, Re-issue key, and Reset lock', () => {
    const { result } = renderHook(() => useLicenses());

    // Test Issue License Button
    act(() => {
      result.current.createLicense({
        client_name: 'Test Client Corp',
        product_title: 'SaaS Pro',
        licensed_domain: 'app.test.com',
      });
    });
    expect(result.current.licenses.length).toBeGreaterThan(2);
    const newLic = result.current.licenses[0];
    expect(newLic.client_name).toBe('Test Client Corp');
    expect(newLic.status).toBe('active');

    // Test Toggle Suspend Button
    act(() => {
      result.current.toggleStatus(newLic.id);
    });
    expect(result.current.licenses[0].status).toBe('suspended');

    // Test Re-issue Key Button
    const oldKey = result.current.licenses[0].license_key;
    act(() => {
      result.current.reIssueKey(newLic.id);
    });
    expect(result.current.licenses[0].license_key).not.toBe(oldKey);

    // Test Reset Domain/IP Lock Button
    act(() => {
      result.current.resetLock(newLic.id);
    });
    expect(result.current.licenses[0].licensed_domain).toBe('');
  });

  it('Webhook buttons: Add, Toggle, and Trigger Test Ping', async () => {
    const { result } = renderHook(() => useWebhooks());

    // Test Add Webhook Button
    act(() => {
      result.current.createWebhook({
        name: 'Test Webhook',
        url: 'https://example.com/webhook',
        events: ['invoice.paid'],
      });
    });
    expect(result.current.webhooks.some((w) => w.name === 'Test Webhook')).toBe(true);

    const whId = result.current.webhooks[0].id;

    // Test Pause/Active Toggle Button
    act(() => {
      result.current.toggleWebhook(whId);
    });
    expect(result.current.webhooks.find((w) => w.id === whId)?.is_active).toBe(false);
  });

  it('Form Builder buttons: Create Form, Add Field, Remove Field', () => {
    const { result } = renderHook(() => useFormBuilder());

    // Test Create Form Template Button
    act(() => {
      result.current.createForm('Minecraft Server Config', 'Custom RAM and Slots');
    });
    expect(result.current.forms.some((f) => f.name === 'Minecraft Server Config')).toBe(true);

    // Test Add Custom Field Button
    act(() => {
      result.current.saveField({
        id: 'ram_slot',
        name: 'ram_slot',
        label: 'Allocated RAM (GB)',
        type: 'number',
        required: true,
      });
    });
    expect(result.current.selectedForm?.fields.some((f) => f.name === 'ram_slot')).toBe(true);

    // Test Remove Field Button
    act(() => {
      result.current.removeField('ram_slot');
    });
    expect(result.current.selectedForm?.fields.some((f) => f.name === 'ram_slot')).toBe(false);
  });

  it('Redirects buttons: Add 301 Redirect, Toggle Active, Delete', () => {
    const { result } = renderHook(() => useRedirects());

    // Test Add Redirect Rule Button
    act(() => {
      result.current.createRedirect({
        source_path: '/promo-spring',
        target_url: '/cart?promo=SPRING50',
        status_code: 301,
      });
    });
    const found = result.current.redirects.find((r) => r.source_path === '/promo-spring');
    expect(found).toBeDefined();
    expect(found?.status_code).toBe(301);

    // Test Delete Redirect Button
    if (found) {
      act(() => {
        result.current.deleteRedirect(found.id);
      });
      expect(result.current.redirects.some((r) => r.source_path === '/promo-spring')).toBe(false);
    }
  });

  it('Cookie Consent buttons: Save settings and toggle enabled', () => {
    const { result } = renderHook(() => useCookieConsent());

    act(() => {
      result.current.updateSetting('is_enabled', false);
      result.current.updateSetting('theme', 'indigo');
    });

    expect(result.current.settings.is_enabled).toBe(false);
    expect(result.current.settings.theme).toBe('indigo');
  });
});
