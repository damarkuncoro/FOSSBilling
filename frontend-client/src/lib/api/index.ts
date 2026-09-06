import { guestApi } from './guest';
import { clientPortalApi } from './client_portal';

export * from '@/types/api';
export * from './client';
export * from './guest';
export * from './client_portal';

// Unified API facade for backward compatibility
export const api = {
  ...guestApi,
  ...clientPortalApi,
};
