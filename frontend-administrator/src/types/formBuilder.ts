export type FormFieldType = 'text' | 'textarea' | 'number' | 'select' | 'radio' | 'checkbox' | 'url';

export interface FormField {
  id: number;
  form_id: number;
  name: string;
  label: string;
  hide_label: boolean;
  description?: string;
  type: FormFieldType;
  default_value?: string;
  required: boolean;
  hidden: boolean;
  readonly: boolean;
  options?: Record<string, any>;
  prefix?: string;
  suffix?: string;
  text_size?: number;
  created_at?: string;
  updated_at?: string;
}

export interface FormStyle {
  type: string;
  show_title: boolean;
}

export interface CustomForm {
  id: number;
  name: string;
  style: FormStyle;
  fields: FormField[];
  created_at: string;
  updated_at: string;
}
