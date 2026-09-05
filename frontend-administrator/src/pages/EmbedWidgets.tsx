import React from 'react';
import { Code2 } from 'lucide-react';
import { useEmbedWidgets } from '../hooks/useEmbedWidgets';
import { WidgetGeneratorCard } from '../components/embedwidgets/WidgetGeneratorCard';
import { WidgetCodePreview } from '../components/embedwidgets/WidgetCodePreview';

export const EmbedWidgets: React.FC = () => {
  const {
    config,
    copied,
    updateConfig,
    generateEmbedCode,
    copyToClipboard,
  } = useEmbedWidgets();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
          <Code2 className="w-7 h-7 text-indigo-600" /> Embed & Order Buttons
        </h1>
        <p className="text-sm text-gray-500 mt-1">
          Generate embeddable order buttons, pricing cards, and checkout widgets for external landing pages.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <WidgetGeneratorCard config={config} onChange={updateConfig} />
        <WidgetCodePreview
          config={config}
          embedCode={generateEmbedCode()}
          copied={copied}
          onCopy={copyToClipboard}
        />
      </div>
    </div>
  );
};
