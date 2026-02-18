import { useMemo, memo } from 'react';
import { useNavigate, type NavigateFunction } from 'react-router-dom';
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
  navigate: NavigateFunction;
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

const TicketCard = memo(function TicketCard({ ticket, density, navigate }: TicketCardProps) {
  // Optimization: use pre-defined constants
  const paddingClass = PADDING_CLASSES[density];
  const fontSizeClass = FONT_SIZE_CLASSES[density];

  const handleKeyDown = (e: React.KeyboardEvent, ticketId: string) => {
    if (e.key === 'Enter' || e.key === ' ') {
      navigate(`/tickets/${ticketId}`);
    }
  };

  const createdDate = formatDate(ticket.created_at);

  return (
    <div
      role="button"
      tabIndex={0}
      onClick={() => navigate(`/tickets/${ticket.id}`)}
      onKeyDown={(e) => handleKeyDown(e, ticket.id)}
      className={clsx(
        "bg-white rounded border border-gray-200 shadow-sm cursor-pointer hover:shadow-md transition-shadow focus:outline-none focus:ring-2 focus:ring-indigo-500",
        paddingClass
      )}
    >
      <div className="flex justify-between items-start mb-2">
          <PriorityLabel priority={ticket.priority_id} />
          <span className="text-xs text-gray-600" aria-label={`Created on ${createdDate}`}>
            {createdDate}
          </span>
      </div>
      <h4 className={clsx("font-medium text-gray-900 mb-2 line-clamp-2", fontSizeClass)}>
        {ticket.title}
      </h4>
      <div className="flex justify-between items-center text-xs text-gray-500 mt-auto">
          <span>{ticket.assignee_name || ticket.assignee_user_id || 'Unassigned'}</span>
      </div>
    </div>
  );
});

interface TicketColumnProps {
  column: typeof TICKET_STATUSES[number];
  tickets: Ticket[] | undefined;
  density: Density;
  navigate: NavigateFunction;
}

// Optimized: Memoize column to prevent re-rendering unchanged columns when other tickets change.
const TicketColumn = memo(function TicketColumn({ column, tickets, density, navigate }: TicketColumnProps) {
  const columnWidthClass = COLUMN_WIDTH_CLASSES[density];

  return (
    <div
      className={clsx(
        'flex-1 bg-gray-100 rounded-lg flex flex-col max-h-[calc(100vh-12rem)]',
        columnWidthClass
      )}
    >
      <div className="p-3 font-semibold text-gray-700 flex justify-between items-center sticky top-0 bg-gray-100 z-10 rounded-t-lg">
        <span>{column.label}</span>
        <span className="bg-gray-200 text-gray-600 text-xs px-2 py-0.5 rounded-full">
          {tickets?.length || 0}
        </span>
      </div>
      <div className="p-2 overflow-y-auto flex-1 space-y-2">
        {tickets?.map((ticket) => (
          <TicketCard
            key={ticket.id}
            ticket={ticket}
            density={density}
            navigate={navigate}
          />
        ))}
        {!tickets?.length && (
          <div className="text-center text-gray-500 text-sm py-4 italic">No tickets</div>
        )}
      </div>
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
  // Optimization: Hoist useNavigate to avoid calling it in every card
  const navigate = useNavigate();

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
          navigate={navigate}
        />
      ))}
    </div>
  );
});

export default TicketBoard;
