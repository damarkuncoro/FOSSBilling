import { clientsApi } from './clients';
import { catalogApi } from './catalog';
import { billingApi } from './billing';
import { systemApi } from './system';

export * from '@/types/api';
export * from '@/types/modules';
export * from './client';
export * from './clients';
export * from './catalog';
export * from './billing';
export * from './system';

// Unified API facade for backward compatibility
export const api = {
  ...clientsApi,
  ...catalogApi,
  ...billingApi,
  ...systemApi,
};
