import { Plus, List, Layout, Search, X } from 'lucide-react';
import { type Density } from './TicketList';
import clsx from 'clsx';
import type { Organization } from '../../types';
import FilterPopover from './FilterPopover';
import { useState, useEffect, memo } from 'react';

interface DashboardHeaderProps {
  currentOrg: Organization | null;
  onOpenNewTicket: () => void;
  viewMode: 'list' | 'board';
  setViewMode: (mode: 'list' | 'board') => void;
  density: Density;
  setDensity: (density: Density) => void;
  onSearch: (query: string) => void;
  priority: string[] | undefined;
  setPriority: (priority: string[] | undefined) => void;
  status: string[] | undefined;
  setStatus: (status: string[] | undefined) => void;
  sortBy: string;
  setSortBy: (sortBy: string) => void;
  sortOrder: 'asc' | 'desc';
  setSortOrder: (sortOrder: 'asc' | 'desc') => void;
}

// Memoized to prevent unnecessary re-renders when parent Dashboard updates (e.g. on data fetch)
const DashboardHeader = memo(function DashboardHeader({
  currentOrg,
  onOpenNewTicket,
  viewMode,
  setViewMode,
  density,
  setDensity,
  onSearch,
  priority,
  setPriority,
  status,
  setStatus,
  sortBy,
  setSortBy,
  sortOrder,
  setSortOrder,
}: DashboardHeaderProps) {
  const [inputValue, setInputValue] = useState('');

  useEffect(() => {
    const timer = setTimeout(() => {
      onSearch(inputValue);
    }, 300);
    return () => clearTimeout(timer);
  }, [inputValue, onSearch]);

  return (
    <div className="mb-6">
      <div className="sm:flex sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-gray-900">Tickets</h1>
          <p className="mt-2 text-sm text-gray-700">
            {currentOrg ? `Manage tickets for ${currentOrg.name}` : 'Select an organization'}
          </p>
        </div>
        <div className="mt-4 sm:mt-0 flex space-x-3">
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
        </div>
      </div>

      <div className="mt-6 flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
        {/* Left controls: Search & Filters */}
        <div className="flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-2 w-full sm:w-auto">
            <div className="relative rounded-md shadow-sm w-full sm:w-64">
                <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                    <Search className="h-4 w-4 text-gray-400" aria-hidden="true" />
                </div>
                <label htmlFor="search" className="sr-only">Search tickets</label>
                <input
                    type="text"
                    name="search"
                    id="search"
                    className="block w-full rounded-md border-gray-300 pl-10 pr-10 focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm py-2 border"
                    placeholder="Search tickets..."
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                />
                {inputValue && (
                  <div className="absolute inset-y-0 right-0 flex items-center pr-3">
                    <button
                      type="button"
                      onClick={() => setInputValue('')}
                      className="text-gray-400 hover:text-gray-500 focus:outline-none focus:text-gray-500"
                      aria-label="Clear search"
                    >
                      <X className="h-4 w-4" aria-hidden="true" />
                    </button>
                  </div>
                )}
            </div>

            <FilterPopover
                status={status}
                setStatus={setStatus}
                priority={priority}
                setPriority={setPriority}
            />
        </div>

        {/* Right controls: View Mode & Density */}
        <div className="flex items-center space-x-4">
             <div className="relative inline-block text-left">
                 <select
                    aria-label="Sort tickets"
                    value={`${sortBy}-${sortOrder}`}
                    onChange={(e) => {
                        const [field, order] = e.target.value.split('-');
                        setSortBy(field);
                        setSortOrder(order as 'asc' | 'desc');
                    }}
                    className="block w-full rounded-md border-gray-300 py-1.5 pl-3 pr-8 text-base focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm border"
                 >
                    <option value="created_at-desc">Newest First</option>
                    <option value="created_at-asc">Oldest First</option>
                    <option value="updated_at-desc">Recently Updated</option>
                    <option value="updated_at-asc">Least Recently Updated</option>
                    <option value="priority-desc">Priority (High-Low)</option>
                    <option value="priority-asc">Priority (Low-High)</option>
                    <option value="status-asc">Status (New-Done)</option>
                    <option value="status-desc">Status (Done-New)</option>
                 </select>
            </div>

            <div className="flex items-center space-x-1 border rounded-md p-1 bg-gray-50">
                <button
                    onClick={() => setViewMode('list')}
                    className={clsx(
                        "p-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500",
                        viewMode === 'list' ? "bg-white shadow text-indigo-600" : "text-gray-500 hover:text-gray-700"
                    )}
                    title="List View"
                    aria-label="List View"
                >
                    <List className="h-4 w-4" />
                </button>
                <button
                    onClick={() => setViewMode('board')}
                    className={clsx(
                        "p-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500",
                        viewMode === 'board' ? "bg-white shadow text-indigo-600" : "text-gray-500 hover:text-gray-700"
                    )}
                    title="Board View"
                    aria-label="Board View"
                >
                    <Layout className="h-4 w-4" />
                </button>
            </div>

            <div className="relative inline-block text-left">
                 <select
                    aria-label="View density"
                    value={density}
                    onChange={(e) => setDensity(e.target.value as Density)}
                    className="block w-full rounded-md border-gray-300 py-1.5 pl-3 pr-8 text-base focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm border"
                 >
                    <option value="compact">Compact</option>
                    <option value="standard">Standard</option>
                    <option value="comfortable">Comfortable</option>
                 </select>
            </div>
        </div>
      </div>
    </div>
  );
});

export default DashboardHeader;
