import { ApiResponse } from '@/types/api';

export const API_BASE = '/api/v1';

export class ApiError extends Error {
  code: string;
  details?: any;
  status?: number;

  constructor(message: string, code = 'API_ERROR', details?: any, status?: number) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.details = details;
    this.status = status;
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

/**
 * requestFull returns the entire ApiResponse object
 */
export async function requestFull<T>(endpoint: string, options: RequestInit = {}): Promise<ApiResponse<T>> {
  const token = getStoredClientToken();
  const headers = new Headers(options.headers || {});

  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }

  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const url = endpoint.startsWith('http') ? endpoint : `${API_BASE}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`;

  try {
    const response = await fetch(url, { ...options, headers });

    if (response.status === 401 && !endpoint.includes('/auth/login')) {
      removeStoredClientToken();
      if (!window.location.pathname.includes('/login')) {
        window.location.href = '/login';
      }
      throw new ApiError('Unauthorized', 'UNAUTHORIZED', null, 401);
    }

    const json: ApiResponse<T> = await response.json().catch(() => {
      throw new ApiError(`HTTP Error: ${response.status} ${response.statusText}`, 'HTTP_ERROR', null, response.status);
    });

    if (!response.ok || !json.success) {
      throw new ApiError(
        json.error?.message || `Status ${response.status}`,
        json.error?.code || 'UNKNOWN_ERROR',
        json.error?.details,
        response.status
      );
    }

    return json;
  } catch (err: any) {
    if (err instanceof ApiError) throw err;
    throw new ApiError(err.message || 'Network error', 'NETWORK_ERROR');
  }
}

/**
 * request returns only the .data portion of the response (default behavior)
 */
export async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const res = await requestFull<T>(endpoint, options);
  return res.data;
}
