import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { getPublicOrganization, getPublicTickets } from '../api/public';
import PublicTicketList from '../components/dashboard/PublicTicketList';
import { LayoutDashboard } from 'lucide-react';

export default function PublicDashboard() {
  const { token } = useParams<{ token: string }>();

  const { data: org, isLoading: orgLoading, error: orgError } = useQuery({
    queryKey: ['publicOrg', token],
    queryFn: () => getPublicOrganization(token!),
    enabled: !!token,
  });

  const { data: tickets, isLoading: ticketsLoading, error: ticketsError } = useQuery({
    queryKey: ['publicTickets', token],
    queryFn: () => getPublicTickets(token!),
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
