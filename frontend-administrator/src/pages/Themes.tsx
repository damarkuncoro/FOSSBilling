import React from 'react';
import { Palette, RefreshCw, CheckCircle2, Settings2 } from 'lucide-react';
import { useThemes } from '@/hooks/useThemes';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

export const Themes: React.FC = () => {
  const {
    clientThemes,
    adminThemes,
    currentClientTheme,
    currentAdminTheme,
    loading,
    fetchThemes,
    handleSelectTheme,
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
        <Button variant="outline" size="sm" onClick={fetchThemes} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      <Tabs defaultValue="client" className="space-y-4">
        <TabsList>
          <TabsTrigger value="client">Client Portal Themes</TabsTrigger>
          <TabsTrigger value="admin">Admin Dashboard Themes</TabsTrigger>
        </TabsList>

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
