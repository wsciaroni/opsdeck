import { useNavigate, type NavigateFunction } from 'react-router-dom';
import type { Ticket } from '../../types';
import { StatusBadge, PriorityLabel } from '../TicketAttributes';
import EmptyState from '../EmptyState';
import { Inbox, Plus } from 'lucide-react';
import clsx from 'clsx';
import { memo, useState, useEffect } from 'react';

export type Density = 'compact' | 'standard' | 'comfortable';

// Optimization: Hoist style constants to prevent object allocation on every render
const PADDING_CLASSES: Record<Density, string> = {
  compact: 'py-2',
  standard: 'py-4',
  comfortable: 'py-6',
};

const FONT_SIZE_CLASSES: Record<Density, string> = {
  compact: 'text-xs',
  standard: 'text-sm',
  comfortable: 'text-base',
};

const METADATA_FONT_SIZE_CLASSES: Record<Density, string> = {
  compact: 'text-xs',
  standard: 'text-xs',
  comfortable: 'text-sm',
};

function useIsDesktop() {
  const [isDesktop, setIsDesktop] = useState(() => {
    if (typeof window === 'undefined') return true;
    return window.matchMedia('(min-width: 768px)').matches;
  });

  useEffect(() => {
    if (typeof window === 'undefined') return;
    const mediaQuery = window.matchMedia('(min-width: 768px)');
    const handler = (e: MediaQueryListEvent) => setIsDesktop(e.matches);
    mediaQuery.addEventListener('change', handler);
    return () => mediaQuery.removeEventListener('change', handler);
  }, []);

  return isDesktop;
}

interface TicketListProps {
  tickets: Ticket[] | undefined;
  isLoading: boolean;
  error: Error | null;
  density: Density;
  onOpenNewTicket: () => void;
}

const MobileTicketCard = memo(function MobileTicketCard({ ticket, density, navigate }: { readonly ticket: Ticket, density: Density, navigate: NavigateFunction }) {
  // Optimization: use pre-defined constants
  const paddingClass = PADDING_CLASSES[density];
  const fontSizeClass = FONT_SIZE_CLASSES[density];
  const metadataFontSizeClass = METADATA_FONT_SIZE_CLASSES[density];

  return (
    <li className="block bg-white hover:bg-gray-50 cursor-pointer">
      <button
        onClick={() => navigate(`/tickets/${ticket.id}`)}
        className={clsx("w-full text-left px-4 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500", paddingClass)}
      >
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center space-x-2">
            <StatusBadge status={ticket.status_id} />
            <PriorityLabel priority={ticket.priority_id} />
          </div>
          <div className={clsx("text-gray-600", metadataFontSizeClass)}>
            {new Date(ticket.created_at).toLocaleDateString()}
          </div>
        </div>
        <div className="mb-2">
          <h3 className={clsx("font-semibold text-gray-900 line-clamp-2", fontSizeClass)}>{ticket.title}</h3>
        </div>
        <div className={clsx("flex items-center text-gray-500", metadataFontSizeClass)}>
          <span>{ticket.assignee_name || ticket.assignee_user_id || 'Unassigned'}</span>
        </div>
      </button>
    </li>
  );
});

const TicketRow = memo(function TicketRow({ ticket, density, navigate }: { ticket: Ticket; density: Density, navigate: NavigateFunction }) {
  // Optimization: use pre-defined constants
  const paddingClass = PADDING_CLASSES[density];
  const fontSizeClass = FONT_SIZE_CLASSES[density];

  return (
    <tr
      onClick={() => navigate(`/tickets/${ticket.id}`)}
      className="cursor-pointer hover:bg-gray-50 focus:outline-none focus:bg-gray-50 focus:ring-2 focus:ring-inset focus:ring-indigo-500"
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          navigate(`/tickets/${ticket.id}`);
        }
      }}
      aria-label={`View ticket: ${ticket.title}`}
    >
      <td className={clsx("whitespace-nowrap px-3 text-sm text-gray-500 text-left", paddingClass)}>
        <StatusBadge status={ticket.status_id} />
      </td>
      <td className={clsx("whitespace-nowrap px-3 font-bold text-gray-900 text-left", paddingClass, fontSizeClass)}>
        {ticket.title}
      </td>
      <td className={clsx("whitespace-nowrap px-3 text-gray-500 text-left", paddingClass, fontSizeClass)}>
        <PriorityLabel priority={ticket.priority_id} />
      </td>
      <td className={clsx("whitespace-nowrap px-3 text-gray-500 text-left", paddingClass, fontSizeClass)}>
        {ticket.assignee_name || ticket.assignee_user_id || 'Unassigned'}
      </td>
      <td className={clsx("whitespace-nowrap px-3 text-gray-500 text-left", paddingClass, fontSizeClass)}>
        {new Date(ticket.created_at).toLocaleDateString()}
      </td>
    </tr>
  );
});

const TicketList = memo(function TicketList({ tickets, isLoading, error, density, onOpenNewTicket }: TicketListProps) {
  // Optimization: Conditionally render mobile or desktop view to reduce DOM nodes by ~50%
  const isDesktop = useIsDesktop();

  // Optimization: Hoist useNavigate to avoid calling it in every row item
  const navigate = useNavigate();

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
          description="Create your first ticket to get started tracking your work."
          icon={Inbox}
          action={
            <button
              type="button"
              onClick={onOpenNewTicket}
              title="Create new ticket"
              aria-keyshortcuts="c"
              className="inline-flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
            >
              <Plus className="h-4 w-4 mr-2" />
              New Ticket
              <span
                className="ml-2 hidden sm:inline-flex items-center rounded bg-white/20 px-1.5 py-0.5 text-xs font-semibold leading-none text-white"
                aria-hidden="true"
              >
                C
              </span>
            </button>
          }
        />
      </div>
    );
  }

  return (
    <>
      {/* Mobile View - Always standard density for mobile mostly, or strictly card based */}
      {!isDesktop && (
        <div className="md:hidden bg-white shadow overflow-hidden rounded-md border border-gray-200">
          <ul className="divide-y divide-gray-200">
            {tickets?.map((ticket: Ticket) => (
              <MobileTicketCard
                key={ticket.id}
                ticket={ticket}
                density={density}
                navigate={navigate}
              />
            ))}
          </ul>
        </div>
      )}

      {/* Desktop View */}
      {isDesktop && (
        <div className="hidden md:flex flex-col">
          <div className="overflow-x-auto">
            <div className="inline-block min-w-full py-2 align-middle">
              <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg bg-white">
                <table className="min-w-full divide-y divide-gray-300">
                  <thead className="bg-gray-50">
                    <tr>
                      <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Status</th>
                      <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Title</th>
                      <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Priority</th>
                      <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Assignee</th>
                      <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Created</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 bg-white">
                    {tickets?.map((ticket: Ticket) => (
                      <TicketRow
                        key={ticket.id}
                        ticket={ticket}
                        density={density}
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

export default TicketList;
