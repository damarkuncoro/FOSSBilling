import { useState } from 'react';

export interface ToastProps {
  title?: string;
  description?: string;
  variant?: 'default' | 'destructive';
}

export function useToast() {
  const [toasts, setToasts] = useState<ToastProps[]>([]);

  const toast = ({ title, description, variant = 'default' }: ToastProps) => {
    console.log(`[Toast] ${title}: ${description} (${variant})`);
    // Simple implementation for now
    setToasts((prev) => [...prev, { title, description, variant }]);

    // Auto-dismiss after 3s
    setTimeout(() => {
      setToasts((prev) => prev.slice(1));
    }, 3000);
  };

  return { toast, toasts };
}
