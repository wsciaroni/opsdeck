import { useNavigate } from 'react-router-dom';
import type { Ticket } from '../../types';
import { PriorityLabel } from '../TicketAttributes';
import clsx from 'clsx';
import { type Density } from './TicketList';

interface TicketBoardProps {
  tickets: Ticket[] | undefined;
  isLoading: boolean;
  error: Error | null;
  density: Density;
  onOpenNewTicket: () => void;
}

const STATUS_COLUMNS = [
  { id: 'new', label: 'New' },
  { id: 'in_progress', label: 'In Progress' },
  { id: 'on_hold', label: 'On Hold' },
  { id: 'done', label: 'Done' },
  { id: 'canceled', label: 'Canceled' },
];

export default function TicketBoard({ tickets, isLoading, error, density }: TicketBoardProps) {
  const navigate = useNavigate();

  if (isLoading) return <div className="p-8 text-center text-gray-500">Loading tickets...</div>;
  if (error) return <div className="p-8 text-center text-red-500">Error loading tickets</div>;

  const ticketsByStatus = (tickets || []).reduce((acc, ticket) => {
    const status = ticket.status_id;
    if (!acc[status]) acc[status] = [];
    acc[status].push(ticket);
    return acc;
  }, {} as Record<string, Ticket[]>);

  // Density styles for cards
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
    <div className="flex h-full overflow-x-auto space-x-4 pb-4">
      {STATUS_COLUMNS.map((column) => (
        <div key={column.id} className="flex-shrink-0 w-72 bg-gray-100 rounded-lg flex flex-col max-h-[calc(100vh-12rem)]">
          <div className="p-3 font-semibold text-gray-700 flex justify-between items-center sticky top-0 bg-gray-100 z-10 rounded-t-lg">
            <span>{column.label}</span>
            <span className="bg-gray-200 text-gray-600 text-xs px-2 py-0.5 rounded-full">
              {ticketsByStatus[column.id]?.length || 0}
            </span>
          </div>
          <div className="p-2 overflow-y-auto flex-1 space-y-2">
            {ticketsByStatus[column.id]?.map((ticket) => (
              <div
                key={ticket.id}
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
                    <span>{ticket.assignee_user_id || 'Unassigned'}</span>
                </div>
              </div>
            ))}
             {ticketsByStatus[column.id]?.length === 0 && (
                <div className="text-center text-gray-400 text-sm py-4 italic">
                    No tickets
                </div>
             )}
          </div>
        </div>
      ))}
    </div>
  );
}
