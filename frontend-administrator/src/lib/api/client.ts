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
