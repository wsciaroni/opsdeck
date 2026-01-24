import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { getPublicOrganization, getPublicTickets } from '../api/public';
import PublicTicketList from '../components/dashboard/PublicTicketList';
import { LayoutDashboard, Search } from 'lucide-react';

export default function PublicDashboard() {
  const { token } = useParams<{ token: string }>();
  const [searchQuery, setSearchQuery] = useState('');

  const { data: org, isLoading: orgLoading, error: orgError } = useQuery({
    queryKey: ['publicOrg', token],
    queryFn: () => getPublicOrganization(token!),
    enabled: !!token,
  });

  const { data: tickets, isLoading: ticketsLoading, error: ticketsError } = useQuery({
    queryKey: ['publicTickets', token, searchQuery],
    queryFn: () => getPublicTickets(token!, searchQuery),
    enabled: !!token,
  });

  if (orgLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  if (orgError || !org) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-red-500">Organization not found or public view disabled.</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <h1 className="text-3xl font-bold leading-tight text-gray-900 flex items-center gap-2">
            <LayoutDashboard className="h-8 w-8 text-indigo-600" />
            {org.name} <span className="text-gray-400 font-normal text-xl">| Public Dashboard</span>
          </h1>
        </div>
      </header>
      <main>
        <div className="max-w-7xl mx-auto sm:px-6 lg:px-8 py-8">
            <div className="mb-6 flex justify-end">
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
                        onChange={(e) => setSearchQuery(e.target.value)}
                    />
                </div>
            </div>
           <PublicTicketList
             tickets={tickets}
             isLoading={ticketsLoading}
             error={ticketsError}
           />
        </div>
      </main>
    </div>
  );
}
