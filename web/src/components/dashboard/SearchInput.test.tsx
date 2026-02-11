import { render, screen, fireEvent } from '@testing-library/react';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import SearchInput from './SearchInput';

describe('SearchInput', () => {
  const mockOnSearch = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders with placeholder', () => {
    render(<SearchInput onSearch={mockOnSearch} placeholder="Test Search" />);
    expect(screen.getByPlaceholderText('Test Search')).toBeInTheDocument();
  });

  it('debounces search input', () => {
    vi.useFakeTimers();
    render(<SearchInput onSearch={mockOnSearch} />);

    const input = screen.getByPlaceholderText('Search...');
    fireEvent.change(input, { target: { value: 'test' } });

    // Should not call immediately
    expect(mockOnSearch).not.toHaveBeenCalled();

    // Advance timer
    vi.advanceTimersByTime(300);

    expect(mockOnSearch).toHaveBeenCalledWith('test');

    vi.useRealTimers();
  });

  it('clears input and calls onSearch with empty string', () => {
    vi.useFakeTimers();
    render(<SearchInput onSearch={mockOnSearch} />);

    const input = screen.getByPlaceholderText('Search...');
    fireEvent.change(input, { target: { value: 'test' } });

    // Fast forward to first search
    vi.advanceTimersByTime(300);
    expect(mockOnSearch).toHaveBeenCalledWith('test');

    // Find clear button (it appears when inputValue is truthy)
    const clearButton = screen.getByLabelText('Clear search');
    fireEvent.click(clearButton);

    expect(input).toHaveValue('');

    // Fast forward debounce for clearing
    vi.advanceTimersByTime(300);
    expect(mockOnSearch).toHaveBeenCalledWith('');

    vi.useRealTimers();
  });
});
