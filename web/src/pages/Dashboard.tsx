import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getTickets, createTicket } from '../api/tickets';
import type { Ticket, CreateTicketRequest } from '../types';
import { Plus, Inbox, Paperclip } from 'lucide-react';
import EmptyState from '../components/EmptyState';
import toast from 'react-hot-toast';

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    new: 'bg-blue-100 text-blue-800',
    in_progress: 'bg-yellow-100 text-yellow-800',
    on_hold: 'bg-orange-100 text-orange-800',
    done: 'bg-green-100 text-green-800',
    canceled: 'bg-gray-100 text-gray-800',
  };

  const labels: Record<string, string> = {
    new: 'New',
    in_progress: 'In Progress',
    on_hold: 'On Hold',
    done: 'Done',
    canceled: 'Canceled',
  };

  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colors[status] || 'bg-gray-100 text-gray-800'}`}>
      {labels[status] || status}
    </span>
  );
}

function PriorityLabel({ priority }: { priority: string }) {
  const colors: Record<string, string> = {
    low: 'text-gray-500',
    medium: 'text-blue-500',
    high: 'text-orange-500 font-bold',
    critical: 'text-red-600 font-bold uppercase',
  };

  return <span className={`text-sm ${colors[priority] || 'text-gray-500'}`}>{priority}</span>;
}

function MobileTicketCard({ ticket, onClick }: { readonly ticket: Ticket; readonly onClick: () => void }) {
  return (
    <li className="block bg-white hover:bg-gray-50 cursor-pointer">
      <button
        onClick={onClick}
        className="w-full text-left px-4 py-4 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500"
      >
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center space-x-2">
            <StatusBadge status={ticket.status_id} />
            <PriorityLabel priority={ticket.priority_id} />
          </div>
          <div className="text-xs text-gray-500">
            {new Date(ticket.created_at).toLocaleDateString()}
          </div>
        </div>
        <div className="mb-2">
          <h3 className="text-sm font-semibold text-gray-900 line-clamp-2">{ticket.title}</h3>
        </div>
        <div className="flex items-center text-xs text-gray-500">
          <span>{ticket.assignee_user_id || 'Unassigned'}</span>
        </div>
      </button>
    </li>
  );
}

export default function Dashboard() {
  const { currentOrg } = useAuth();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newTicket, setNewTicket] = useState({
    title: '',
    description: '',
    priority_id: 'medium',
  });
  const [files, setFiles] = useState<FileList | null>(null);

  const { data: tickets, isLoading, error } = useQuery({
    queryKey: ['tickets', currentOrg?.id],
    queryFn: () => getTickets(currentOrg!.id),
    enabled: !!currentOrg,
  });

  const mutation = useMutation({
    mutationFn: (data: FormData | CreateTicketRequest) => createTicket(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tickets'] });
      setIsModalOpen(false);
      setNewTicket({ title: '', description: '', priority_id: 'medium' });
      setFiles(null);
      toast.success("Ticket created!");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!currentOrg) return;

    const formData = new FormData();
    formData.append('organization_id', currentOrg.id);
    formData.append('title', newTicket.title);
    formData.append('description', newTicket.description);
    formData.append('priority_id', newTicket.priority_id);
    if (files) {
      for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i]);
      }
    }

    mutation.mutate(formData);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      if (e.target.files && e.target.files.length > 0) {
        setFiles(e.target.files);
      }
  };

  if (!currentOrg) {
    return (
        <div className="p-8 text-center text-gray-500">
            Please select an organization.
        </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="sm:flex sm:items-center">
        <div className="sm:flex-auto">
          <h1 className="text-2xl font-semibold text-gray-900">Tickets</h1>
          <p className="mt-2 text-sm text-gray-700">
            A list of all tickets in <strong>{currentOrg.name}</strong> including their title, status, and priority.
          </p>
        </div>
        <div className="mt-4 sm:mt-0 sm:ml-16 sm:flex-none">
          <button
            type="button"
            onClick={() => setIsModalOpen(true)}
            className="inline-flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:w-auto"
          >
            <Plus className="h-4 w-4 mr-2" />
            New Ticket
          </button>
        </div>
      </div>

      <div className="mt-8">
        {isLoading ? (
          <div className="bg-white shadow rounded-lg p-8 text-center text-gray-500">Loading tickets...</div>
        ) : error ? (
          <div className="bg-white shadow rounded-lg p-8 text-center text-red-500">Error loading tickets</div>
        ) : tickets && tickets.length === 0 ? (
          <div className="bg-white shadow rounded-lg">
            <EmptyState
              title="No tickets found"
              description="Create your first ticket to get started tracking your work."
              icon={Inbox}
              action={
                <button
                  type="button"
                  onClick={() => setIsModalOpen(true)}
                  className="inline-flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                >
                  <Plus className="h-4 w-4 mr-2" />
                  New Ticket
                </button>
              }
            />
          </div>
        ) : (
          <>
            {/* Mobile View */}
            <div className="md:hidden bg-white shadow overflow-hidden rounded-md border border-gray-200">
              <ul className="divide-y divide-gray-200">
                {tickets?.map((ticket: Ticket) => (
                  <MobileTicketCard
                    key={ticket.id}
                    ticket={ticket}
                    onClick={() => navigate(`/tickets/${ticket.id}`)}
                  />
                ))}
              </ul>
            </div>

            {/* Desktop View */}
            <div className="hidden md:flex flex-col">
              <div className="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div className="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
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
                          <tr
                            key={ticket.id}
                            onClick={() => navigate(`/tickets/${ticket.id}`)}
                            className="cursor-pointer hover:bg-gray-50"
                          >
                            <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                              <StatusBadge status={ticket.status_id} />
                            </td>
                            <td className="whitespace-nowrap px-3 py-4 text-sm font-bold text-gray-900">
                              {ticket.title}
                            </td>
                            <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                              <PriorityLabel priority={ticket.priority_id} />
                            </td>
                            <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                              {ticket.assignee_user_id || 'Unassigned'}
                            </td>
                            <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                              {new Date(ticket.created_at).toLocaleDateString()}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>
          </>
        )}
      </div>

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 z-10 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" onClick={() => setIsModalOpen(false)}></div>
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <form onSubmit={handleSubmit}>
                <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                  <div className="sm:flex sm:items-start">
                    <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                      <h3 className="text-lg leading-6 font-medium text-gray-900" id="modal-title">Create New Ticket</h3>
                      <div className="mt-4 space-y-4">
                        <div>
                          <label htmlFor="title" className="block text-sm font-medium text-gray-700">Title</label>
                          <input
                            type="text"
                            name="title"
                            id="title"
                            required
                            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
                            value={newTicket.title}
                            onChange={(e) => setNewTicket({ ...newTicket, title: e.target.value })}
                          />
                        </div>
                        <div>
                          <label htmlFor="description" className="block text-sm font-medium text-gray-700">Description</label>
                          <textarea
                            name="description"
                            id="description"
                            rows={3}
                            required
                            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
                            value={newTicket.description}
                            onChange={(e) => setNewTicket({ ...newTicket, description: e.target.value })}
                          />
                        </div>
                        <div>
                          <label htmlFor="priority" className="block text-sm font-medium text-gray-700">Priority</label>
                          <select
                            id="priority"
                            name="priority"
                            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
                            value={newTicket.priority_id}
                            onChange={(e) => setNewTicket({ ...newTicket, priority_id: e.target.value })}
                          >
                            <option value="low">Low</option>
                            <option value="medium">Medium</option>
                            <option value="high">High</option>
                            <option value="critical">Critical</option>
                          </select>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Attachments</label>
                            <div className="mt-1 flex items-center">
                                <label htmlFor="file-upload" className="cursor-pointer bg-white py-2 px-3 border border-gray-300 rounded-md shadow-sm text-sm leading-4 font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 flex items-center gap-2">
                                    <Paperclip className="h-4 w-4" />
                                    <span>Upload files</span>
                                    <input id="file-upload" name="file-upload" type="file" className="sr-only" multiple onChange={handleFileChange} />
                                </label>
                                {files && files.length > 0 && (
                                    <span className="ml-3 text-sm text-gray-500">{files.length} file(s) selected</span>
                                )}
                            </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                  <button
                    type="submit"
                    className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
                    disabled={mutation.isPending}
                  >
                    {mutation.isPending ? 'Creating...' : 'Create'}
                  </button>
                  <button
                    type="button"
                    className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                    onClick={() => setIsModalOpen(false)}
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
