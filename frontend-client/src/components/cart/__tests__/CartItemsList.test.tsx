import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { CartItemsList } from '../CartItemsList';
import { CartItem } from '@/lib/cart';

describe('CartItemsList - Component Tests', () => {
  const mockItems: CartItem[] = [
    {
      id: 'item-1',
      product_id: 101,
      title: 'Cloud Hosting Basic',
      price: 9.99,
      period: '1M',
      type: 'hosting',
    },
    {
      id: 'item-2',
      product_id: 102,
      title: 'Domain .com',
      price: 12.99,
      period: '1Y',
      type: 'domain',
    },
  ];

  it('renders cart items and count', () => {
    render(
      <CartItemsList
        items={mockItems}
        onRemoveItem={vi.fn()}
        onClearCart={vi.fn()}
      />
    );

    expect(screen.getByText('Selected Services (2)')).toBeDefined();
    expect(screen.getByText('Cloud Hosting Basic')).toBeDefined();
    expect(screen.getByText('Domain .com')).toBeDefined();
    expect(screen.getByText('$9.99')).toBeDefined();
    expect(screen.getByText('$12.99')).toBeDefined();
  });

  it('triggers onClearCart when Clear Cart button is clicked', () => {
    const clearCartMock = vi.fn();
    render(
      <CartItemsList
        items={mockItems}
        onRemoveItem={vi.fn()}
        onClearCart={clearCartMock}
      />
    );

    const clearBtn = screen.getByText('Clear Cart');
    fireEvent.click(clearBtn);
    expect(clearCartMock).toHaveBeenCalledTimes(1);
  });
});
