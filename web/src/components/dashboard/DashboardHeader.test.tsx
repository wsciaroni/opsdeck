import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import DashboardHeader from './DashboardHeader';
import { BrowserRouter } from 'react-router-dom';
import type { Organization } from '../../types';

// Mock dependencies
vi.mock('./FilterPopover', () => ({
  default: () => <div data-testid="filter-popover">FilterPopover</div>
}));

describe('DashboardHeader', () => {
  const defaultProps = {
    currentOrg: {
      id: 'org-1',
      name: 'Test Org',
      slug: 'test-org',
      role: 'admin',
      created_at: '2023-01-01T00:00:00Z',
      updated_at: '2023-01-01T00:00:00Z',
    } as Organization,
    onOpenNewTicket: vi.fn(),
    viewMode: 'list' as const,
    setViewMode: vi.fn(),
    density: 'standard' as const,
    setDensity: vi.fn(),
    onSearch: vi.fn(),
    priority: undefined,
    setPriority: vi.fn(),
    status: undefined,
    setStatus: vi.fn(),
    sortBy: 'created_at',
    setSortBy: vi.fn(),
    sortOrder: 'desc' as const,
    setSortOrder: vi.fn(),
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('debounces search input', async () => {
    vi.useFakeTimers();
    render(
      <BrowserRouter>
        <DashboardHeader {...defaultProps} />
      </BrowserRouter>
    );

    const input = screen.getByPlaceholderText('Search tickets...');
    fireEvent.change(input, { target: { value: 'test query' } });

    // Should not call immediately
    expect(defaultProps.onSearch).not.toHaveBeenCalled();

    // Advance timer by 300ms
    vi.advanceTimersByTime(300);

    expect(defaultProps.onSearch).toHaveBeenCalledWith('test query');

    vi.useRealTimers();
  });

  it('clears search input', async () => {
     vi.useFakeTimers();
    render(
      <BrowserRouter>
        <DashboardHeader {...defaultProps} />
      </BrowserRouter>
    );

    const input = screen.getByPlaceholderText('Search tickets...');
    fireEvent.change(input, { target: { value: 'test' } });

    vi.advanceTimersByTime(300);
    expect(defaultProps.onSearch).toHaveBeenCalledWith('test');

    // Click clear button
    const clearButton = screen.getByLabelText('Clear search');
    fireEvent.click(clearButton);

    expect(input).toHaveValue('');

    // Should call onSearch with empty string after debounce
    vi.advanceTimersByTime(300);
    expect(defaultProps.onSearch).toHaveBeenCalledWith('');

    vi.useRealTimers();
  });
});
