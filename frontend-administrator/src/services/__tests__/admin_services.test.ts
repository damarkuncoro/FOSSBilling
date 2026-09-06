import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AdminClientService } from '../admin_client.service';
import { AdminOrderService } from '../admin_order.service';
import { AdminInvoiceService } from '../admin_invoice.service';
import { AdminSupportService } from '../admin_support.service';
import { AdminMassMailService } from '../admin_massmail.service';
import { AdminStatsService } from '../admin_stats.service';
import { AdminProductService } from '../admin_product.service';
import { AdminServerService } from '../admin_server.service';
import { AdminGatewayService } from '../admin_gateway.service';
import { AdminCurrencyService } from '../admin_currency.service';
import { AdminCouponService } from '../admin_coupon.service';
import { AdminCompanyService } from '../admin_company.service';
import { AdminLicenseService } from '../admin_license.service';
import { AdminSystemService } from '../admin_system.service';
import { AdminWebhookService } from '../admin_webhook.service';
import { AdminFormBuilderService } from '../admin_form_builder.service';

describe('Admin Clean Architecture Services Suite', () => {
  describe('AdminClientService', () => {
    let mockRepo: any;
    let service: AdminClientService;

    beforeEach(() => {
      mockRepo = {
        listClients: vi.fn().mockResolvedValue([{ id: 1, first_name: 'John', last_name: 'Doe', email: 'john@example.com' }]),
        getClient: vi.fn().mockResolvedValue({ id: 1, first_name: 'John', last_name: 'Doe', email: 'john@example.com' }),
        createClient: vi.fn().mockResolvedValue({ id: 2, email: 'new@example.com' }),
        updateClient: vi.fn().mockResolvedValue({ id: 1, company: 'Acme Corp' }),
        deleteClient: vi.fn().mockResolvedValue({ success: true }),
      };
      service = new AdminClientService(mockRepo);
    });

    it('should validate and list clients', async () => {
      const clients = await service.listClients();
      expect(clients).toHaveLength(1);
      expect(mockRepo.listClients).toHaveBeenCalledTimes(1);
    });

    it('should reject invalid client ID', async () => {
      await expect(service.getClientDetail(0)).rejects.toThrow('Valid client ID is required');
    });

    it('should filter clients by search query', () => {
      const clients = [
        { id: 1, first_name: 'Alice', last_name: 'Smith', email: 'alice@test.com' } as any,
        { id: 2, first_name: 'Bob', last_name: 'Jones', email: 'bob@domain.org' } as any,
      ];
      const results = service.filterClients(clients, 'alice');
      expect(results).toHaveLength(1);
      expect(results[0].email).toBe('alice@test.com');
    });
  });

  describe('AdminOrderService', () => {
    let mockRepo: any;
    let service: AdminOrderService;

    beforeEach(() => {
      mockRepo = {
        listOrders: vi.fn().mockResolvedValue([{ id: 101, status: 'active', client_id: 1, title: 'VPS Hosting' }]),
        getOrder: vi.fn().mockResolvedValue({ id: 101, status: 'active' }),
        activateOrder: vi.fn().mockResolvedValue({ success: true }),
        suspendOrder: vi.fn().mockResolvedValue({ success: true }),
        unsuspendOrder: vi.fn().mockResolvedValue({ success: true }),
        cancelOrder: vi.fn().mockResolvedValue({ success: true }),
      };
      service = new AdminOrderService(mockRepo);
    });

    it('should require suspension reason', async () => {
      await expect(service.suspendOrder(101, '   ')).rejects.toThrow('Suspension reason is required');
      expect(mockRepo.suspendOrder).not.toHaveBeenCalled();
    });

    it('should activate order successfully', async () => {
      const res = await service.activateOrder(101);
      expect(res.success).toBe(true);
      expect(mockRepo.activateOrder).toHaveBeenCalledWith(101);
    });

    it('should filter orders by status', () => {
      const orders = [
        { id: 1, status: 'active' } as any,
        { id: 2, status: 'suspended' } as any,
      ];
      const activeOnly = service.filterByStatus(orders, 'active');
      expect(activeOnly).toHaveLength(1);
      expect(activeOnly[0].id).toBe(1);
    });
  });

  describe('AdminInvoiceService', () => {
    let mockRepo: any;
    let service: AdminInvoiceService;

    beforeEach(() => {
      mockRepo = {
        listInvoices: vi.fn().mockResolvedValue([{ id: 50, total: 100, status: 'unpaid' }]),
        createInvoice: vi.fn().mockResolvedValue({ id: 51, total: 150, status: 'unpaid' }),
      };
      service = new AdminInvoiceService(mockRepo);
    });

    it('should validate invoice items before creation', async () => {
      await expect(service.createInvoice(1, [])).rejects.toThrow('Invoice must contain at least one line item');
    });

    it('should create invoice successfully with valid payload', async () => {
      const res = await service.createInvoice(1, [{ title: 'Dedicated Server', price: 150 }]);
      expect(res.id).toBe(51);
      expect(mockRepo.createInvoice).toHaveBeenCalledTimes(1);
    });
  });

  describe('AdminMassMailService', () => {
    let mockRepo: any;
    let service: AdminMassMailService;

    beforeEach(() => {
      mockRepo = {
        sendBroadcast: vi.fn().mockResolvedValue({ queued: 42 }),
      };
      service = new AdminMassMailService(mockRepo);
    });

    it('should validate subject and content', async () => {
      await expect(service.sendMassMail('', 'Hello')).rejects.toThrow('Email subject is required');
      await expect(service.sendMassMail('Announcement', '')).rejects.toThrow('Email body content is required');
    });

    it('should trigger broadcast when params are valid', async () => {
      const res = await service.sendMassMail('System Upgrade', 'Maintenance scheduled at midnight');
      expect(res.queued).toBe(42);
    });
  });

  describe('AdminProductService', () => {
    let mockRepo: any;
    let service: AdminProductService;

    beforeEach(() => {
      mockRepo = {
        getProducts: vi.fn().mockResolvedValue([{ id: 1, title: 'VPS Pro', slug: 'vps-pro', type: 'hosting' }]),
        getProduct: vi.fn().mockResolvedValue({ id: 1, title: 'VPS Pro' }),
        createProduct: vi.fn().mockResolvedValue({ id: 2, title: 'Shared Hosting', slug: 'shared-hosting' }),
        updateProduct: vi.fn().mockResolvedValue({ id: 1, title: 'VPS Ultra' }),
        deleteProduct: vi.fn().mockResolvedValue({ success: true }),
        getProductCategories: vi.fn().mockResolvedValue([{ id: 1, title: 'Cloud' }]),
      };
      service = new AdminProductService(mockRepo);
    });

    it('should validate and create product', async () => {
      await expect(service.createProduct({ title: '' })).rejects.toThrow('Product title is required');
      const res = await service.createProduct({ title: 'Shared Hosting', price_monthly: 10 });
      expect(res.id).toBe(2);
      expect(mockRepo.createProduct).toHaveBeenCalledTimes(1);
    });

    it('should filter products by search and type', () => {
      const items = [
        { id: 1, title: 'Cloud Alpha', slug: 'cloud-alpha', type: 'hosting', description: 'Fast' },
        { id: 2, title: 'Dot COM', slug: 'dot-com', type: 'domain', description: 'TLD' },
      ] as any;
      const res = service.filterProducts(items, 'alpha', 'hosting');
      expect(res).toHaveLength(1);
      expect(res[0].id).toBe(1);
    });
  });

  describe('AdminServerService', () => {
    let mockRepo: any;
    let service: AdminServerService;

    beforeEach(() => {
      mockRepo = {
        getServers: vi.fn().mockResolvedValue([{ id: 1, name: 'cPanel Srv 1', status: 'online' }]),
        createServer: vi.fn().mockResolvedValue({ id: 2, name: 'Hestia Node 1' }),
        testServerConnection: vi.fn().mockResolvedValue({ success: true, message: 'OK' }),
      };
      service = new AdminServerService(mockRepo);
    });

    it('should validate server inputs on creation', async () => {
      await expect(service.createServer({ name: '' })).rejects.toThrow('Server name is required');
      await expect(service.createServer({ name: 'Node 1' })).rejects.toThrow('Server hostname or IP address is required');
    });

    it('should test connection successfully', async () => {
      const res = await service.testConnection(1);
      expect(res.success).toBe(true);
      expect(mockRepo.testServerConnection).toHaveBeenCalledWith(1);
    });
  });

  describe('AdminGatewayService', () => {
    let mockRepo: any;
    let service: AdminGatewayService;

    beforeEach(() => {
      mockRepo = {
        getPaymentGateways: vi.fn().mockResolvedValue([{ id: 'stripe', name: 'Stripe' }]),
        updatePaymentGateway: vi.fn().mockResolvedValue({ id: 'stripe', enabled: true }),
        getTaxRules: vi.fn().mockResolvedValue([{ id: 1, name: 'VAT', rate: 20 }]),
        createTaxRule: vi.fn().mockResolvedValue({ id: 2, name: 'PPN', rate: 11 }),
        deleteTaxRule: vi.fn().mockResolvedValue({ success: true }),
      };
      service = new AdminGatewayService(mockRepo);
    });

    it('should validate tax rule creation', async () => {
      await expect(service.createTaxRule({ name: '' })).rejects.toThrow('Tax rule name is required');
      const res = await service.createTaxRule({ name: 'PPN', rate: 11 });
      expect(res.id).toBe(2);
    });
  });

  describe('AdminCurrencyService', () => {
    let mockRepo: any;
    let service: AdminCurrencyService;

    beforeEach(() => {
      mockRepo = {
        getCurrencies: vi.fn().mockResolvedValue([{ code: 'USD', title: 'US Dollar', conversion_rate: 1.0 }]),
        createCurrency: vi.fn().mockResolvedValue({ code: 'EUR', title: 'Euro', conversion_rate: 0.92 }),
        setDefaultCurrency: vi.fn().mockResolvedValue({ success: true }),
        deleteCurrency: vi.fn().mockResolvedValue({ success: true }),
      };
      service = new AdminCurrencyService(mockRepo);
    });

    it('should validate currency code length', async () => {
      await expect(service.createCurrency({ code: 'US', title: 'Dollar' })).rejects.toThrow(
        'A valid 3-letter currency code (e.g. USD, IDR) is required'
      );
    });

    it('should create currency with uppercase code', async () => {
      const res = await service.createCurrency({ code: 'eur', title: 'Euro', conversion_rate: 0.92 });
      expect(res.code).toBe('EUR');
      expect(mockRepo.createCurrency).toHaveBeenCalledWith(
        expect.objectContaining({ code: 'EUR', title: 'Euro' })
      );
    });
  });

  describe('AdminCouponService', () => {
    let mockRepo: any;
    let service: AdminCouponService;

    beforeEach(() => {
      mockRepo = {
        getCoupons: vi.fn().mockResolvedValue([{ id: 1, code: 'SAVE10', value: 10 }]),
        createCoupon: vi.fn().mockResolvedValue({ id: 2, code: 'FLASH50', value: 50 }),
        deleteCoupon: vi.fn().mockResolvedValue({ success: true }),
      };
      service = new AdminCouponService(mockRepo);
    });

    it('should validate coupon code and discount value', async () => {
      await expect(service.createCoupon({ code: '', value: 10 })).rejects.toThrow('Coupon code is required');
      await expect(service.createCoupon({ code: 'DISC', value: 0 })).rejects.toThrow(
        'Coupon discount value must be greater than zero'
      );
    });

    it('should generate random coupon code', () => {
      const code = service.generateRandomCode('SUMMER');
      expect(code.startsWith('SUMMER')).toBe(true);
      expect(code.length).toBeGreaterThan(6);
    });
  });

  describe('AdminCompanyService', () => {
    let mockRepo: any;
    let service: AdminCompanyService;

    beforeEach(() => {
      mockRepo = {
        getCompany: vi.fn().mockResolvedValue({ name: 'FOSSBilling Corp', email: 'billing@fossbilling.org' }),
        updateCompany: vi.fn().mockResolvedValue({ name: 'FOSSBilling Inc', email: 'support@fossbilling.org' }),
      };
      service = new AdminCompanyService(mockRepo);
    });

    it('should validate email format on update', async () => {
      await expect(service.updateCompany({ email: 'invalid-email' })).rejects.toThrow(
        'Valid company email address is required'
      );
    });

    it('should update company settings successfully', async () => {
      const res = await service.updateCompany({ name: 'FOSSBilling Inc', email: 'support@fossbilling.org' });
      expect(res.name).toBe('FOSSBilling Inc');
      expect(mockRepo.updateCompany).toHaveBeenCalledTimes(1);
    });
  });

  describe('AdminLicenseService', () => {
    let mockRepo: any;
    let service: AdminLicenseService;

    beforeEach(() => {
      mockRepo = {
        getLicenses: vi.fn().mockResolvedValue([{ id: 1, license_key: 'FOSS-ENT-1122', status: 'active' }]),
        createLicense: vi.fn().mockResolvedValue({ id: 2, license_key: 'FOSS-PRO-3344', status: 'active' }),
        resetLock: vi.fn().mockResolvedValue({ success: true }),
        getValidationLogs: vi.fn().mockResolvedValue([{ id: 1, license_key: 'FOSS-ENT-1122', result: 'valid' }]),
      };
      service = new AdminLicenseService(mockRepo);
    });

    it('should generate valid license key format', () => {
      const key = service.generateKey('TEST');
      expect(key.startsWith('TEST-')).toBe(true);
      expect(key.split('-')).toHaveLength(5);
    });

    it('should reset lock successfully', async () => {
      const res = await service.resetLock(1);
      expect(res.success).toBe(true);
      expect(mockRepo.resetLock).toHaveBeenCalledWith(1);
    });

    it('should filter licenses by search query and status', () => {
      const lics = [
        { id: 1, license_key: 'FOSS-ENT-1111', client_name: 'Alpha', status: 'active' },
        { id: 2, license_key: 'FOSS-VPN-2222', client_name: 'Beta', status: 'suspended' },
      ] as any;
      const res = service.filterLicenses(lics, 'alpha', 'active');
      expect(res).toHaveLength(1);
      expect(res[0].id).toBe(1);
    });
  });

  describe('AdminSystemService', () => {
    let mockRepo: any;
    let service: AdminSystemService;

    beforeEach(() => {
      mockRepo = {
        getSystemInfo: vi.fn().mockResolvedValue({ version: '1.0.0', database: 'postgres' }),
        getSystemStatus: vi.fn().mockResolvedValue({ engine_version: 'v0.9.0', cron_status: 'healthy' }),
        triggerCron: vi.fn().mockResolvedValue({ success: true, message: 'Cron executed' }),
        clearCache: vi.fn().mockResolvedValue({ success: true, message: 'Cache cleared' }),
        createNews: vi.fn().mockResolvedValue({ id: 1, title: 'Release 1.0' }),
        deleteNews: vi.fn().mockResolvedValue({ success: true }),
      };
      service = new AdminSystemService(mockRepo);
    });

    it('should trigger cron and clear cache', async () => {
      const cronRes = await service.triggerCron();
      expect(cronRes.success).toBe(true);
      const cacheRes = await service.clearCache();
      expect(cacheRes.success).toBe(true);
    });

    it('should validate news article inputs', async () => {
      await expect(service.createNewsArticle('', 'Content')).rejects.toThrow('Article title is required');
      await expect(service.createNewsArticle('Title', '')).rejects.toThrow('Article content is required');
    });
  });

  describe('AdminWebhookService', () => {
    let mockRepo: any;
    let service: AdminWebhookService;

    beforeEach(() => {
      mockRepo = {
        getWebhooks: vi.fn().mockResolvedValue([{ id: 1, name: 'Discord', url: 'https://discord.com/hook' }]),
        createWebhook: vi.fn().mockResolvedValue({ id: 2, name: 'Slack', url: 'https://slack.com/hook' }),
        deleteWebhook: vi.fn().mockResolvedValue({ success: true }),
        triggerTestPing: vi.fn().mockResolvedValue({ success: true }),
      };
      service = new AdminWebhookService(mockRepo);
    });

    it('should validate webhook name and URL', async () => {
      await expect(service.createWebhook('', 'https://test.com', ['order.created'])).rejects.toThrow(
        'Webhook name is required'
      );
      await expect(service.createWebhook('Slack', 'invalid-url', ['order.created'])).rejects.toThrow(
        'Valid webhook URL (HTTP/HTTPS) is required'
      );
      await expect(service.createWebhook('Slack', 'https://slack.com', [])).rejects.toThrow(
        'At least one event trigger must be selected'
      );
    });

    it('should create webhook with generated secret', async () => {
      const res = await service.createWebhook('Slack', 'https://slack.com/hook', ['invoice.paid']);
      expect(res.id).toBe(2);
      expect(mockRepo.createWebhook).toHaveBeenCalledWith(
        expect.objectContaining({ name: 'Slack', url: 'https://slack.com/hook' })
      );
    });
  });

  describe('AdminFormBuilderService', () => {
    let mockRepo: any;
    let service: AdminFormBuilderService;

    beforeEach(() => {
      mockRepo = {
        getForms: vi.fn().mockResolvedValue([{ id: 1, name: 'VPS Setup', fields: [] }]),
        createForm: vi.fn().mockResolvedValue({ id: 2, name: 'Minecraft Config', fields: [] }),
        deleteForm: vi.fn().mockResolvedValue({ success: true }),
      };
      service = new AdminFormBuilderService(mockRepo);
    });

    it('should validate form name on creation', async () => {
      await expect(service.createForm('', 'Description')).rejects.toThrow('Form name is required');
      const res = await service.createForm('Minecraft Config', 'Game server config');
      expect(res.id).toBe(2);
    });

    it('should add new field or update existing field in list', () => {
      const initialFields: any[] = [{ id: 'f1', name: 'hostname', label: 'Host' }];
      const updated = service.buildUpdatedFieldList(initialFields, { id: 'f1', name: 'hostname', label: 'Updated Host' } as any);
      expect(updated[0].label).toBe('Updated Host');

      const withNew = service.buildUpdatedFieldList(initialFields, { name: 'port', label: 'Port' } as any);
      expect(withNew).toHaveLength(2);
    });
  });
});
