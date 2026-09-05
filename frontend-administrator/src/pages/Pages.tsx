import React from 'react';
import { RefreshCw } from 'lucide-react';
import { usePages } from '@/hooks/usePages';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { PagesListTab } from '@/components/pages/PagesListTab';
import { KnowledgebaseListTab } from '@/components/pages/KnowledgebaseListTab';
import { EditPageDialog } from '@/components/pages/EditPageDialog';

export const Pages: React.FC = () => {
  const {
    pages,
    articles,
    loading,
    openPageModal,
    setOpenPageModal,
    pageForm,
    setPageForm,
    fetchData,
    handleSavePage,
    handleDeletePage,
  } = usePages();

  return (
    <div className="space-y-6 animate-in fade-in-50 duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Custom Pages & Knowledgebase CMS</h1>
          <p className="text-sm text-muted-foreground">
            Manage legal documents, static content pages, and customer self-help articles.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchData} disabled={loading} className="gap-2">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      <Tabs defaultValue="pages" className="space-y-4">
        <TabsList>
          <TabsTrigger value="pages">Custom Static Pages ({pages.length})</TabsTrigger>
          <TabsTrigger value="kb">Knowledgebase & FAQs ({articles.length})</TabsTrigger>
        </TabsList>

        <TabsContent value="pages">
          <Card className="border-border/60 shadow-sm">
            <CardHeader>
              <CardTitle className="text-base font-semibold">Published Site Pages</CardTitle>
              <CardDescription>Rendered on the client portal footer and navigation</CardDescription>
            </CardHeader>
            <CardContent>
              <PagesListTab
                pages={pages}
                loading={loading}
                onEdit={(p) => {
                  setPageForm(p);
                  setOpenPageModal(true);
                }}
                onDelete={handleDeletePage}
                onNew={() => {
                  setPageForm({ title: '', slug: '', content: '', published: true });
                  setOpenPageModal(true);
                }}
              />
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="kb">
          <Card className="border-border/60 shadow-sm">
            <CardHeader>
              <CardTitle className="text-base font-semibold">Knowledgebase Articles</CardTitle>
              <CardDescription>Troubleshooting tutorials and onboarding guides</CardDescription>
            </CardHeader>
            <CardContent>
              <KnowledgebaseListTab articles={articles} loading={loading} />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <EditPageDialog
        open={openPageModal}
        onOpenChange={setOpenPageModal}
        form={pageForm}
        onChange={setPageForm}
        onSubmit={handleSavePage}
      />
    </div>
  );
};
