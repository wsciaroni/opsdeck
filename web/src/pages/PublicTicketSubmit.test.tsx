import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import PublicTicketSubmit from './PublicTicketSubmit';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter, Routes, Route } from 'react-router-dom';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: false },
    mutations: { retry: false },
  },
});

vi.mock('../api/tickets', () => ({
  createPublicTicket: vi.fn(),
}));

vi.mock('../components/FileUpload', () => ({
  default: () => <div data-testid="file-upload">File Upload</div>,
}));

describe('PublicTicketSubmit', () => {
  it('renders form when token is present', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={['/submit?token=valid-token']}>
          <Routes>
            <Route path="/submit" element={<PublicTicketSubmit />} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText(/Submit a Ticket/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Email/i)).toBeInTheDocument();

    // Check that inputs are required
    expect(screen.getByLabelText(/Name/i)).toBeRequired();
    expect(screen.getByLabelText(/Email/i)).toBeRequired();
    expect(screen.getByLabelText(/Title/i)).toBeRequired();
    expect(screen.getByLabelText(/Description/i)).toBeRequired();
  });

  it('renders error when token is missing', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={['/submit']}>
           <Routes>
            <Route path="/submit" element={<PublicTicketSubmit />} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText(/Invalid Link/i)).toBeInTheDocument();
  });
});
