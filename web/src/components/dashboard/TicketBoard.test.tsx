import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import TicketBoard from './TicketBoard';
import { BrowserRouter } from 'react-router-dom';
import type { Ticket } from '../../types';

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

  it('renders tickets as accessible links', () => {
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

    const link1 = screen.getByRole('link', { name: /Test Ticket 1/i });
    expect(link1).toHaveAttribute('href', '/tickets/1');
    expect(link1).toHaveAttribute('aria-label', expect.stringContaining('Test Ticket 1'));
    expect(link1).toHaveAttribute('aria-label', expect.stringContaining('high priority'));

    const link2 = screen.getByRole('link', { name: /Test Ticket 2/i });
    expect(link2).toHaveAttribute('href', '/tickets/2');
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
});
