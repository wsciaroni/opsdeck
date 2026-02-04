import clsx from 'clsx';
import { useMemo } from 'react';

interface AvatarProps {
  name: string;
  src?: string | null;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
}

const COLORS = [
  'bg-red-500',
  'bg-orange-500',
  'bg-amber-500',
  'bg-green-500',
  'bg-emerald-500',
  'bg-teal-500',
  'bg-cyan-500',
  'bg-sky-500',
  'bg-blue-500',
  'bg-indigo-500',
  'bg-violet-500',
  'bg-purple-500',
  'bg-fuchsia-500',
  'bg-pink-500',
  'bg-rose-500',
];

function getInitials(name: string): string {
  const parts = name.trim().split(/\s+/);
  if (parts.length === 0) return '?';
  if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

function getColor(name: string): string {
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  return COLORS[Math.abs(hash) % COLORS.length];
}

export default function Avatar({ name, src, className, size = 'md' }: AvatarProps) {
  const initials = useMemo(() => getInitials(name), [name]);
  const colorClass = useMemo(() => getColor(name), [name]);

  const sizeClasses = {
    sm: 'h-8 w-8 text-xs',
    md: 'h-10 w-10 text-sm',
    lg: 'h-12 w-12 text-base',
  };

  if (src) {
    return (
      <img
        className={clsx('rounded-full object-cover', sizeClasses[size], className)}
        src={src}
        alt={name}
      />
    );
  }

  return (
    <div
      className={clsx(
        'rounded-full flex items-center justify-center text-white font-medium select-none',
        sizeClasses[size],
        colorClass,
        className
      )}
      aria-hidden="true"
    >
      {initials}
    </div>
  );
}
