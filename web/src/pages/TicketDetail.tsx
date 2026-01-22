import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getTicket, updateTicket } from '../api/tickets';
import TicketComments from '../components/TicketComments';
import { ArrowLeft } from 'lucide-react';
import type { Ticket } from '../types';
import toast from 'react-hot-toast';

export default function TicketDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: ticket, isLoading, isError } = useQuery({
    queryKey: ['ticket', id],
    queryFn: () => getTicket(id!),
    enabled: !!id,
  });

  const mutation = useMutation({
    mutationFn: (data: { status_id?: string; priority_id?: string }) =>
      updateTicket(id!, data),
    onSuccess: (updatedTicket) => {
      queryClient.setQueryData(['ticket', id], (oldData: Ticket) => ({
        ...oldData,
        ...updatedTicket,
      }));
       // Also invalidate the list to update the dashboard
      queryClient.invalidateQueries({ queryKey: ['tickets'] });
      toast.success("Ticket updated!");
    },
  });

  if (isLoading) {
    return <div className="p-8 text-center text-gray-500">Loading ticket...</div>;
  }

  if (isError || !ticket) {
    return <div className="p-8 text-center text-red-500">Ticket not found</div>;
  }

  const handleStatusChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    mutation.mutate({ status_id: e.target.value });
  };

  const handlePriorityChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    mutation.mutate({ priority_id: e.target.value });
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="flex items-center mb-6">
        <button
          onClick={() => navigate(-1)}
          className="mr-4 p-2 rounded-full hover:bg-gray-100 text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          aria-label="Go back"
          title="Go back"
        >
          <ArrowLeft className="h-6 w-6" aria-hidden="true" />
        </button>
        <h1 className="text-3xl font-bold text-gray-900">{ticket.title}</h1>
      </div>

      {/* Toolbar */}
      <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200 mb-8 flex flex-wrap gap-4 items-center">
        <div>
          <label htmlFor="status" className="block text-sm font-medium text-gray-700 mr-2 mb-1">Status</label>
          <select
            id="status"
            value={ticket.status_id}
            onChange={handleStatusChange}
            disabled={mutation.isPending}
            className="block w-40 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
          >
            <option value="new">New</option>
            <option value="in_progress">In Progress</option>
            <option value="on_hold">On Hold</option>
            <option value="done">Done</option>
            <option value="canceled">Canceled</option>
          </select>
        </div>

        <div>
          <label htmlFor="priority" className="block text-sm font-medium text-gray-700 mr-2 mb-1">Priority</label>
          <select
            id="priority"
            value={ticket.priority_id}
            onChange={handlePriorityChange}
            disabled={mutation.isPending}
            className="block w-40 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
          >
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="critical">Critical</option>
          </select>
        </div>
      </div>

      {/* Info Grid */}
      <div className="bg-white shadow overflow-hidden sm:rounded-lg">
        <div className="px-4 py-5 sm:px-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900">Ticket Details</h3>
        </div>
        <div className="border-t border-gray-200">
          <dl>
            <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="text-sm font-medium text-gray-500">Reporter</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">{ticket.reporter_name}</dd>
            </div>
            <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="text-sm font-medium text-gray-500">Created At</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {new Date(ticket.created_at).toLocaleString()}
              </dd>
            </div>
            <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="text-sm font-medium text-gray-500">Description</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                <div className="bg-gray-100 p-4 rounded-md border border-gray-200 whitespace-pre-wrap">
                  {ticket.description}
                </div>
              </dd>
            </div>
          </dl>
        </div>
      </div>

      {/* Comments Section */}
      <TicketComments ticketId={id!} />
    </div>
  );
}
