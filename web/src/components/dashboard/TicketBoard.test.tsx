import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import TicketBoard from './TicketBoard';
import { BrowserRouter } from 'react-router-dom';
import type { Ticket } from '../../types';

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => vi.fn(),
  };
});

const mockTickets: Ticket[] = [
  {
    id: '1',
    organization_id: 'org1',
    title: 'Test Ticket 1',
    description: 'Desc',
    location: 'Loc',
    status_id: 'new',
    priority_id: 'high',
    reporter_id: 'rep1',
    assignee_user_id: 'assignee1',
    assignee_name: 'Assignee One',
    sensitive: false,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    completed_at: null,
  },
  {
    id: '2',
    organization_id: 'org1',
    title: 'Test Ticket 2',
    description: 'Desc',
    location: 'Loc',
    status_id: 'in_progress',
    priority_id: 'medium',
    reporter_id: 'rep1',
    assignee_user_id: 'assignee2',
    sensitive: false,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    completed_at: null,
  }
];

describe('TicketBoard', () => {
  it('renders tickets in correct columns', () => {
    render(
      <BrowserRouter>
        <TicketBoard
            tickets={mockTickets}
            isLoading={false}
            error={null}
            density="standard"
            onOpenNewTicket={() => {}}
        />
      </BrowserRouter>
    );

    expect(screen.getByText('Test Ticket 1')).toBeInTheDocument();
    expect(screen.getByText('Test Ticket 2')).toBeInTheDocument();
    expect(screen.getByText('Assignee One')).toBeInTheDocument();
  });

  it('renders loading state', () => {
    render(
      <BrowserRouter>
        <TicketBoard
            tickets={[]}
            isLoading={true}
            error={null}
            density="standard"
            onOpenNewTicket={() => {}}
        />
      </BrowserRouter>
    );
    expect(screen.getByText('Loading tickets...')).toBeInTheDocument();
  });

  it('renders only active columns by default', () => {
    render(
      <BrowserRouter>
        <TicketBoard
            tickets={[]}
            isLoading={false}
            error={null}
            density="standard"
            onOpenNewTicket={() => {}}
        />
      </BrowserRouter>
    );

    expect(screen.getByText('New')).toBeInTheDocument();
    expect(screen.getByText('In Progress')).toBeInTheDocument();
    expect(screen.getByText('On Hold')).toBeInTheDocument();
    expect(screen.queryByText('Done')).not.toBeInTheDocument();
    expect(screen.queryByText('Canceled')).not.toBeInTheDocument();
  });

  it('renders filtered columns when visibleStatuses is provided', () => {
    render(
      <BrowserRouter>
        <TicketBoard
            tickets={[]}
            isLoading={false}
            error={null}
            density="standard"
            visibleStatuses={['new', 'done']}
            onOpenNewTicket={() => {}}
        />
      </BrowserRouter>
    );

    expect(screen.getByText('New')).toBeInTheDocument();
    expect(screen.getByText('Done')).toBeInTheDocument();
    expect(screen.queryByText('In Progress')).not.toBeInTheDocument();
  });

  it('renders semantic HTML structure', () => {
    render(
      <BrowserRouter>
        <TicketBoard
            tickets={mockTickets}
            isLoading={false}
            error={null}
            density="standard"
            visibleStatuses={['new']}
            onOpenNewTicket={() => {}}
        />
      </BrowserRouter>
    );

    // Check for Column Header (h2)
    const columnHeader = screen.getByRole('heading', { level: 2, name: /New/i });
    expect(columnHeader).toBeInTheDocument();

    // Check for Ticket List (ul)
    const list = screen.getByRole('list');
    expect(list).toBeInTheDocument();

    // Check for Ticket Item (li)
    const listItems = screen.getAllByRole('listitem');
    expect(listItems.length).toBeGreaterThan(0);

    // Check for Ticket Title (h3)
    const ticketTitle = screen.getByRole('heading', { level: 3, name: 'Test Ticket 1' });
    expect(ticketTitle).toBeInTheDocument();

    // Check for Ticket Link (a)
    const link = screen.getByRole('link', { name: /Test Ticket 1/i });
    expect(link).toHaveAttribute('href', '/tickets/1');
  });
});
