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
});
