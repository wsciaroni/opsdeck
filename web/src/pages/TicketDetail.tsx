import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getTicket, updateTicket } from '../api/tickets';
import { ArrowLeft } from 'lucide-react';
import { useState, useEffect } from 'react';

export default function TicketDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // Local state for optimistic/immediate UI updates if needed,
  // though we rely on query invalidation usually.
  // We will control the selects with local state derived from data,
  // or just use the data directly if we want to wait for server roundtrip.
  // The requirement says "Trigger the mutation immediately on change."

  const { data: ticket, isLoading, isError } = useQuery({
    queryKey: ['ticket', id],
    queryFn: () => getTicket(id!),
    enabled: !!id,
  });

  const mutation = useMutation({
    mutationFn: (data: { status_id?: string; priority_id?: string }) =>
      updateTicket(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ticket', id] });
      queryClient.invalidateQueries({ queryKey: ['tickets'] });
    },
  });

  const handleStatusChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    mutation.mutate({ status_id: e.target.value });
  };

  const handlePriorityChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    mutation.mutate({ priority_id: e.target.value });
  };

  if (isLoading) {
    return <div className="p-8 text-center text-gray-500">Loading ticket details...</div>;
  }

  if (isError || !ticket) {
    return (
      <div className="p-8 text-center">
        <p className="text-red-500 mb-4">Ticket not found or error loading ticket.</p>
        <button
            onClick={() => navigate('/')}
            className="text-indigo-600 hover:text-indigo-800"
        >
            Back to Dashboard
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="flex items-center mb-8">
        <button
          onClick={() => navigate('/')}
          className="mr-4 p-2 rounded-full hover:bg-gray-100 text-gray-500"
          aria-label="Back to Dashboard"
        >
          <ArrowLeft className="h-6 w-6" />
        </button>
        <h1 className="text-3xl font-bold text-gray-900">{ticket.title}</h1>
      </div>

      {/* Toolbar */}
      <div className="bg-white shadow rounded-lg p-6 mb-8 flex flex-col sm:flex-row gap-4 items-start sm:items-center border border-gray-200">
        <div className="flex-1">
            <h3 className="text-sm font-medium text-gray-500 uppercase tracking-wide mb-1">Status</h3>
            <select
                value={ticket.status_id}
                onChange={handleStatusChange}
                disabled={mutation.isPending}
                className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
            >
                <option value="new">New</option>
                <option value="in_progress">In Progress</option>
                <option value="on_hold">On Hold</option>
                <option value="done">Done</option>
                <option value="canceled">Canceled</option>
            </select>
        </div>

        <div className="flex-1">
            <h3 className="text-sm font-medium text-gray-500 uppercase tracking-wide mb-1">Priority</h3>
            <select
                value={ticket.priority_id}
                onChange={handlePriorityChange}
                disabled={mutation.isPending}
                className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
            >
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
                <option value="critical">Critical</option>
            </select>
        </div>
      </div>

      {/* Info Grid */}
      <div className="bg-white shadow rounded-lg border border-gray-200 overflow-hidden">
        <div className="px-4 py-5 sm:px-6 bg-gray-50 border-b border-gray-200">
          <h3 className="text-lg leading-6 font-medium text-gray-900">Ticket Information</h3>
        </div>
        <div className="border-t border-gray-200">
          <dl>
            <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 border-b border-gray-100">
              <dt className="text-sm font-medium text-gray-500">Reporter</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">{ticket.reporter_name || 'Unknown'}</dd>
            </div>
            <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 border-b border-gray-100">
              <dt className="text-sm font-medium text-gray-500">Created At</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {new Date(ticket.created_at).toLocaleString()}
              </dd>
            </div>
             <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="text-sm font-medium text-gray-500">Description</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                 <div className="bg-gray-50 p-4 rounded border border-gray-200 whitespace-pre-wrap">
                    {ticket.description}
                 </div>
              </dd>
            </div>
          </dl>
        </div>
      </div>
    </div>
  );
}
