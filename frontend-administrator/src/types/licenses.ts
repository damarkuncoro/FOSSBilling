export interface SoftwareLicense {
  id: number;
  client_id: number;
  client_name: string;
  product_id: number;
  product_title: string;
  license_key: string;
  status: 'active' | 'suspended' | 'expired' | 'revoked';
  licensed_domain?: string;
  licensed_ip?: string;
  version?: string;
  max_instances: number;
  instances_count: number;
  expires_at?: string;
  created_at: string;
  updated_at: string;
}

export interface LicenseValidationLog {
  id: number;
  license_key: string;
  ip_address: string;
  domain: string;
  result: 'valid' | 'invalid_domain' | 'invalid_ip' | 'expired' | 'suspended';
  created_at: string;
}
