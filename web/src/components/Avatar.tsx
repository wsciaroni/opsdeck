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

// Define the avatar service URL as a constant
const AVATAR_SERVICE_URL = 'https://ui-avatars.com/api/';

export default function Avatar({ src, alt, name, className }: AvatarProps) {
  const [imgError, setImgError] = useState(false);
  const [uiAvatarError, setUiAvatarError] = useState(false);

  const showOriginal = src && !imgError;
  const showUiAvatar = !showOriginal && name && !uiAvatarError;

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

  if (showUiAvatar) {
    // Use initials to avoid sending full PII to third-party service
    const initials = getInitials(name!);

    // Construct URL safely using URL API
    const url = new URL(AVATAR_SERVICE_URL);
    url.searchParams.set('name', initials);
    url.searchParams.set('background', 'random');

    return (
        <img
            className={clsx("object-cover", className)}
            src={url.toString()}
            alt={alt || name}
            onError={() => setUiAvatarError(true)}
            referrerPolicy="no-referrer"
        />
    )
  }

  return (
    <div className={clsx("flex items-center justify-center bg-gray-100 text-gray-400 overflow-hidden", className)}>
      <UserIcon className="h-3/5 w-3/5" />
    </div>
  );
}
