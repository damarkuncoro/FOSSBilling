import React, { useState } from 'react';
import { X, Radio } from 'lucide-react';
import { AVAILABLE_EVENTS } from '../../hooks/useWebhooks';

interface AddWebhookDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: { name: string; url: string; events: string[] }) => void;
}

export const AddWebhookDialog: React.FC<AddWebhookDialogProps> = ({
  isOpen,
  onClose,
  onSubmit,
}) => {
  const [name, setName] = useState('');
  const [url, setUrl] = useState('');
  const [selectedEvents, setSelectedEvents] = useState<string[]>(['order.activated', 'invoice.paid']);

  if (!isOpen) return null;

  const toggleEvent = (eventId: string) => {
    setSelectedEvents((prev) =>
      prev.includes(eventId) ? prev.filter((e) => e !== eventId) : [...prev, eventId]
    );
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (name.trim() && url.trim() && selectedEvents.length > 0) {
      onSubmit({ name: name.trim(), url: url.trim(), events: selectedEvents });
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl animate-in fade-in zoom-in duration-200">
        <div className="flex items-center justify-between pb-4 border-b border-gray-100">
          <div className="flex items-center gap-2 text-indigo-600 font-semibold">
            <Radio className="w-5 h-5" />
            <h3 className="text-gray-900 font-bold">Add Webhook Endpoint</h3>
          </div>
          <button onClick={onClose} className="p-1 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100">
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Endpoint Name</label>
            <input
              type="text"
              required
              placeholder="e.g. Discord Billing Bot"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3.5 py-2 border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-1">Payload URL</label>
            <input
              type="url"
              required
              placeholder="https://your-server.com/api/webhook"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              className="w-full px-3.5 py-2 font-mono border border-gray-200 rounded-lg text-xs focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-700 uppercase mb-2">Events to Subscribe</label>
            <div className="grid grid-cols-2 gap-2 max-h-48 overflow-y-auto p-1">
              {AVAILABLE_EVENTS.map((ev) => (
                <label
                  key={ev.id}
                  className={`flex items-center gap-2 p-2 rounded-lg border text-xs cursor-pointer transition-all ${
                    selectedEvents.includes(ev.id)
                      ? 'bg-indigo-50/70 border-indigo-200 text-indigo-900 font-medium'
                      : 'bg-white border-gray-200 text-gray-600 hover:bg-gray-50'
                  }`}
                >
                  <input
                    type="checkbox"
                    checked={selectedEvents.includes(ev.id)}
                    onChange={() => toggleEvent(ev.id)}
                    className="rounded text-indigo-600 border-gray-300"
                  />
                  <span>{ev.label}</span>
                </label>
              ))}
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-4 border-t border-gray-100">
            <button type="button" onClick={onClose} className="px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 rounded-lg">
              Cancel
            </button>
            <button type="submit" className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-sm">
              Save Webhook
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
