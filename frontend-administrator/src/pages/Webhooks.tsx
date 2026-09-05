import React from 'react';
import { Radio, Plus } from 'lucide-react';
import { useWebhooks } from '../hooks/useWebhooks';
import { WebhookTable } from '../components/webhooks/WebhookTable';
import { AddWebhookDialog } from '../components/webhooks/AddWebhookDialog';
import { WebhookDeliveryLogs } from '../components/webhooks/WebhookDeliveryLogs';

export const Webhooks: React.FC = () => {
  const {
    webhooks,
    logs,
    isAddModal,
    setIsAddModal,
    testingWebhookId,
    createWebhook,
    toggleWebhook,
    deleteWebhook,
    triggerTestPayload,
  } = useWebhooks();

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
            <Radio className="w-7 h-7 text-indigo-600" /> Event Webhooks
          </h1>
          <p className="text-sm text-gray-500 mt-1">
            Dispatch real-time HTTP POST notifications to external services (Discord, Slack, ERP) on system events.
          </p>
        </div>
        <button
          onClick={() => setIsAddModal(true)}
          className="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-sm rounded-xl shadow-sm transition-all"
        >
          <Plus className="w-4 h-4" /> Add Webhook Endpoint
        </button>
      </div>

      <WebhookTable
        webhooks={webhooks}
        testingId={testingWebhookId}
        onToggle={toggleWebhook}
        onDelete={deleteWebhook}
        onTest={triggerTestPayload}
      />

      <WebhookDeliveryLogs logs={logs} />

      <AddWebhookDialog
        isOpen={isAddModal}
        onClose={() => setIsAddModal(false)}
        onSubmit={createWebhook}
      />
    </div>
  );
};
