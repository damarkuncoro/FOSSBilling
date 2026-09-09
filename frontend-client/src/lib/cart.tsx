import React, { createContext, useContext, useEffect, useState } from 'react';
import { api } from './api';

export interface CartItem {
  id: string;
  product_id: number;
  title: string;
  type: string; // 'hosting' | 'vps' | 'domain' | 'license' | 'downloadable'
  price: number;
  period: string; // '1M' | '3M' | '1Y' | 'ONETIME'
  domain_name?: string;
  form_id?: number;
  config?: Record<string, any>;
}

interface CartContextType {
  items: CartItem[];
  promoCode: string;
  discount: number;
  subtotal: number;
  tax: number;
  total: number;
  addItem: (item: CartItem) => void;
  removeItem: (id: string) => void;
  clearCart: () => void;
  updateItemConfig: (id: string, config: Record<string, any>) => void;
  setPromoCode: (code: string) => void;
  applyPromo: (code: string) => Promise<void>;
  calculateTotals: () => Promise<void>;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [items, setItems] = useState<CartItem[]>(() => {
    try {
      const saved = localStorage.getItem('fossbilling_cart_items');
      if (!saved) return [];
      const parsed = JSON.parse(saved);
      return Array.isArray(parsed) ? parsed : [];
    } catch {
      return [];
    }
  });
  const [promoCode, setPromoCodeState] = useState<string>('');
  const [discount, setDiscount] = useState<number>(0);
  const [subtotal, setSubtotal] = useState<number>(0);
  const [tax, setTax] = useState<number>(0);
  const [total, setTotal] = useState<number>(0);

  useEffect(() => {
    localStorage.setItem('fossbilling_cart_items', JSON.stringify(items));
    calculateTotals();
  }, [items]);

  const calculateTotals = async () => {
    if (items.length === 0) {
      setSubtotal(0);
      setDiscount(0);
      setTax(0);
      setTotal(0);
      return;
    }

    try {
      const calc = await api.calculateCart(
        items.map((i) => ({
          product_id: i.product_id,
          title: i.title,
          period: i.period,
          price: Number(i.price),
          quantity: 1,
          config: i.domain_name ? { domain_name: i.domain_name } : undefined,
        })),
        promoCode || undefined
      );
      setSubtotal(Number(calc.subtotal) || 0);
      setDiscount(Number(calc.discount) || 0);
      setTax(Number(calc.tax) || 0);
      setTotal(Number(calc.total) || 0);
    } catch {
      // Fallback local calc
      const rawSub = items.reduce((acc, curr) => acc + (Number(curr.price) || 0), 0);
      const rawTax = rawSub * 0.11;
      setSubtotal(rawSub);
      setTax(rawTax);
      setTotal(rawSub + rawTax);
    }
  };

  const addItem = (item: CartItem) => {
    setItems((prev) => [...prev, item]);
  };

  const removeItem = (id: string) => {
    setItems((prev) => prev.filter((i) => i.id !== id));
  };

  const updateItemConfig = (id: string, config: Record<string, any>) => {
    setItems((prev) => prev.map((i) => (i.id === id ? { ...i, config } : i)));
  };

  const clearCart = () => {
    setItems([]);
    setPromoCodeState('');
    setDiscount(0);
  };

  const applyPromo = async (code: string) => {
    setPromoCodeState(code);
    try {
      const calc = await api.calculateCart(
        items.map((i) => ({
          product_id: i.product_id,
          title: i.title,
          period: i.period,
          price: Number(i.price),
          quantity: 1,
          config: i.domain_name ? { domain_name: i.domain_name } : undefined,
        })),
        code
      );
      setSubtotal(Number(calc.subtotal) || 0);
      setDiscount(Number(calc.discount) || 0);
      setTax(Number(calc.tax) || 0);
      setTotal(Number(calc.total) || 0);
    } catch (err: any) {
      throw new Error(err.message || 'Invalid promotional voucher code');
    }
  };

  return (
    <CartContext.Provider
      value={{
        items,
        promoCode,
        discount,
        subtotal,
        tax,
        total,
        addItem,
        removeItem,
        clearCart,
        updateItemConfig,
        setPromoCode: setPromoCodeState,
        applyPromo,
        calculateTotals,
      }}
    >
      {children}
    </CartContext.Provider>
  );
};

export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};
