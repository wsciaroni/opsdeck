import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import PublicTicketDetail from './PublicTicketDetail';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useParams: () => ({ token: 'test-token', ticketId: '1' }),
  };
});

// Mock API
vi.mock('../api/public', () => ({
  getPublicTicket: vi.fn().mockResolvedValue({
    id: '1',
    organization_id: 'org1',
    title: 'Test Ticket',
    status_id: 'new',
    priority_id: 'high',
    created_at: new Date().toISOString(),
    description: 'Desc',
    reporter_id: 'rep1',
    assignee_user_id: 'assignee1',
    sensitive: false,
    location: 'Loc',
    updated_at: new Date().toISOString(),
    completed_at: null,
  }),
}));

// Mock PublicTicketComments to avoid testing its internals
vi.mock('../components/PublicTicketComments', () => ({
  default: () => <div>Comments Section</div>,
}));

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
    },
});

describe('PublicTicketDetail', () => {
  it('does not render Reporter or Assignee details', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
            <PublicTicketDetail />
        </BrowserRouter>
      </QueryClientProvider>
    );

    // Wait for loading to finish (by finding title)
    await screen.findByText('Test Ticket');

    expect(screen.queryByText('Reporter')).not.toBeInTheDocument();
    expect(screen.queryByText('Assignee')).not.toBeInTheDocument();
  });
});
