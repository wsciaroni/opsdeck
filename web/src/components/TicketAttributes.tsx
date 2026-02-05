import { memo } from 'react';
import clsx from 'clsx';

// Optimized: Hoist constants outside component to prevent object recreation on every render.
const STATUS_COLORS: Record<string, string> = {
  new: 'bg-blue-100 text-blue-800',
  in_progress: 'bg-yellow-100 text-yellow-800',
  on_hold: 'bg-orange-100 text-orange-800',
  done: 'bg-green-100 text-green-800',
  canceled: 'bg-gray-100 text-gray-800',
};

const STATUS_LABELS: Record<string, string> = {
  new: 'New',
  in_progress: 'In Progress',
  on_hold: 'On Hold',
  done: 'Done',
  canceled: 'Canceled',
};

// Optimized: Memoize component to prevent re-renders in large lists when props are stable.
export const StatusBadge = memo(function StatusBadge({ status }: { status: string }) {
  return (
    <span className={clsx("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium", STATUS_COLORS[status] || 'bg-gray-100 text-gray-800')}>
      {STATUS_LABELS[status] || status}
    </span>
  );
});

const PRIORITY_COLORS: Record<string, string> = {
  low: 'text-gray-600',
  medium: 'text-blue-700',
  high: 'text-orange-700 font-bold',
  critical: 'text-red-700 font-bold uppercase',
};

const PRIORITY_LABELS: Record<string, string> = {
  low: 'Low',
  medium: 'Medium',
  high: 'High',
  critical: 'Critical',
};

// Optimized: Memoize component to prevent re-renders in large lists when props are stable.
export const PriorityLabel = memo(function PriorityLabel({ priority }: { priority: string }) {
  return <span className={clsx("text-sm", PRIORITY_COLORS[priority] || 'text-gray-500')}>{PRIORITY_LABELS[priority] || priority}</span>;
});
