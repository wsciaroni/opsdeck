import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createScheduledTask, updateScheduledTask } from '../../api/scheduled_tasks';
import toast from 'react-hot-toast';
import { Dialog, DialogPanel, DialogTitle, Transition, TransitionChild } from '@headlessui/react';
import { Loader2 } from 'lucide-react';
import { FREQUENCIES, type ScheduledTask } from '../../types';

interface CreateScheduledTaskModalProps {
  isOpen: boolean;
  onClose: () => void;
  organizationId: string;
  initialData?: ScheduledTask;
}

export default function CreateScheduledTaskModal({ isOpen, onClose, organizationId, initialData }: CreateScheduledTaskModalProps) {
  const queryClient = useQueryClient();

  // Initialize state directly from props.
  // This component is expected to be re-mounted (via key change in parent)
  // whenever it is opened or the target task changes.
  const [formData, setFormData] = useState(() => {
    if (initialData) {
        return {
            title: initialData.title,
            description: initialData.description,
            priority_id: initialData.priority_id,
            frequency: initialData.frequency,
            start_date: new Date(initialData.next_run_at).toISOString().split('T')[0],
            location: initialData.location,
            enabled: initialData.enabled,
        };
    }
    return {
        title: '',
        description: '',
        priority_id: 'medium',
        frequency: 'daily',
        start_date: new Date().toISOString().split('T')[0],
        location: '',
        enabled: true,
    };
  });

  const createMutation = useMutation({
    mutationFn: createScheduledTask,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['scheduledTasks'] });
      onClose();
      toast.success("Scheduled Task created!");
    },
  });

  const updateMutation = useMutation({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    mutationFn: (data: any) => updateScheduledTask(initialData!.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['scheduledTasks'] });
      onClose();
      toast.success("Scheduled Task updated!");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!organizationId) return;

    const payload = {
        ...formData,
        organization_id: organizationId,
        start_date: new Date(formData.start_date).toISOString(),
        next_run_at: new Date(formData.start_date).toISOString(),
    };

    if (initialData) {
        updateMutation.mutate(payload);
    } else {
        createMutation.mutate(payload);
    }
  };

  return (
    <Transition show={isOpen}>
      <Dialog className="relative z-10" onClose={onClose}>
        <TransitionChild
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        </TransitionChild>

        <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <TransitionChild
              enter="ease-out duration-300"
              enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
              enterTo="opacity-100 translate-y-0 sm:scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 translate-y-0 sm:scale-100"
              leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            >
              <DialogPanel className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg">
                <form onSubmit={handleSubmit}>
                  <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                    <div className="sm:flex sm:items-start">
                      <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                        <DialogTitle as="h3" className="text-lg leading-6 font-medium text-gray-900">
                            {initialData ? 'Edit Scheduled Task' : 'Create Scheduled Task'}
                        </DialogTitle>
                        <div className="mt-4 space-y-4">
                          <div>
                            <label htmlFor="title" className="block text-sm font-medium text-gray-700">Title</label>
                            <input
                              type="text"
                              name="title"
                              id="title"
                              required
                              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
                              value={formData.title}
                              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                            />
                          </div>
                          <div>
                            <label htmlFor="description" className="block text-sm font-medium text-gray-700">Description</label>
                            <textarea
                              name="description"
                              id="description"
                              rows={3}
                              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
                              value={formData.description}
                              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                            />
                          </div>
                          <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label htmlFor="frequency" className="block text-sm font-medium text-gray-700">Frequency</label>
                                <select
                                id="frequency"
                                name="frequency"
                                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
                                value={formData.frequency}
                                onChange={(e) => setFormData({ ...formData, frequency: e.target.value })}
                                >
                                {FREQUENCIES.map((freq) => (
                                    <option key={freq.id} value={freq.id}>{freq.label}</option>
                                ))}
                                </select>
                            </div>
                            <div>
                                <label htmlFor="start_date" className="block text-sm font-medium text-gray-700">Next Run Date</label>
                                <input
                                type="date"
                                name="start_date"
                                id="start_date"
                                required
                                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
                                value={formData.start_date}
                                onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                                />
                            </div>
                          </div>
                          <div>
                            <label htmlFor="priority" className="block text-sm font-medium text-gray-700">Priority</label>
                            <select
                              id="priority"
                              name="priority"
                              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
                              value={formData.priority_id}
                              onChange={(e) => setFormData({ ...formData, priority_id: e.target.value })}
                            >
                              <option value="low">Low</option>
                              <option value="medium">Medium</option>
                              <option value="high">High</option>
                              <option value="critical">Critical</option>
                            </select>
                          </div>
                          <div>
                            <label htmlFor="location" className="block text-sm font-medium text-gray-700">Location</label>
                            <input
                              type="text"
                              name="location"
                              id="location"
                              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm border p-2"
                              value={formData.location}
                              onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                            />
                          </div>
                          <div className="flex items-center">
                            <input
                              id="enabled"
                              name="enabled"
                              type="checkbox"
                              className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                              checked={formData.enabled}
                              onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
                            />
                            <label htmlFor="enabled" className="ml-2 block text-sm text-gray-900">
                              Enabled
                            </label>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                    <button
                      type="submit"
                      className="w-full inline-flex justify-center items-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                      disabled={createMutation.isPending || updateMutation.isPending}
                    >
                      {createMutation.isPending || updateMutation.isPending ? (
                        <>
                          <Loader2 className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" />
                          Saving...
                        </>
                      ) : (
                        initialData ? 'Save Changes' : 'Create Task'
                      )}
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
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}
