import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { NewsAnnouncements } from '../NewsAnnouncements';

describe('NewsAnnouncements - Component Tests', () => {
  it('returns null when news list is empty', () => {
    const { container } = render(<NewsAnnouncements news={[]} />);
    expect(container.firstChild).toBeNull();
  });

  it('renders news announcements when data is provided', () => {
    const mockNews = [
      {
        id: 1,
        title: 'Upgrade Jaringan 2026',
        content: 'Kami meningkatkan kapasitas bandwidth ke 10Gbps unmetered.',
        created_at: '2026-01-15T00:00:00Z',
      },
    ];

    render(<NewsAnnouncements news={mockNews} />);

    expect(screen.getByText('System Announcements & News')).toBeDefined();
    expect(screen.getByText('Upgrade Jaringan 2026')).toBeDefined();
    expect(
      screen.getByText('Kami meningkatkan kapasitas bandwidth ke 10Gbps unmetered.')
    ).toBeDefined();
  });
});
