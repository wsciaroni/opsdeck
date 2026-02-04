import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import FileUpload from './FileUpload';

// Mock toast to avoid errors in test output
vi.mock('react-hot-toast', () => ({
  default: {
    error: vi.fn(),
  },
}));

describe('FileUpload', () => {
  it('appends files instead of replacing them', async () => {
    const onFilesChange = vi.fn();
    const initialFiles: File[] = [];

    const { rerender } = render(
      <FileUpload files={initialFiles} onFilesChange={onFilesChange} />
    );

    // Using the input directly might be easier if labels are duplicated
    // But let's try to get by label.
    // Since there are two labels pointing to the same ID, testing-library might return the input for either.
    const input = screen.getByLabelText('Attachments');

    // Create first file
    const file1 = new File(['content1'], 'file1.txt', { type: 'text/plain' });

    // Simulate upload
    fireEvent.change(input, { target: { files: [file1] } });

    // First call should have the new file
    expect(onFilesChange).toHaveBeenCalledWith([file1]);

    // Update the component prop to simulate parent state update
    rerender(
      <FileUpload files={[file1]} onFilesChange={onFilesChange} />
    );

    // Create second file
    const file2 = new File(['content2'], 'file2.txt', { type: 'text/plain' });

    // Simulate upload again
    // We need to clear the value first because the component does it, but fireEvent doesn't stick
    // Actually the component does e.target.value = ''

    fireEvent.change(input, { target: { files: [file2] } });

    // With the BUG, this will be called with ONLY [file2]
    // With the FIX, this should be called with [file1, file2]
    expect(onFilesChange).toHaveBeenLastCalledWith([file1, file2]);
  });
});
