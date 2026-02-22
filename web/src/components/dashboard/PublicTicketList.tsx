import { useNavigate, useParams, type NavigateFunction } from 'react-router-dom';
import type { PublicTicket } from '../../api/public';
import { StatusBadge, PriorityLabel } from '../TicketAttributes';
import EmptyState from '../EmptyState';
import { Inbox } from 'lucide-react';
import { memo } from 'react';
import { useIsDesktop } from '../../hooks/useIsDesktop';
import { formatDate } from '../../utils';

interface PublicTicketListProps {
  tickets: PublicTicket[] | undefined;
  isLoading: boolean;
  error: Error | null;
}

const MobileTicketCard = memo(function MobileTicketCard({ ticket, token, navigate }: { readonly ticket: PublicTicket; readonly token: string; readonly navigate: NavigateFunction }) {
  return (
    <li className="block bg-white hover:bg-gray-50 cursor-pointer">
      <button
        onClick={() => navigate(`/public/${token}/tickets/${ticket.id}`)}
        className="w-full text-left px-4 py-4 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500"
      >
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center space-x-2">
            <StatusBadge status={ticket.status_id} />
            <PriorityLabel priority={ticket.priority_id} />
          </div>
          <div className="text-xs text-gray-500">
            {formatDate(ticket.created_at)}
          </div>
        </div>
        <div className="mb-2">
          <h3 className="text-sm font-semibold text-gray-900 line-clamp-2">{ticket.title}</h3>
        </div>
      </button>
    </li>
  );
});

const PublicTicketRow = memo(function PublicTicketRow({ ticket, token, navigate }: { readonly ticket: PublicTicket; readonly token: string; readonly navigate: NavigateFunction }) {
  const handleClick = () => navigate(`/public/${token}/tickets/${ticket.id}`);

  return (
    <tr
      onClick={handleClick}
      className="cursor-pointer hover:bg-gray-50 focus:outline-none focus:bg-gray-50 focus:ring-2 focus:ring-inset focus:ring-indigo-500"
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          handleClick();
        }
      }}
      aria-label={`View ticket: ${ticket.title}`}
    >
      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500 text-left">
        <StatusBadge status={ticket.status_id} />
      </td>
      <td className="whitespace-nowrap px-3 py-4 text-sm font-bold text-gray-900 text-left">
        {ticket.title}
      </td>
      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500 text-left">
        <PriorityLabel priority={ticket.priority_id} />
      </td>
      <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500 text-left">
        {formatDate(ticket.created_at)}
      </td>
    </tr>
  );
});

const PublicTicketList = memo(function PublicTicketList({ tickets, isLoading, error }: PublicTicketListProps) {
  const navigate = useNavigate();
  const { token } = useParams<{ token: string }>();
  // Optimization: Conditionally render mobile or desktop view to reduce DOM nodes by ~50%
  const isDesktop = useIsDesktop();

  if (isLoading) {
    return <div className="bg-white shadow rounded-lg p-8 text-center text-gray-500">Loading tickets...</div>;
  }

  if (error) {
    return <div className="bg-white shadow rounded-lg p-8 text-center text-red-500">Error loading tickets</div>;
  }

  if (tickets?.length === 0) {
    return (
      <div className="bg-white shadow rounded-lg">
        <EmptyState
          title="No tickets found"
          description="There are no public tickets visible at the moment."
          icon={Inbox}
          action={
             // Optional: Link to public submit if needed, but maybe not here.
             null
          }
        />
      </div>
    );
  }

  return (
    <>
      {/* Mobile View */}
      {!isDesktop && (
        <div className="md:hidden bg-white shadow overflow-hidden rounded-md border border-gray-200">
          <ul className="divide-y divide-gray-200">
            {tickets?.map((ticket) => (
              <MobileTicketCard
                key={ticket.id}
                ticket={ticket}
                token={token!}
                navigate={navigate}
              />
            ))}
          </ul>
        </div>
      )}

      {/* Desktop View */}
      {isDesktop && (
        <div className="hidden md:flex flex-col">
          <div className="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div className="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
              <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg bg-white">
                <table className="min-w-full divide-y divide-gray-300">
                  <thead className="bg-gray-50">
                    <tr>
                      <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Status</th>
                      <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Title</th>
                      <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Priority</th>
                      <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Created</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 bg-white">
                    {tickets?.map((ticket) => (
                      <PublicTicketRow
                        key={ticket.id}
                        ticket={ticket}
                        token={token!}
                        navigate={navigate}
                      />
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
});

export default PublicTicketList;
