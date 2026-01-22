import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createTicket } from '../../api/tickets';
import toast from 'react-hot-toast';

interface CreateTicketModalProps {
  isOpen: boolean;
  onClose: () => void;
  organizationId: string;
}

export default function CreateTicketModal({ isOpen, onClose, organizationId }: CreateTicketModalProps) {
  const queryClient = useQueryClient();
  const [newTicket, setNewTicket] = useState({
    title: '',
    description: '',
    priority_id: 'medium',
  });

  const mutation = useMutation({
    mutationFn: createTicket,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tickets'] });
      onClose();
      setNewTicket({ title: '', description: '', priority_id: 'medium' });
      toast.success("Ticket created!");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!organizationId) return;
    mutation.mutate({
      ...newTicket,
      organization_id: organizationId,
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-10 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" onClick={onClose}></div>
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
                onClick={onClose}
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
