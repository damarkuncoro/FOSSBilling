import { useEffect, useState } from 'react';
import { useToast } from './use-toast';

export function useLiveNotifications() {
  const { toast } = useToast();
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('fossbilling_admin_token');
    if (!token) return;

    // Determine WS protocol (ws or wss)
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host === 'localhost:3000' ? 'localhost:8080' : window.location.host;
    const wsUrl = `${protocol}//${host}/api/v1/admin/system/ws?token=${token}`;

    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log('✅ Live notification stream connected');
      setIsConnected(true);
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.title) {
          toast({
            title: data.title,
            description: data.message,
            variant: data.priority === 'high' || data.priority === 'urgent' ? 'destructive' : 'default',
          });
        }
      } catch (err) {
        console.error('Failed to parse WS message', err);
      }
    };

    socket.onclose = () => {
      console.log('❌ Live notification stream disconnected');
      setIsConnected(false);
    };

    return () => {
      socket.close();
    };
  }, [toast]);

  return { isConnected };
}
