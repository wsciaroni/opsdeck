import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import TicketList from './TicketList';
import { BrowserRouter } from 'react-router-dom';
import type { Ticket } from '../../types';

// Mock matchMedia for Mobile view (matches: false for min-width: 768px)
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false, // Force mobile view
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

const mockTickets: Ticket[] = [
  {
    id: '123',
    organization_id: 'org1',
    title: 'Mobile Ticket',
    description: 'Desc',
    location: 'Loc',
    status_id: 'new',
    priority_id: 'high',
    reporter_id: 'rep1',
    assignee_user_id: 'assignee1',
    assignee_name: 'John Doe',
    sensitive: false,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    completed_at: null,
  },
];

describe('TicketList Mobile View', () => {
  it('renders mobile card with compact density', () => {
    render(
      <BrowserRouter>
        <TicketList
          tickets={mockTickets}
          isLoading={false}
          error={null}
          density="compact"
          onOpenNewTicket={() => {}}
        />
      </BrowserRouter>
    );

    // The card itself is a list item containing a button
    // We find the button that wraps the content
    const buttons = screen.getAllByRole('button');
    // Filter to find the card button (it contains the title)
    const cardButton = buttons.find(b => b.textContent?.includes('Mobile Ticket'));

    expect(cardButton).toBeDefined();
    expect(cardButton).toHaveClass('py-2');

    // Check for font size class (text-xs for compact title)
    const title = screen.getByRole('heading', { level: 3 });
    expect(title).toHaveClass('text-xs');
  });

  it('renders mobile card with standard density', () => {
    render(
      <BrowserRouter>
        <TicketList
          tickets={mockTickets}
          isLoading={false}
          error={null}
          density="standard"
          onOpenNewTicket={() => {}}
        />
      </BrowserRouter>
    );

    const buttons = screen.getAllByRole('button');
    const cardButton = buttons.find(b => b.textContent?.includes('Mobile Ticket'));

    expect(cardButton).toBeDefined();
    expect(cardButton).toHaveClass('py-4');

    const title = screen.getByRole('heading', { level: 3 });
    expect(title).toHaveClass('text-sm');
  });

  it('renders mobile card with comfortable density', () => {
    render(
      <BrowserRouter>
        <TicketList
          tickets={mockTickets}
          isLoading={false}
          error={null}
          density="comfortable"
          onOpenNewTicket={() => {}}
        />
      </BrowserRouter>
    );

    const buttons = screen.getAllByRole('button');
    const cardButton = buttons.find(b => b.textContent?.includes('Mobile Ticket'));

    expect(cardButton).toBeDefined();
    expect(cardButton).toHaveClass('py-6');

    const title = screen.getByRole('heading', { level: 3 });
    expect(title).toHaveClass('text-base');
  });
});
