import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import CreateTicketModal from './CreateTicketModal';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import * as ticketsApi from '../../api/tickets';
import type { Ticket } from '../../types';

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
    },
});

describe('CreateTicketModal', () => {
    it('renders and submits form with file', async () => {
        const mockTicket: Ticket = {
            id: '1',
            organization_id: 'org1',
            title: 'Test Ticket',
            description: 'Test Description',
            status_id: 'new',
            priority_id: 'medium',
            reporter_id: 'user1',
            assignee_user_id: null,
            sensitive: false,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            completed_at: null,
            location: '',
        };
        const createTicketSpy = vi.spyOn(ticketsApi, 'createTicket').mockResolvedValue(mockTicket);
        const onCloseMock = vi.fn();

        render(
            <QueryClientProvider client={queryClient}>
                <CreateTicketModal isOpen={true} onClose={onCloseMock} organizationId="org1" />
            </QueryClientProvider>
        );

        // Fill out form
        fireEvent.change(screen.getByLabelText('Title'), { target: { value: 'Test Ticket' } });
        fireEvent.change(screen.getByLabelText('Description'), { target: { value: 'Test Description' } });

        // Upload file
        const file = new File(['hello'], 'hello.png', { type: 'image/png' });
        const input = screen.getByLabelText('Upload files');
        fireEvent.change(input, { target: { files: [file] } });

        expect(screen.getByText('hello.png')).toBeInTheDocument();

        // Submit
        fireEvent.click(screen.getByText('Create'));

        await waitFor(() => {
            expect(createTicketSpy).toHaveBeenCalled();
        });

        const formData = createTicketSpy.mock.calls[0][0] as FormData;
        expect(formData).toBeInstanceOf(FormData);
        expect(formData.get('title')).toBe('Test Ticket');
        expect(formData.get('description')).toBe('Test Description');
        expect(formData.get('organization_id')).toBe('org1');
        expect(formData.get('files')).toBeInstanceOf(File);
        expect((formData.get('files') as File).name).toBe('hello.png');
    });
});
