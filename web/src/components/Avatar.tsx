import { useState } from 'react';
import { User as UserIcon } from 'lucide-react';
import clsx from 'clsx';

interface AvatarProps {
  src?: string;
  alt?: string;
  name?: string;
  className?: string;
}

const getInitials = (name: string) => {
  return name
    .trim()
    .split(/\s+/)
    .map(part => part[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);
};

const getBackgroundColor = (name: string) => {
  const colors = [
    'bg-red-500',
    'bg-orange-500',
    'bg-amber-500',
    'bg-yellow-500',
    'bg-lime-500',
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
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  return colors[Math.abs(hash) % colors.length];
};

export default function Avatar({ src, alt, name, className }: AvatarProps) {
  const [imgError, setImgError] = useState(false);

  const showOriginal = src && !imgError;
  // If no image or image failed, and we have a name, show initials
  const showInitials = !showOriginal && name;

  if (showOriginal) {
    return (
      <img
        className={clsx("object-cover", className)}
        src={src}
        alt={alt || ""}
        onError={() => setImgError(true)}
        referrerPolicy="no-referrer"
      />
    );
  }

  if (showInitials) {
    const initials = getInitials(name!);
    const bgColor = getBackgroundColor(name!);

    return (
      <div
        className={clsx(
          "flex items-center justify-center text-white font-medium select-none",
          bgColor,
          className
        )}
        aria-label={alt || name}
      >
        <span className="text-[40%] leading-none">{initials}</span>
      </div>
    );
  }

  return (
    <div className={clsx("flex items-center justify-center bg-gray-100 text-gray-400 overflow-hidden", className)}>
      <UserIcon className="h-3/5 w-3/5" />
    </div>
  );
}
