import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { listScheduledTasks, deleteScheduledTask } from '../api/scheduled_tasks';
import { Plus, Edit2, Trash2, Calendar, Repeat } from 'lucide-react';
import { format } from 'date-fns';
import CreateScheduledTaskModal from '../components/scheduled_tasks/CreateScheduledTaskModal';
import type { ScheduledTask } from '../types';

export default function ScheduledTasks() {
    const { currentOrg } = useAuth();
    const queryClient = useQueryClient();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingTask, setEditingTask] = useState<ScheduledTask | undefined>(undefined);

    const { data: tasks, isLoading, error } = useQuery({
        queryKey: ['scheduledTasks', currentOrg?.id],
        queryFn: () => listScheduledTasks(currentOrg!.id),
        enabled: !!currentOrg,
    });

    const deleteMutation = useMutation({
        mutationFn: deleteScheduledTask,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['scheduledTasks', currentOrg?.id] });
        },
    });

    const handleEdit = (task: ScheduledTask) => {
        setEditingTask(task);
        setIsModalOpen(true);
    };

    const handleDelete = async (id: string) => {
        if (confirm('Are you sure you want to delete this scheduled task?')) {
            await deleteMutation.mutateAsync(id);
        }
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setEditingTask(undefined);
    };

    if (!currentOrg) return <div>Please select an organization</div>;

    return (
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">Scheduled Tasks</h1>
                <button
                    onClick={() => setIsModalOpen(true)}
                    className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                    <Plus className="h-4 w-4 mr-2" />
                    New Task
                </button>
            </div>

            {isLoading ? (
                <div>Loading...</div>
            ) : error ? (
                <div>Error loading tasks</div>
            ) : tasks?.length === 0 ? (
                <div className="text-center py-12 bg-white rounded-lg shadow">
                    <Calendar className="mx-auto h-12 w-12 text-gray-400" />
                    <h3 className="mt-2 text-sm font-medium text-gray-900">No scheduled tasks</h3>
                    <p className="mt-1 text-sm text-gray-500">Get started by creating a new recurring task.</p>
                </div>
            ) : (
                <div className="bg-white shadow overflow-hidden sm:rounded-md">
                    <ul className="divide-y divide-gray-200">
                        {tasks?.map((task: ScheduledTask) => (
                            <li key={task.id}>
                                <div className="px-4 py-4 sm:px-6 hover:bg-gray-50 flex justify-between items-center">
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-center justify-between">
                                            <p className="text-sm font-medium text-indigo-600 truncate">{task.title}</p>
                                            <div className="ml-2 flex-shrink-0 flex">
                                                <p className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${task.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                                                    {task.enabled ? 'Active' : 'Disabled'}
                                                </p>
                                            </div>
                                        </div>
                                        <div className="mt-2 sm:flex sm:justify-between">
                                            <div className="sm:flex">
                                                <p className="flex items-center text-sm text-gray-500 mr-6">
                                                    <Repeat className="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-400" />
                                                    <span className="capitalize">{task.frequency}</span>
                                                </p>
                                                <p className="mt-2 flex items-center text-sm text-gray-500 sm:mt-0 sm:ml-6">
                                                    <Calendar className="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-400" />
                                                    Next: {format(new Date(task.next_run_at), 'MMM d, yyyy h:mm a')}
                                                </p>
                                            </div>
                                            <div className="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                                                <button onClick={() => handleEdit(task)} className="text-gray-400 hover:text-gray-600 mr-4">
                                                    <Edit2 className="h-5 w-5" />
                                                </button>
                                                <button onClick={() => handleDelete(task.id)} className="text-red-400 hover:text-red-600">
                                                    <Trash2 className="h-5 w-5" />
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </li>
                        ))}
                    </ul>
                </div>
            )}

            <CreateScheduledTaskModal
                isOpen={isModalOpen}
                onClose={handleCloseModal}
                organizationId={currentOrg.id}
                initialData={editingTask}
            />
        </div>
    );
}
