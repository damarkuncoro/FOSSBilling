export interface WidgetConfig {
  product_id: number;
  product_title: string;
  button_text: string;
  button_color: string;
  text_color: string;
  border_radius: number;
  layout: 'button' | 'pricing_card' | 'iframe_checkout';
  action_type: 'popup' | 'redirect';
  show_price: boolean;
  price_display: string;
}
