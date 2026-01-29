import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import { useQuery } from '@tanstack/react-query';
import { getTickets } from '../api/tickets';
import CreateTicketModal from '../components/dashboard/CreateTicketModal';
import TicketList, { type Density } from '../components/dashboard/TicketList';
import TicketBoard from '../components/dashboard/TicketBoard';
import DashboardHeader from '../components/dashboard/DashboardHeader';

export default function Dashboard() {
  const { currentOrg } = useAuth();
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleOpenNewTicket = useCallback(() => {
    setIsModalOpen(true);
  }, []);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Ignore if user is typing in an input, textarea, or contentEditable element
      const target = e.target as HTMLElement;
      if (
        target.tagName === 'INPUT' ||
        target.tagName === 'TEXTAREA' ||
        target.tagName === 'SELECT' ||
        target.isContentEditable
      ) {
        return;
      }

      // Ignore if modifiers are pressed
      if (e.ctrlKey || e.altKey || e.metaKey || e.shiftKey) return;

      if (e.key.toLowerCase() === 'c') {
        e.preventDefault();
        handleOpenNewTicket();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [handleOpenNewTicket]);

  // State for view preferences with local storage persistence
  const [viewMode, setViewMode] = useState<'list' | 'board'>(() => {
    return (localStorage.getItem('dashboard_view_mode') as 'list' | 'board') || 'list';
  });

  const [density, setDensity] = useState<Density>(() => {
    return (localStorage.getItem('dashboard_density') as Density) || 'standard';
  });

  // Filter states
  const [search, setSearch] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [priority, setPriority] = useState<string[] | undefined>(undefined);
  const [status, setStatus] = useState<string[] | undefined>(undefined);
  const [sortBy, setSortBy] = useState<string>(() => localStorage.getItem('dashboard_sort_by') || 'created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>(() => (localStorage.getItem('dashboard_sort_order') as 'asc' | 'desc') || 'desc');

  // Debounce search
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(search);
    }, 300);
    return () => clearTimeout(timer);
  }, [search]);

  // Persist preferences
  useEffect(() => {
    localStorage.setItem('dashboard_view_mode', viewMode);
  }, [viewMode]);

  useEffect(() => {
    localStorage.setItem('dashboard_density', density);
  }, [density]);

  useEffect(() => {
    localStorage.setItem('dashboard_sort_by', sortBy);
  }, [sortBy]);

  useEffect(() => {
    localStorage.setItem('dashboard_sort_order', sortOrder);
  }, [sortOrder]);

  const { data: tickets, isLoading, error } = useQuery({
    queryKey: ['tickets', currentOrg?.id, debouncedSearch, priority, status, sortBy, sortOrder],
    queryFn: () => getTickets(currentOrg!.id, { search: debouncedSearch, priority, status, sort_by: sortBy, sort_order: sortOrder }),
    enabled: !!currentOrg,
  });

  if (!currentOrg) {
    return (
      <div className="p-8 text-center text-gray-500">
        Please select an organization.
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 h-full flex flex-col flex-1">
      <DashboardHeader
        currentOrg={currentOrg}
        onOpenNewTicket={handleOpenNewTicket}
        viewMode={viewMode}
        setViewMode={setViewMode}
        density={density}
        setDensity={setDensity}
        search={search}
        setSearch={setSearch}
        priority={priority}
        setPriority={setPriority}
        status={status}
        setStatus={setStatus}
        sortBy={sortBy}
        setSortBy={setSortBy}
        sortOrder={sortOrder}
        setSortOrder={setSortOrder}
      />

      <div className="flex-1 overflow-hidden">
        {viewMode === 'list' ? (
          <div className="h-full overflow-y-auto">
             <TicketList
                tickets={tickets}
                isLoading={isLoading}
                error={error}
                density={density}
                onOpenNewTicket={handleOpenNewTicket}
            />
          </div>
        ) : (
          <div className="h-full">
            <TicketBoard
                tickets={tickets}
                isLoading={isLoading}
                error={error}
                density={density}
                visibleStatuses={status}
                onOpenNewTicket={handleOpenNewTicket}
            />
          </div>
        )}
      </div>

      <CreateTicketModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        organizationId={currentOrg.id}
      />
    </div>
  );
}
