import { request } from '../lib/api/client';
import type { CustomForm, FormField } from '../types/formBuilder';

export interface IAdminFormBuilderRepository {
  getForms(): Promise<CustomForm[]>;
  getForm(id: number): Promise<CustomForm>;
  createForm(dto: Partial<CustomForm>): Promise<CustomForm>;
  updateForm(id: number, dto: Partial<CustomForm>): Promise<CustomForm>;
  deleteForm(id: number): Promise<any>;

  addField(formId: number, dto: Partial<FormField>): Promise<FormField>;
  updateField(fieldId: number, dto: Partial<FormField>): Promise<FormField>;
  deleteField(fieldId: number): Promise<any>;
}

export class AdminFormBuilderRepository implements IAdminFormBuilderRepository {
  async getForms(): Promise<CustomForm[]> {
    const res = await request<{ list: CustomForm[] }>('/admin/forms');
    return res.list;
  }

  async getForm(id: number): Promise<CustomForm> {
    return request<CustomForm>(`/admin/forms/${id}`);
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

  async addField(formId: number, dto: Partial<FormField>): Promise<FormField> {
    return request<FormField>(`/admin/forms/${formId}/fields`, {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async updateField(fieldId: number, dto: Partial<FormField>): Promise<FormField> {
    return request<FormField>(`/admin/forms/fields/${fieldId}`, {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }

  async deleteField(fieldId: number): Promise<any> {
    return request<any>(`/admin/forms/fields/${fieldId}`, {
      method: 'DELETE',
    });
  }
}

export const adminFormBuilderRepository = new AdminFormBuilderRepository();
