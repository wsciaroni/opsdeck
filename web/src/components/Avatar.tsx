import { useState } from 'react';
import { User as UserIcon } from 'lucide-react';
import clsx from 'clsx';

interface AvatarProps {
  src?: string;
  alt?: string;
  name?: string;
  className?: string;
}

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
      />
    );
  }

  if (showUiAvatar) {
    return (
        <img
            className={clsx("object-cover", className)}
            src={`https://ui-avatars.com/api/?name=${encodeURIComponent(name!)}&background=random`}
            alt={alt || name}
            onError={() => setUiAvatarError(true)}
        />
    )
  }

  return (
    <div className={clsx("flex items-center justify-center bg-gray-100 text-gray-400 overflow-hidden", className)}>
      <UserIcon className="h-3/5 w-3/5" />
    </div>
  );
}
