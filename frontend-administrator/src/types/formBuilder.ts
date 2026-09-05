export type FormFieldType = 'text' | 'textarea' | 'number' | 'dropdown' | 'radio' | 'checkbox';

export interface FormField {
  id: string;
  name: string;
  label: string;
  type: FormFieldType;
  required: boolean;
  placeholder?: string;
  description?: string;
  options?: string[]; // for dropdown, radio, checkbox
  default_value?: string;
}

export interface CustomForm {
  id: number;
  name: string;
  description: string;
  product_ids: number[];
  fields: FormField[];
  created_at: string;
  updated_at: string;
}
