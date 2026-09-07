import { ApiResponse } from '@/types/api';

export const API_BASE = '/api/v1';

export class ApiError extends Error {
  code: string;
  details?: any;

  constructor(message: string, code = 'API_ERROR', details?: any) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.details = details;
  }
}

export const getStoredToken = () => localStorage.getItem('fossbilling_admin_token');
export const setStoredToken = (t: string) => localStorage.setItem('fossbilling_admin_token', t);
export const removeStoredToken = () => {
  localStorage.removeItem('fossbilling_admin_token');
  localStorage.removeItem('fossbilling_admin_user');
};

export async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = getStoredToken();
  const headers = new Headers(options.headers || {});
  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const url = endpoint.startsWith('http') ? endpoint : `${API_BASE}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`;
  const response = await fetch(url, { ...options, headers });
  const json: ApiResponse<T> = await response.json().catch(() => {
    throw new ApiError(`HTTP Error: ${response.status} ${response.statusText}`, 'HTTP_ERROR');
  });

  if (!response.ok || !json.success) {
    throw new ApiError(
      json.error?.message || `Status ${response.status}`,
      json.error?.code || 'UNKNOWN_ERROR',
      json.error?.details
    );
  }
  return json.data;
}

export const apiClient = {
  get: async <T>(url: string, params?: Record<string, any>) => {
    let finalUrl = url;
    if (params) {
      const q = new URLSearchParams();
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== null) q.append(k, String(v));
      });
      const qs = q.toString();
      if (qs) finalUrl += `?${qs}`;
    }
    const data = await request<T>(finalUrl, { method: 'GET' });
    return { data };
  },
  post: async <T>(url: string, body?: any) => {
    const data = await request<T>(url, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    });
    return { data };
  },
  put: async <T>(url: string, body?: any) => {
    const data = await request<T>(url, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    });
    return { data };
  },
  delete: async <T>(url: string) => {
    const data = await request<T>(url, { method: 'DELETE' });
    return { data };
  },
  patch: async <T>(url: string, body?: any) => {
    const data = await request<T>(url, {
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    });
    return { data };
  },
};

