import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import TicketList from './TicketList';
import { BrowserRouter } from 'react-router-dom';
import type { Ticket } from '../../types';

const mockNavigate = vi.fn();

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

const mockTickets: Ticket[] = [
  {
    id: '123',
    organization_id: 'org1',
    title: 'Test Ticket',
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

describe('TicketList', () => {
  it('renders ticket list correctly', () => {
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
    // Use getAllByText because it renders in both mobile and desktop views
    expect(screen.getAllByText('Test Ticket').length).toBeGreaterThan(0);
    expect(screen.getAllByText('John Doe').length).toBeGreaterThan(0);
  });

  it('navigates on row click', () => {
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

    // Find the text specifically inside the table row (td)
    // We can look for the row directly if we can identify it, but text is easier.
    const titleCells = screen.getAllByText('Test Ticket');
    // Find the one that is inside a td (desktop view) or we can just pick the second one if we know order?
    // Better: find closest tr.
    const desktopTitle = titleCells.find(el => el.closest('tr'));

    expect(desktopTitle).toBeInTheDocument();

    if (desktopTitle) {
        const row = desktopTitle.closest('tr');
        expect(row).toBeInTheDocument();
        if (row) {
            fireEvent.click(row);
            expect(mockNavigate).toHaveBeenCalledWith('/tickets/123');
        }
    }
  });

  it('navigates on row Enter key press', () => {
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

    const titleCells = screen.getAllByText('Test Ticket');
    const desktopTitle = titleCells.find(el => el.closest('tr'));
    expect(desktopTitle).toBeInTheDocument();

    if (desktopTitle) {
        const row = desktopTitle.closest('tr');
        if (row) {
            fireEvent.keyDown(row, { key: 'Enter', code: 'Enter' });
            expect(mockNavigate).toHaveBeenCalledWith('/tickets/123');
        }
    }
  });
});
