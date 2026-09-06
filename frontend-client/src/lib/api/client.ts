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

export function getStoredClientToken(): string | null {
  return localStorage.getItem('fossbilling_client_token');
}

export function setStoredClientToken(token: string) {
  localStorage.setItem('fossbilling_client_token', token);
}

export function removeStoredClientToken() {
  localStorage.removeItem('fossbilling_client_token');
  localStorage.removeItem('fossbilling_client_user');
}

export async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = getStoredClientToken();
  const headers = new Headers(options.headers || {});

  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }

  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const url = endpoint.startsWith('http') ? endpoint : `${API_BASE}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`;

  const response = await fetch(url, {
    ...options,
    headers,
  });

  const json: ApiResponse<T> = await response.json().catch(() => {
    throw new ApiError(`HTTP Error: ${response.status} ${response.statusText}`, 'HTTP_ERROR');
  });

  if (!response.ok || !json.success) {
    const errorMsg = json.error?.message || `Request failed with status ${response.status}`;
    const errorCode = json.error?.code || 'UNKNOWN_ERROR';
    throw new ApiError(errorMsg, errorCode, json.error?.details);
  }

  return json.data;
}
