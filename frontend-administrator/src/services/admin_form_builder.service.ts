import { AdminFormBuilderRepository, adminFormBuilderRepository, IAdminFormBuilderRepository } from '../repositories/admin_form_builder.repository';
import type { CustomForm, FormField } from '../types/formBuilder';

export class AdminFormBuilderService {
  constructor(private repo: IAdminFormBuilderRepository = adminFormBuilderRepository) {}

  async listForms(): Promise<CustomForm[]> {
    return this.repo.getForms();
  }

  async createForm(name: string, description: string): Promise<CustomForm> {
    if (!name.trim()) {
      throw new Error('Form name is required');
    }
    return this.repo.createForm({
      name: name.trim(),
      description: description.trim(),
      product_ids: [],
      fields: [],
    });
  }

  async updateForm(id: number, dto: Partial<CustomForm>): Promise<CustomForm> {
    if (!id || id <= 0) {
      throw new Error('Valid form ID is required');
    }
    return this.repo.updateForm(id, dto);
  }

  async deleteForm(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid form ID is required');
    }
    return this.repo.deleteForm(id);
  }

  buildUpdatedFieldList(fields: FormField[], field: FormField): FormField[] {
    const exists = fields.some((f) => f.id === field.id);
    if (exists) {
      return fields.map((f) => (f.id === field.id ? field : f));
    }
    return [...fields, { ...field, id: field.id || `f_${Date.now()}` }];
  }
}

export const adminFormBuilderService = new AdminFormBuilderService();
