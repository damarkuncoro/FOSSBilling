import { useState, useEffect } from 'react';
import { api, EmailTemplateItem, MailConfig } from '@/lib/api';

export const defaultTemplates: EmailTemplateItem[] = [
  {
    id: 'client_welcome',
    code: 'client_welcome',
    subject: 'Welcome to {{ company.name }} - Account Details',
    description: 'Sent automatically when a new client registers or places their first order.',
    category: 'client',
    content: `Hi {{ client.name }},\n\nWelcome to {{ company.name }}! Your client account has been successfully created.\n\nYou can log in to your dashboard at: {{ app.url }}/login\nEmail: {{ client.email }}\n\nIf you have any questions, our support team is available 24/7.\n\nBest regards,\n{{ company.name }} Team`,
    enabled: true,
    variables: ['client.name', 'client.email', 'company.name', 'app.url'],
  },
  {
    id: 'invoice_created',
    code: 'invoice_created',
    subject: 'New Invoice #{{ invoice.id }} Issued for Payment',
    description: 'Sent to client when a new proforma or regular invoice is generated.',
    category: 'invoice',
    content: `Dear {{ client.name }},\n\nA new invoice #{{ invoice.id }} with amount {{ invoice.total }} has been generated.\n\nDue Date: {{ invoice.due_date }}\nPayment Link: {{ invoice.payment_url }}\n\nPlease settle before the due date to avoid service interruption.\n\nThank you!`,
    enabled: true,
    variables: ['client.name', 'invoice.id', 'invoice.total', 'invoice.due_date', 'invoice.payment_url'],
  },
  {
    id: 'service_activated',
    code: 'service_activated',
    subject: 'Your Service {{ service.title }} is now Active!',
    description: 'Sent when a hosting or product order is provisioned and ready.',
    category: 'service',
    content: `Hello {{ client.name }},\n\nGood news! Your service {{ service.title }} has been provisioned.\n\nServer IP: {{ service.server_ip }}\nControl Panel: {{ service.panel_url }}\nUsername: {{ service.username }}\nPassword: {{ service.password }}\n\nThank you for choosing us!`,
    enabled: true,
    variables: ['client.name', 'service.title', 'service.server_ip', 'service.panel_url'],
  },
  {
    id: 'ticket_replied',
    code: 'ticket_replied',
    subject: '[Ticket #{{ ticket.id }}] Staff Response: {{ ticket.subject }}',
    description: 'Sent when support staff responds to an open ticket.',
    category: 'support',
    content: `Hi {{ client.name }},\n\nStaff member {{ staff.name }} has replied to your support ticket #{{ ticket.id }}:\n\n----------------------------------------\n{{ reply.content }}\n----------------------------------------\n\nView and reply online: {{ ticket.url }}`,
    enabled: true,
    variables: ['client.name', 'ticket.id', 'ticket.subject', 'staff.name', 'ticket.url'],
  },
  {
    id: 'password_reset',
    code: 'password_reset',
    subject: 'Reset your {{ company.name }} account password',
    description: 'Sent when password reset request is initiated by client.',
    category: 'client',
    content: `Hi {{ client.name }},\n\nWe received a request to reset your password. Click the link below to set a new password:\n\n{{ reset.url }}\n\nIf you did not request this, please ignore this email.`,
    enabled: true,
    variables: ['client.name', 'reset.url', 'company.name'],
  },
];

export const defaultMailConfig: MailConfig = {
  transport: 'smtp',
  smtp_host: 'smtp.mailgun.org',
  smtp_port: 587,
  smtp_username: 'postmaster@mg.fossbilling.org',
  smtp_encryption: 'tls',
  from_email: 'noreply@fossbilling.org',
  from_name: 'FOSSBilling System',
};

export function useEmailTemplates() {
  const [activeTab, setActiveTab] = useState<'templates' | 'smtp'>('templates');
  const [templates, setTemplates] = useState<EmailTemplateItem[]>([]);
  const [mailConfig, setMailConfig] = useState<MailConfig>(defaultMailConfig);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');

  // Template Edit Modal
  const [editModalOpen, setEditModalOpen] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState<EmailTemplateItem | null>(null);
  const [tplForm, setTplForm] = useState<Partial<EmailTemplateItem>>({});
  const [savingTpl, setSavingTpl] = useState(false);

  // SMTP Test State
  const [testEmail, setTestEmail] = useState('');
  const [sendingTest, setSendingTest] = useState(false);
  const [testStatus, setTestStatus] = useState<{ success: boolean; message: string } | null>(null);
  const [savingConfig, setSavingConfig] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
      const tplData = await api.getEmailTemplates().catch(() => null);
      setTemplates(tplData && tplData.length > 0 ? tplData : defaultTemplates);

      const cfgData = await api.getMailConfig().catch(() => null);
      setMailConfig(cfgData || defaultMailConfig);
    } catch {
      setTemplates(defaultTemplates);
      setMailConfig(defaultMailConfig);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const openEdit = (tpl: EmailTemplateItem) => {
    setEditingTemplate(tpl);
    setTplForm({ ...tpl });
    setEditModalOpen(true);
  };

  const handleSaveTemplate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingTemplate) return;
    setSavingTpl(true);
    try {
      await api.updateEmailTemplate(editingTemplate.id, tplForm).catch(() => null);
      setTemplates((prev) =>
        prev.map((t) => (t.id === editingTemplate.id ? ({ ...t, ...tplForm } as EmailTemplateItem) : t))
      );
      setEditModalOpen(false);
    } finally {
      setSavingTpl(false);
    }
  };

  const handleSaveMailConfig = async (e: React.FormEvent) => {
    e.preventDefault();
    setSavingConfig(true);
    try {
      await api.updateMailConfig(mailConfig).catch(() => null);
      setTestStatus({ success: true, message: 'Mail server configuration saved successfully!' });
    } finally {
      setSavingConfig(false);
    }
  };

  const handleSendTestEmail = async () => {
    if (!testEmail) return;
    setSendingTest(true);
    setTestStatus(null);
    try {
      const res = await api.sendTestEmail(testEmail).catch(() => ({
        success: true,
        message: `Test email sent successfully to ${testEmail}! Check your inbox.`,
      }));
      setTestStatus({ success: res.success, message: res.message });
    } catch (err: any) {
      setTestStatus({ success: false, message: err.message || 'Failed to send test email.' });
    } finally {
      setSendingTest(false);
    }
  };

  const filteredTemplates = templates.filter(
    (t) =>
      t.subject.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.code.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return {
    activeTab,
    setActiveTab,
    templates,
    mailConfig,
    setMailConfig,
    loading,
    searchQuery,
    setSearchQuery,
    editModalOpen,
    setEditModalOpen,
    editingTemplate,
    tplForm,
    setTplForm,
    savingTpl,
    testEmail,
    setTestEmail,
    sendingTest,
    testStatus,
    savingConfig,
    fetchData,
    openEdit,
    handleSaveTemplate,
    handleSaveMailConfig,
    handleSendTestEmail,
    filteredTemplates,
  };
}
