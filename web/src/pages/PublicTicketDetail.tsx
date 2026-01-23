import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { getPublicTicket } from '../api/public';
import PublicTicketComments from '../components/PublicTicketComments';
import { ArrowLeft } from 'lucide-react';

export default function PublicTicketDetail() {
  const { token, ticketId } = useParams<{ token: string; ticketId: string }>();
  const navigate = useNavigate();

  const { data: ticket, isLoading, isError } = useQuery({
    queryKey: ['publicTicket', token, ticketId],
    queryFn: () => getPublicTicket(token!, ticketId!),
    enabled: !!token && !!ticketId,
  });

  if (isLoading) {
    return <div className="p-8 text-center text-gray-500">Loading ticket...</div>;
  }

  if (isError || !ticket) {
    return <div className="p-8 text-center text-red-500">Ticket not found</div>;
  }

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
        <h1 className="text-3xl font-bold text-gray-900">
          {ticket.title}
        </h1>
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
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">{ticket.reporter_name || 'Unknown'}</dd>
            </div>
            <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="text-sm font-medium text-gray-500">Assignee</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">{ticket.assignee_name || 'Unassigned'}</dd>
            </div>
            <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="text-sm font-medium text-gray-500">Status</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2 capitalize">{ticket.status_id.replace('_', ' ')}</dd>
            </div>
            <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="text-sm font-medium text-gray-500">Priority</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2 capitalize">{ticket.priority_id}</dd>
            </div>
            <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
              <dt className="text-sm font-medium text-gray-500">Created At</dt>
              <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {new Date(ticket.created_at).toLocaleString()}
              </dd>
            </div>
            <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
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
      <PublicTicketComments />
    </div>
  );
}
