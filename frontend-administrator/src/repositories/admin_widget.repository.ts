import { apiClient } from '@/lib/api';

export interface AdminWidgetRepositoryInterface {
  listWidgets(): Promise<Record<string, any[]>>;
  registerWidget(widget: any): Promise<any>;
}

export class AdminWidgetRepository implements AdminWidgetRepositoryInterface {
  async listWidgets(): Promise<Record<string, any[]>> {
    const res = await apiClient.get<any>('/admin/widgets');
    return res.data?.data || res.data || {};
  }

  async registerWidget(widget: any): Promise<any> {
    const res = await apiClient.post<any>('/admin/widgets', widget);
    return res.data?.data || res.data;
  }
}
