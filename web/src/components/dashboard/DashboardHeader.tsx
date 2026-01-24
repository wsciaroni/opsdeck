import { Plus, List, Layout, Search, Filter, SlidersHorizontal } from 'lucide-react';
import { type Density } from './TicketList';
import clsx from 'clsx';
import type { Organization } from '../../types';

interface DashboardHeaderProps {
  currentOrg: Organization | null;
  onOpenNewTicket: () => void;
  viewMode: 'list' | 'board';
  setViewMode: (mode: 'list' | 'board') => void;
  density: Density;
  setDensity: (density: Density) => void;
  searchQuery: string;
  onSearchChange: (query: string) => void;
}

export default function DashboardHeader({
  currentOrg,
  onOpenNewTicket,
  viewMode,
  setViewMode,
  density,
  setDensity,
  searchQuery,
  onSearchChange,
}: DashboardHeaderProps) {
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
            className="inline-flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
          >
            <Plus className="h-4 w-4 mr-2" />
            New Ticket
          </button>
        </div>
      </div>

      <div className="mt-6 flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
        {/* Left controls: Search & Filter Placeholders */}
        <div className="flex space-x-2 w-full sm:w-auto">
            <div className="relative rounded-md shadow-sm w-full sm:w-64">
                <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                    <Search className="h-4 w-4 text-gray-400" aria-hidden="true" />
                </div>
                <input
                    type="text"
                    name="search"
                    id="search"
                    className="block w-full rounded-md border-gray-300 pl-10 focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm py-2 border"
                    placeholder="Search tickets..."
                    value={searchQuery}
                    onChange={(e) => onSearchChange(e.target.value)}
                />
            </div>
            <button
                type="button"
                className="inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none disabled:opacity-50"
                disabled // Placeholder
                title="Filter (Coming soon)"
            >
                <Filter className="h-4 w-4 mr-2 text-gray-500" />
                Filter
            </button>
             <button
                type="button"
                className="inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none disabled:opacity-50"
                disabled // Placeholder
                title="Sort (Coming soon)"
            >
                <SlidersHorizontal className="h-4 w-4 mr-2 text-gray-500" />
                Sort
            </button>
        </div>

        {/* Right controls: View Mode & Density */}
        <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-1 border rounded-md p-1 bg-gray-50">
                <button
                    onClick={() => setViewMode('list')}
                    className={clsx(
                        "p-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500",
                        viewMode === 'list' ? "bg-white shadow text-indigo-600" : "text-gray-500 hover:text-gray-700"
                    )}
                    title="List View"
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
                >
                    <Layout className="h-4 w-4" />
                </button>
            </div>

            <div className="relative inline-block text-left">
                 <select
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
}
