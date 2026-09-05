import React, { useEffect, useState } from 'react';
import {
  Building2,
  Save,
  RefreshCw,
  CheckCircle,
  AlertCircle,
  Globe,
  MapPin,
  Image as ImageIcon,
  Mail,
  ShieldCheck,
  FileText,
} from 'lucide-react';
import { api, CompanySettings } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

export const Company: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const [form, setForm] = useState<CompanySettings>({
    name: '',
    email: '',
    phone: '',
    address_1: '',
    address_2: '',
    city: '',
    state: '',
    postcode: '',
    country: '',
    vat_number: '',
    logo_url: '',
    logo_dark_url: '',
    favicon_url: '',
    terms_url: '',
    email_signature: '',
  });

  const fetchCompanySettings = async () => {
    setLoading(true);
    setErrorMessage(null);
    try {
      const data = await api.getCompany();
      if (data) {
        setForm({
          name: data.name || '',
          email: data.email || '',
          phone: data.phone || '',
          address_1: data.address_1 || '',
          address_2: data.address_2 || '',
          city: data.city || '',
          state: data.state || '',
          postcode: data.postcode || '',
          country: data.country || '',
          vat_number: data.vat_number || '',
          logo_url: data.logo_url || '',
          logo_dark_url: data.logo_dark_url || '',
          favicon_url: data.favicon_url || '',
          terms_url: data.terms_url || '',
          email_signature: data.email_signature || '',
          updated_at: data.updated_at,
        });
      }
    } catch (err: any) {
      setErrorMessage(err.message || 'Failed to load company settings');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCompanySettings();
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setSuccessMessage(null);
    setErrorMessage(null);

    try {
      const res = await api.updateCompany(form);
      if (res) {
        setForm((prev) => ({
          ...prev,
          ...res,
        }));
      }
      setSuccessMessage('Company settings updated successfully.');
      setTimeout(() => setSuccessMessage(null), 4000);
    } catch (err: any) {
      setErrorMessage(err.message || 'Failed to save settings');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
        <RefreshCw className="w-8 h-8 animate-spin text-primary" />
        <p className="text-muted-foreground text-sm">Loading company profile & branding...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-6xl mx-auto pb-12">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b pb-4">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-3xl font-bold tracking-tight">Company Settings</h1>
            <Badge variant="outline" className="text-xs bg-primary/10 text-primary border-primary/20">
              Official Branding
            </Badge>
          </div>
          <p className="text-muted-foreground mt-1 text-sm">
            Manage your organization identity, public invoice details, branding logos, and client communications.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" onClick={fetchCompanySettings} disabled={saving}>
            <RefreshCw className={`w-4 h-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
            Reload
          </Button>
          <Button onClick={handleSave} disabled={saving} size="sm" className="gap-2">
            <Save className="w-4 h-4" />
            {saving ? 'Saving Changes...' : 'Save Settings'}
          </Button>
        </div>
      </div>

      {/* Status Notifications */}
      {successMessage && (
        <div className="flex items-center gap-3 p-4 rounded-lg bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30">
          <CheckCircle className="w-5 h-5 flex-shrink-0" />
          <p className="text-sm font-medium">{successMessage}</p>
        </div>
      )}

      {errorMessage && (
        <div className="flex items-center gap-3 p-4 rounded-lg bg-destructive/15 text-destructive border border-destructive/30">
          <AlertCircle className="w-5 h-5 flex-shrink-0" />
          <p className="text-sm font-medium">{errorMessage}</p>
        </div>
      )}

      {/* Tabs Container */}
      <form onSubmit={handleSave}>
        <Tabs defaultValue="general" className="w-full space-y-6">
          <TabsList className="grid grid-cols-2 md:grid-cols-4 w-full h-auto p-1.5 gap-1 bg-muted/60">
            <TabsTrigger value="general" className="gap-2 py-2.5">
              <Building2 className="w-4 h-4" />
              General Info
            </TabsTrigger>
            <TabsTrigger value="address" className="gap-2 py-2.5">
              <MapPin className="w-4 h-4" />
              Address & Tax
            </TabsTrigger>
            <TabsTrigger value="branding" className="gap-2 py-2.5">
              <ImageIcon className="w-4 h-4" />
              Logos & Identity
            </TabsTrigger>
            <TabsTrigger value="signature" className="gap-2 py-2.5">
              <Mail className="w-4 h-4" />
              Email Signature
            </TabsTrigger>
          </TabsList>

          {/* Tab 1: General Info */}
          <TabsContent value="general" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Building2 className="w-5 h-5 text-primary" />
                  Organization Details
                </CardTitle>
                <CardDescription>
                  Basic profile used on invoices, system emails, and public client portal headers.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Company / Brand Name <span className="text-destructive">*</span></label>
                    <Input
                      name="name"
                      value={form.name}
                      onChange={handleChange}
                      placeholder="e.g. FOSSBilling Cloud Solutions"
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Official Email Address <span className="text-destructive">*</span></label>
                    <Input
                      type="email"
                      name="email"
                      value={form.email}
                      onChange={handleChange}
                      placeholder="billing@example.com"
                      required
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Contact Phone Number</label>
                    <Input
                      name="phone"
                      value={form.phone}
                      onChange={handleChange}
                      placeholder="+62 812-3456-7890"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Terms & Conditions URL</label>
                    <Input
                      name="terms_url"
                      value={form.terms_url}
                      onChange={handleChange}
                      placeholder="https://example.com/terms"
                    />
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Tab 2: Address & Tax */}
          <TabsContent value="address" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <MapPin className="w-5 h-5 text-primary" />
                  Registered Business Address & Tax / VAT
                </CardTitle>
                <CardDescription>
                  This address and tax identification number will be printed on all generated PDF invoices and receipts.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Street Address Line 1</label>
                  <Input
                    name="address_1"
                    value={form.address_1}
                    onChange={handleChange}
                    placeholder="e.g. Jl. Sudirman No. 88, SCBD Area"
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Address Line 2 (Optional)</label>
                  <Input
                    name="address_2"
                    value={form.address_2}
                    onChange={handleChange}
                    placeholder="e.g. Tower 2, Floor 18"
                  />
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium">City</label>
                    <Input
                      name="city"
                      value={form.city}
                      onChange={handleChange}
                      placeholder="Jakarta Selatan"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">State / Province</label>
                    <Input
                      name="state"
                      value={form.state}
                      onChange={handleChange}
                      placeholder="DKI Jakarta"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Postal / Zip Code</label>
                    <Input
                      name="postcode"
                      value={form.postcode}
                      onChange={handleChange}
                      placeholder="12190"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium">Country Code</label>
                    <Input
                      name="country"
                      value={form.country}
                      onChange={handleChange}
                      placeholder="ID (2-letter ISO code)"
                      maxLength={2}
                    />
                  </div>
                </div>

                <div className="pt-2">
                  <div className="space-y-2">
                    <label className="text-sm font-medium flex items-center gap-2">
                      <ShieldCheck className="w-4 h-4 text-primary" />
                      Tax / VAT Identification Number (NPWP / TIN)
                    </label>
                    <Input
                      name="vat_number"
                      value={form.vat_number}
                      onChange={handleChange}
                      placeholder="e.g. 01.234.567.8-901.000 / EU123456789"
                    />
                    <p className="text-xs text-muted-foreground">
                      Included in invoice headers to comply with local tax regulations.
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Tab 3: Branding & Logos */}
          <TabsContent value="branding" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <ImageIcon className="w-5 h-5 text-primary" />
                  Brand Assets & Logos
                </CardTitle>
                <CardDescription>
                  Configure logos for light and dark modes, as well as browser favicon.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* Light Mode Logo */}
                  <div className="space-y-3">
                    <label className="text-sm font-medium">Logo URL (Light Mode)</label>
                    <Input
                      name="logo_url"
                      value={form.logo_url}
                      onChange={handleChange}
                      placeholder="/branding/logo-light.svg or https://..."
                    />
                    <div className="p-4 rounded-lg bg-white border flex flex-col items-center justify-center min-h-[120px]">
                      <span className="text-xs text-slate-400 mb-2 font-mono">Light Background Preview</span>
                      {form.logo_url ? (
                        <img
                          src={form.logo_url}
                          alt="Light Logo"
                          className="max-h-12 max-w-full object-contain"
                          onError={(e) => {
                            (e.target as HTMLElement).style.display = 'none';
                          }}
                        />
                      ) : (
                        <span className="text-xs text-muted-foreground">No logo configured</span>
                      )}
                    </div>
                  </div>

                  {/* Dark Mode Logo */}
                  <div className="space-y-3">
                    <label className="text-sm font-medium">Logo URL (Dark Mode)</label>
                    <Input
                      name="logo_dark_url"
                      value={form.logo_dark_url}
                      onChange={handleChange}
                      placeholder="/branding/logo-dark.svg or https://..."
                    />
                    <div className="p-4 rounded-lg bg-slate-900 border border-slate-800 flex flex-col items-center justify-center min-h-[120px]">
                      <span className="text-xs text-slate-500 mb-2 font-mono">Dark Background Preview</span>
                      {form.logo_dark_url ? (
                        <img
                          src={form.logo_dark_url}
                          alt="Dark Logo"
                          className="max-h-12 max-w-full object-contain"
                          onError={(e) => {
                            (e.target as HTMLElement).style.display = 'none';
                          }}
                        />
                      ) : (
                        <span className="text-xs text-slate-500">No logo configured</span>
                      )}
                    </div>
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">Favicon URL</label>
                  <div className="flex gap-4 items-center">
                    <Input
                      name="favicon_url"
                      value={form.favicon_url}
                      onChange={handleChange}
                      placeholder="/branding/favicon.svg or https://..."
                      className="flex-1"
                    />
                    {form.favicon_url && (
                      <div className="w-9 h-9 rounded-lg border bg-muted flex items-center justify-center overflow-hidden flex-shrink-0">
                        <img src={form.favicon_url} alt="Favicon" className="w-5 h-5 object-contain" />
                      </div>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Tab 4: Email Signature */}
          <TabsContent value="signature" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Mail className="w-5 h-5 text-primary" />
                  Default Email Signature
                </CardTitle>
                <CardDescription>
                  Appended to system-generated notifications, invoices, and staff ticket replies.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Signature Template</label>
                  <Textarea
                    name="email_signature"
                    value={form.email_signature}
                    onChange={handleChange}
                    placeholder="--&#10;Best regards,&#10;Support Team"
                    rows={6}
                    className="font-mono text-sm"
                  />
                </div>

                {form.email_signature && (
                  <div className="space-y-2">
                    <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                      Signature Live Preview
                    </label>
                    <div className="p-4 rounded-lg bg-muted/40 border text-sm whitespace-pre-line text-muted-foreground">
                      {form.email_signature}
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* Action Footer */}
        <div className="flex justify-end pt-4">
          <Button type="submit" disabled={saving} className="gap-2 px-6">
            <Save className="w-4 h-4" />
            {saving ? 'Saving...' : 'Save All Changes'}
          </Button>
        </div>
      </form>
    </div>
  );
};
