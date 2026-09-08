import { AdminFormBuilderRepository, adminFormBuilderRepository, IAdminFormBuilderRepository } from '../repositories/admin_form_builder.repository';
import type { CustomForm, FormField } from '../types/formBuilder';

export class AdminFormBuilderService {
  constructor(private repo: IAdminFormBuilderRepository = adminFormBuilderRepository) {}

  async listForms(): Promise<CustomForm[]> {
    return this.repo.getForms();
  }

  async createForm(name: string, type: string = 'horizontal'): Promise<CustomForm> {
    if (!name.trim()) {
      throw new Error('Form name is required');
    }
    return this.repo.createForm({
      name: name.trim(),
      style: { type, show_title: true },
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

  async addField(formId: number, dto: Partial<FormField>): Promise<FormField> {
    return this.repo.addField(formId, dto);
  }

  async updateField(fieldId: number, dto: Partial<FormField>): Promise<FormField> {
    return this.repo.updateField(fieldId, dto);
  }

  async deleteField(fieldId: number): Promise<any> {
    return this.repo.deleteField(fieldId);
  }

  buildUpdatedFieldList(existingFields: FormField[], newOrUpdated: FormField): FormField[] {
    if (newOrUpdated.id) {
      const idx = existingFields.findIndex(f => f.id === newOrUpdated.id);
      if (idx !== -1) {
        const copy = [...existingFields];
        copy[idx] = { ...copy[idx], ...newOrUpdated };
        return copy;
      }
    }
    return [...existingFields, newOrUpdated];
  }
}

export const adminFormBuilderService = new AdminFormBuilderService();
