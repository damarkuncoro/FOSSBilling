import { useEffect, useState, useRef } from 'react';
import { useToast } from './use-toast';

export function useLiveNotifications() {
  const { toast } = useToast();
  const [isConnected, setIsConnected] = useState(false);
  const toastRef = useRef(toast);
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<any>(null);

  // Keep toast updated but stable for the effect
  useEffect(() => {
    toastRef.current = toast;
  }, [toast]);

  useEffect(() => {
    let isMounted = true;

    const connect = () => {
      // If already connected or connecting, don't start another one
      if (socketRef.current?.readyState === WebSocket.OPEN ||
          socketRef.current?.readyState === WebSocket.CONNECTING) {
        return;
      }

      const token = localStorage.getItem('fossbilling_admin_token');
      if (!token || !isMounted) return;

      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host === 'localhost:3000' ? 'localhost:8080' : window.location.host;
      const wsUrl = `${protocol}//${host}/api/v1/admin/system/ws?token=${token}`;

      console.log('🔌 Attempting WebSocket connection...');
      const socket = new WebSocket(wsUrl);
      socketRef.current = socket;

      socket.onopen = () => {
        if (isMounted) {
          console.log('✅ Live notification stream connected');
          setIsConnected(true);
        }
      };

      socket.onmessage = (event) => {
        if (!isMounted) return;
        try {
          const data = JSON.parse(event.data);
          if (data.title) {
            toastRef.current({
              title: data.title,
              description: data.message,
              variant: data.priority === 'high' || data.priority === 'urgent' ? 'destructive' : 'default',
            });
          }
        } catch (err) {
          console.error('Failed to parse WS message', err);
        }
      };

      socket.onclose = (event) => {
        socketRef.current = null;
        if (isMounted) {
          setIsConnected(false);
          // Only reconnect if it wasn't a clean close from our end
          if (event.code !== 1000) {
            console.log(`❌ WS disconnected (code: ${event.code}). Retrying in 5s...`);
            reconnectTimeoutRef.current = setTimeout(connect, 5000);
          }
        }
      };

      socket.onerror = (err) => {
        console.error('WebSocket Error:', err);
      };
    };

    connect();

    return () => {
      isMounted = false;
      if (socketRef.current) {
        console.log('🧹 Cleaning up WebSocket connection');
        socketRef.current.close(1000);
        socketRef.current = null;
      }
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
    };
  }, []);

  return { isConnected };
}
