import { useMemo, memo } from 'react';
import { useNavigate } from 'react-router-dom';
import { type Ticket, TICKET_STATUSES } from '../../types';
import { PriorityLabel } from '../TicketAttributes';
import clsx from 'clsx';
import { type Density } from './TicketList';

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

const TicketCard = memo(function TicketCard({ ticket, density }: TicketCardProps) {
  const navigate = useNavigate();

  const paddingClass = {
    compact: 'p-2',
    standard: 'p-4',
    comfortable: 'p-6',
  }[density];

  const fontSizeClass = {
    compact: 'text-xs',
    standard: 'text-sm',
    comfortable: 'text-base',
  }[density];

  const handleKeyDown = (e: React.KeyboardEvent, ticketId: string) => {
    if (e.key === 'Enter' || e.key === ' ') {
      navigate(`/tickets/${ticketId}`);
    }
  };

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
          <span className="text-xs text-gray-400">{new Date(ticket.created_at).toLocaleDateString()}</span>
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

const TicketBoard = memo(function TicketBoard({
  tickets,
  isLoading,
  error,
  density,
  visibleStatuses,
}: TicketBoardProps) {
  // Memoize grouping logic to prevent O(N) recalculation on every render (e.g. density change or modal open)
  const ticketsByStatus = useMemo(() => {
    return (tickets || []).reduce((acc, ticket) => {
      const status = ticket.status_id;
      if (!acc[status]) acc[status] = [];
      acc[status].push(ticket);
      return acc;
    }, {} as Record<string, Ticket[]>);
  }, [tickets]);

  const columns = useMemo(() => {
    if (visibleStatuses && visibleStatuses.length > 0) {
      return TICKET_STATUSES.filter((status) => visibleStatuses.includes(status.id));
    }
    // Default view: Show active statuses (not finished)
    return TICKET_STATUSES.filter((status) => !status.isFinished);
  }, [visibleStatuses]);

  const columnWidthClass = {
    compact: 'min-w-[14rem]',
    standard: 'min-w-[16rem]',
    comfortable: 'min-w-[18rem]',
  }[density];

  if (isLoading) return <div className="p-8 text-center text-gray-500">Loading tickets...</div>;
  if (error) return <div className="p-8 text-center text-red-500">Error loading tickets</div>;

  return (
    <div className="flex h-full overflow-x-auto space-x-4 pb-4">
      {columns.map((column) => (
        <div
          key={column.id}
          className={clsx(
            'flex-1 bg-gray-100 rounded-lg flex flex-col max-h-[calc(100vh-12rem)]',
            columnWidthClass
          )}
        >
          <div className="p-3 font-semibold text-gray-700 flex justify-between items-center sticky top-0 bg-gray-100 z-10 rounded-t-lg">
            <span>{column.label}</span>
            <span className="bg-gray-200 text-gray-600 text-xs px-2 py-0.5 rounded-full">
              {ticketsByStatus[column.id]?.length || 0}
            </span>
          </div>
          <div className="p-2 overflow-y-auto flex-1 space-y-2">
            {ticketsByStatus[column.id]?.map((ticket) => (
              <TicketCard key={ticket.id} ticket={ticket} density={density} />
            ))}
            {!ticketsByStatus[column.id]?.length && (
              <div className="text-center text-gray-400 text-sm py-4 italic">No tickets</div>
            )}
          </div>
        </div>
      ))}
    </div>
  );
});

export default TicketBoard;
