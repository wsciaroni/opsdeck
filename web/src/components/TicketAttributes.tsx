import clsx from 'clsx';

export function StatusBadge({ status }: { status: string }) {
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
    <span className={clsx("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium", colors[status] || 'bg-gray-100 text-gray-800')}>
      {labels[status] || status}
    </span>
  );
}

export function PriorityLabel({ priority }: { priority: string }) {
  const colors: Record<string, string> = {
    low: 'text-gray-500',
    medium: 'text-blue-500',
    high: 'text-orange-500 font-bold',
    critical: 'text-red-600 font-bold uppercase',
  };

  return <span className={clsx("text-sm", colors[priority] || 'text-gray-500')}>{priority}</span>;
}
