export interface WebhookEndpoint {
  id: number;
  name: string;
  url: string;
  secret: string;
  events: string[];
  is_active: boolean;
  total_deliveries: number;
  last_status_code?: number;
  last_delivery_at?: string;
  created_at: string;
}

export interface WebhookDeliveryLog {
  id: number;
  webhook_id: number;
  event: string;
  url: string;
  response_code: number;
  response_body?: string;
  execution_time_ms: number;
  created_at: string;
}
