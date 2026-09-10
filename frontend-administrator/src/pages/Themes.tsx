import React from 'react';
import { Palette, RefreshCw, CheckCircle2, Settings2 } from 'lucide-react';
import { useThemes } from '@/hooks/useThemes';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export const Themes: React.FC = () => {
  const {
    clientThemes,
    adminThemes,
    currentClientTheme,
    currentAdminTheme,
    branding,
    setBranding,
    loading,
    savingBranding,
    fetchThemes,
    handleSelectTheme,
    handleUpdateBranding,
  } = useThemes();

  const ThemeGrid = ({ themes, currentCode, target }: { themes: any[], currentCode: string, target: 'client' | 'admin' }) => (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {themes.map((theme) => {
        const isActive = theme.code === currentCode;
        return (
          <Card key={theme.code} className={`border-border/60 shadow-sm flex flex-col justify-between transition-all ${isActive ? 'ring-2 ring-primary ring-offset-2' : ''}`}>
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <Badge variant="outline" className="text-[10px] uppercase font-mono">
                  {theme.code}
                </Badge>
                {isActive && (
                  <Badge variant="success" className="gap-1">
                    <CheckCircle2 className="h-3 w-3" /> Active
                  </Badge>
                )}
              </div>
              <CardTitle className="text-base mt-2">{theme.name}</CardTitle>
              <CardDescription className="text-xs line-clamp-2">{theme.description}</CardDescription>
            </CardHeader>
            <CardContent>
               <div className="aspect-video rounded-lg bg-muted flex items-center justify-center border overflow-hidden">
                  {theme.screenshot ? (
                    <img src={theme.screenshot} alt={theme.name} className="w-full h-full object-cover" />
                  ) : (
                    <Palette className="h-10 w-10 text-muted-foreground/20" />
                  )}
               </div>
               <div className="flex items-center justify-between text-[10px] text-muted-foreground mt-4">
                  <span>v{theme.version}</span>
                  <span>By {theme.author}</span>
               </div>
            </CardContent>
            <CardFooter className="gap-2">
              {!isActive && (
                <Button
                  variant="default"
                  size="sm"
                  className="w-full text-xs font-semibold"
                  onClick={() => handleSelectTheme(theme.code, target)}
                >
                  Activate Theme
                </Button>
              )}
              <Button
                variant="outline"
                size="sm"
                className={`text-xs font-semibold gap-1.5 ${isActive ? 'w-full' : 'w-auto'}`}
              >
                <Settings2 className="h-3.5 w-3.5" />
                Configure
              </Button>
            </CardFooter>
          </Card>
        );
      })}
    </div>
  );

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2.5">
            <Palette className="w-7 h-7 text-primary" /> Appearance & Theme Manager
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            Customize the visual identity of your client portal and administrator dashboard.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={() => fetchThemes()} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      <Tabs defaultValue="client" className="space-y-4">
        <TabsList>
          <TabsTrigger value="client">Client Portal Themes</TabsTrigger>
          <TabsTrigger value="admin">Admin Dashboard Themes</TabsTrigger>
          <TabsTrigger value="branding">Global Branding</TabsTrigger>
        </TabsList>

        <TabsContent value="branding">
           <Card className="border-border/60">
             <CardHeader>
               <CardTitle className="text-base font-semibold">White-label & Brand Identity</CardTitle>
               <CardDescription>Customize how your customers see your brand across all portals and emails.</CardDescription>
             </CardHeader>
             <CardContent>
                <form onSubmit={handleUpdateBranding} className="max-w-2xl space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="space-y-2">
                      <Label htmlFor="company_name">Company Display Name</Label>
                      <Input
                        id="company_name"
                        value={branding.company_name}
                        onChange={(e) => setBranding({ ...branding, company_name: e.target.value })}
                        placeholder="e.g. My Hosting Co."
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="primary_color">Primary Brand Color</Label>
                      <div className="flex gap-2">
                        <Input
                          id="primary_color"
                          type="color"
                          className="w-12 p-1 h-10"
                          value={branding.primary_color}
                          onChange={(e) => setBranding({ ...branding, primary_color: e.target.value })}
                        />
                        <Input
                          value={branding.primary_color}
                          onChange={(e) => setBranding({ ...branding, primary_color: e.target.value })}
                          className="font-mono"
                        />
                      </div>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="logo_url">Custom Logo URL (Horizontal)</Label>
                    <Input
                      id="logo_url"
                      value={branding.logo_url || ''}
                      onChange={(e) => setBranding({ ...branding, logo_url: e.target.value })}
                      placeholder="https://your-domain.com/logo.png"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="favicon_url">Favicon URL (16x16 or 32x32)</Label>
                    <Input
                      id="favicon_url"
                      value={branding.favicon_url || ''}
                      onChange={(e) => setBranding({ ...branding, favicon_url: e.target.value })}
                      placeholder="https://your-domain.com/favicon.ico"
                    />
                  </div>

                  <Button type="submit" disabled={savingBranding} className="font-semibold">
                    {savingBranding ? 'Saving...' : 'Save Branding Identity'}
                  </Button>
                </form>
             </CardContent>
           </Card>
        </TabsContent>

        <TabsContent value="client">
          {loading ? (
            <div className="py-20 text-center text-muted-foreground">Loading themes...</div>
          ) : (
            <ThemeGrid themes={clientThemes} currentCode={currentClientTheme?.code} target="client" />
          )}
        </TabsContent>

        <TabsContent value="admin">
          {loading ? (
            <div className="py-20 text-center text-muted-foreground">Loading themes...</div>
          ) : (
            <ThemeGrid themes={adminThemes} currentCode={currentAdminTheme?.code} target="admin" />
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};
