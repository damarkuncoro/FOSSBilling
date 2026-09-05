export interface CookieConsentSettings {
  is_enabled: boolean;
  banner_position: 'bottom' | 'top' | 'floating_left' | 'floating_right';
  theme: 'dark' | 'light' | 'indigo';
  message: string;
  accept_button_text: string;
  decline_button_text: string;
  show_decline_button: boolean;
  privacy_policy_url: string;
  cookie_expiration_days: number;
}

export interface ConsentLog {
  id: number;
  ip_address: string;
  country: string;
  decision: 'accepted' | 'declined';
  created_at: string;
}
