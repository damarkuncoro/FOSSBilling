import { request } from '../lib/api/client';
import type { CustomForm } from '../types/formBuilder';

export interface IAdminFormBuilderRepository {
  getForms(): Promise<CustomForm[]>;
  createForm(dto: Partial<CustomForm>): Promise<CustomForm>;
  updateForm(id: number, dto: Partial<CustomForm>): Promise<CustomForm>;
  deleteForm(id: number): Promise<any>;
}

export class AdminFormBuilderRepository implements IAdminFormBuilderRepository {
  async getForms(): Promise<CustomForm[]> {
    return request<CustomForm[]>('/admin/forms');
  }

  async createForm(dto: Partial<CustomForm>): Promise<CustomForm> {
    return request<CustomForm>('/admin/forms', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async updateForm(id: number, dto: Partial<CustomForm>): Promise<CustomForm> {
    return request<CustomForm>(`/admin/forms/${id}`, {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }

  async deleteForm(id: number): Promise<any> {
    return request<any>(`/admin/forms/${id}`, {
      method: 'DELETE',
    });
  }
}

export const adminFormBuilderRepository = new AdminFormBuilderRepository();
