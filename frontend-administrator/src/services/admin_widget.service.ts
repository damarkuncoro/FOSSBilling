import { AdminWidgetRepositoryInterface } from '../repositories/admin_widget.repository';

export class AdminWidgetService {
  constructor(private readonly repo: AdminWidgetRepositoryInterface) {}

  async getAllWidgets(): Promise<Record<string, any[]>> {
    return this.repo.listWidgets();
  }

  async addWidget(widget: { id: string; slot: string; title: string; component: string; priority?: number }): Promise<any> {
    if (!widget.id || !widget.slot || !widget.component) {
      throw new Error('Widget ID, slot, and component are required');
    }
    return this.repo.registerWidget(widget);
  }
}
