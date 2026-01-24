import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import PublicTicketList from './PublicTicketList';
import { BrowserRouter } from 'react-router-dom';
import type { PublicTicket } from '../../api/public';

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useParams: () => ({ token: 'test-token' }),
    useNavigate: () => vi.fn(),
  };
});

const mockTickets: PublicTicket[] = [
  {
    id: '1',
    organization_id: 'org1',
    title: 'Test Ticket',
    description: 'Desc',
    location: 'Loc',
    status_id: 'new',
    priority_id: 'high',
    reporter_id: 'rep1',
    assignee_user_id: 'assignee1',
    sensitive: false,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    completed_at: null,
  },
];

describe('PublicTicketList', () => {
  it('does not render Assignee column in desktop view', () => {
    render(
      <BrowserRouter>
        <PublicTicketList tickets={mockTickets} isLoading={false} error={null} />
      </BrowserRouter>
    );
    expect(screen.queryByText('Assignee')).not.toBeInTheDocument();
  });
});
