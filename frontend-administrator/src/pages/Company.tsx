import React, { useState } from 'react';
import {
  Building2,
  RefreshCw,
  CheckCircle,
  AlertCircle,
  Globe,
  Image as ImageIcon,
} from 'lucide-react';
import { useCompany } from '@/hooks/useCompany';
import { CompanyGeneralTab } from '@/components/company/CompanyGeneralTab';
import { CompanyBrandingTab } from '@/components/company/CompanyBrandingTab';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

export const Company: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'general' | 'branding'>('general');
  const {
    form,
    loading,
    saving,
    successMessage,
    errorMessage,
    fetchCompanySettings,
    handleChange,
    handleSave,
  } = useCompany();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300 max-w-4xl mx-auto">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
            <Building2 className="h-6 w-6 text-primary" />
            Company & Business Profile
          </h1>
          <p className="text-sm text-muted-foreground">
            Configure global business information, invoice headers, tax identification, and branding assets.
          </p>
        </div>
        <div className="flex items-center gap-2">
          {form.updated_at && (
            <Badge variant="outline" className="text-xs text-muted-foreground">
              Updated: {new Date(form.updated_at).toLocaleDateString()}
            </Badge>
          )}
          <Button variant="outline" size="sm" onClick={fetchCompanySettings} disabled={loading} className="gap-2">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {successMessage && (
        <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 flex items-center gap-3 text-sm">
          <CheckCircle className="h-5 w-5 shrink-0" />
          <p>{successMessage}</p>
        </div>
      )}

      {errorMessage && (
        <div className="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive flex items-center gap-3 text-sm">
          <AlertCircle className="h-5 w-5 shrink-0" />
          <p>{errorMessage}</p>
        </div>
      )}

      <form onSubmit={handleSave}>
        <Card className="border-border/60 shadow-sm">
          <CardHeader className="border-b pb-4">
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-base font-semibold">Settings & Configuration</CardTitle>
                <CardDescription>Shown on invoices, client portal headers, and transactional emails</CardDescription>
              </div>
              <div className="flex bg-muted/60 p-1 rounded-lg border">
                <button
                  type="button"
                  onClick={() => setActiveTab('general')}
                  className={`px-3 py-1 text-xs font-semibold rounded-md transition-all flex items-center gap-1.5 ${
                    activeTab === 'general' ? 'bg-card text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  <Globe className="h-3.5 w-3.5" />
                  General & Address
                </button>
                <button
                  type="button"
                  onClick={() => setActiveTab('branding')}
                  className={`px-3 py-1 text-xs font-semibold rounded-md transition-all flex items-center gap-1.5 ${
                    activeTab === 'branding' ? 'bg-card text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  <ImageIcon className="h-3.5 w-3.5" />
                  Logos & Branding
                </button>
              </div>
            </div>
          </CardHeader>

          <CardContent className="pt-6">
            {activeTab === 'general' ? (
              <CompanyGeneralTab form={form} onChange={handleChange} />
            ) : (
              <CompanyBrandingTab form={form} onChange={handleChange} />
            )}
          </CardContent>

          <CardFooter className="flex justify-between border-t bg-muted/20 px-6 py-4">
            <span className="text-xs text-muted-foreground">Changes reflect in real-time across invoices</span>
            <Button type="submit" disabled={saving}>
              {saving ? 'Saving...' : 'Save Company Details'}
            </Button>
          </CardFooter>
        </Card>
      </form>
    </div>
  );
};

export default Company;
