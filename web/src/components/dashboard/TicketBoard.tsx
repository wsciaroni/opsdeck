import { useMemo, memo } from 'react';
import { Link } from 'react-router-dom';
import { type Ticket, TICKET_STATUSES } from '../../types';
import { PriorityLabel } from '../TicketAttributes';
import clsx from 'clsx';
import { type Density } from './TicketList';
import { formatDate } from '../../utils';

interface TicketBoardProps {
  tickets: Ticket[] | undefined;
  isLoading: boolean;
  error: Error | null;
  density: Density;
  visibleStatuses?: string[];
  onOpenNewTicket: () => void;
}

interface TicketCardProps {
  ticket: Ticket;
  density: Density;
}

// Optimization: Hoist style constants to prevent object allocation on every render
const PADDING_CLASSES: Record<Density, string> = {
  compact: 'p-2',
  standard: 'p-4',
  comfortable: 'p-6',
};

const FONT_SIZE_CLASSES: Record<Density, string> = {
  compact: 'text-xs',
  standard: 'text-sm',
  comfortable: 'text-base',
};

const COLUMN_WIDTH_CLASSES: Record<Density, string> = {
  compact: 'min-w-[14rem]',
  standard: 'min-w-[16rem]',
  comfortable: 'min-w-[18rem]',
};

const TicketCard = memo(function TicketCard({ ticket, density }: TicketCardProps) {
  // Optimization: use pre-defined constants
  const paddingClass = PADDING_CLASSES[density];
  const fontSizeClass = FONT_SIZE_CLASSES[density];

  const createdDate = formatDate(ticket.created_at);

  return (
    <li className="list-none">
      <Link
        to={`/tickets/${ticket.id}`}
        className={clsx(
          "block bg-white rounded border border-gray-200 shadow-sm hover:shadow-md transition-shadow focus:outline-none focus:ring-2 focus:ring-indigo-500",
          paddingClass
        )}
      >
        <div className="flex justify-between items-start mb-2">
            <PriorityLabel priority={ticket.priority_id} />
            <span className="text-xs text-gray-600" aria-label={`Created on ${createdDate}`}>
              {createdDate}
            </span>
        </div>
        <h3 className={clsx("font-medium text-gray-900 mb-2 line-clamp-2", fontSizeClass)}>
          {ticket.title}
        </h3>
        <div className="flex justify-between items-center text-xs text-gray-500 mt-auto">
            <span>{ticket.assignee_name || ticket.assignee_user_id || 'Unassigned'}</span>
        </div>
      </Link>
    </li>
  );
});

interface TicketColumnProps {
  column: typeof TICKET_STATUSES[number];
  tickets: Ticket[] | undefined;
  density: Density;
}

// Optimized: Memoize column to prevent re-rendering unchanged columns when other tickets change.
const TicketColumn = memo(function TicketColumn({ column, tickets, density }: TicketColumnProps) {
  const columnWidthClass = COLUMN_WIDTH_CLASSES[density];

  return (
    <div
      className={clsx(
        'flex-1 bg-gray-100 rounded-lg flex flex-col max-h-[calc(100vh-12rem)]',
        columnWidthClass
      )}
    >
      <div className="p-3 font-semibold text-gray-700 flex justify-between items-center sticky top-0 bg-gray-100 z-10 rounded-t-lg">
        <h2 className="text-sm font-semibold m-0">{column.label}</h2>
        <span className="bg-gray-200 text-gray-600 text-xs px-2 py-0.5 rounded-full">
          {tickets?.length || 0}
        </span>
      </div>
      <ul className="p-2 overflow-y-auto flex-1 space-y-2">
        {tickets?.map((ticket) => (
          <TicketCard
            key={ticket.id}
            ticket={ticket}
            density={density}
          />
        ))}
        {!tickets?.length && (
          <li className="text-center text-gray-500 text-sm py-4 italic list-none">No tickets</li>
        )}
      </ul>
    </div>
  );
}, (prevProps, nextProps) => {
  // Custom comparison to prevent re-renders when tickets array reference changes but content is same
  if (prevProps.density !== nextProps.density) return false;
  if (prevProps.column.id !== nextProps.column.id) return false;

  const prevTickets = prevProps.tickets || [];
  const nextTickets = nextProps.tickets || [];

  if (prevTickets.length !== nextTickets.length) return false;

  // Since React Query uses structural sharing, unchanged ticket objects are referentially stable.
  // We can just check reference equality for each item.
  for (let i = 0; i < prevTickets.length; i++) {
    if (prevTickets[i] !== nextTickets[i]) return false;
  }

  return true;
});

const TicketBoard = memo(function TicketBoard({
  tickets,
  isLoading,
  error,
  density,
  visibleStatuses,
}: TicketBoardProps) {
  // Memoize grouping logic to prevent O(N) recalculation on every render (e.g. density change or modal open)
  const ticketsByStatus = useMemo(() => {
    const groups: Record<string, Ticket[]> = {};
    if (!tickets) return groups;
    for (const ticket of tickets) {
      const status = ticket.status_id;
      if (!groups[status]) groups[status] = [];
      groups[status].push(ticket);
    }
    return groups;
  }, [tickets]);

  const columns = useMemo(() => {
    if (visibleStatuses && visibleStatuses.length > 0) {
      return TICKET_STATUSES.filter((status) => visibleStatuses.includes(status.id));
    }
    // Default view: Show active statuses (not finished)
    return TICKET_STATUSES.filter((status) => !status.isFinished);
  }, [visibleStatuses]);

  if (isLoading) return <div className="p-8 text-center text-gray-500">Loading tickets...</div>;
  if (error) return <div className="p-8 text-center text-red-500">Error loading tickets</div>;

  return (
    <div className="flex h-full overflow-x-auto space-x-4 pb-4">
      {columns.map((column) => (
        <TicketColumn
          key={column.id}
          column={column}
          tickets={ticketsByStatus[column.id]}
          density={density}
        />
      ))}
    </div>
  );
});

export default TicketBoard;
