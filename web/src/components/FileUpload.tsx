import { Paperclip, X } from 'lucide-react';
import toast from 'react-hot-toast';
import { formatBytes } from '../utils';

interface FileUploadProps {
  files: File[];
  onFilesChange: (files: File[]) => void;
  maxSize?: number;
  allowedExtensions?: string[];
}

export default function FileUpload({
  files,
  onFilesChange,
  maxSize = 32 * 1024 * 1024, // 32MB
  allowedExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.pdf', '.doc', '.docx', '.xls', '.xlsx', '.csv', '.txt']
}: FileUploadProps) {
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const selectedFiles = Array.from(e.target.files);
      const validFiles: File[] = [];

      selectedFiles.forEach(file => {
        if (file.size > maxSize) {
          toast.error(`File ${file.name} is too large (max ${formatBytes(maxSize)})`);
          return;
        }

        const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase();
        if (!allowedExtensions.includes(fileExtension) && file.name.includes('.')) {
             toast.error(`File type ${fileExtension} is not allowed`);
             return;
        }

        validFiles.push(file);
      });

      if (validFiles.length > 0) {
        onFilesChange([...files, ...validFiles]);
      }

      // Reset input value to allow re-selecting the same file if needed
      e.target.value = '';
    }
  };

  const removeFile = (indexToRemove: number) => {
    onFilesChange(files.filter((_, index) => index !== indexToRemove));
  };

  return (
    <div>
      <label htmlFor="file-upload" className="block text-sm font-medium text-gray-700">Attachments</label>
      <div className="mt-1 flex items-center">
        <label htmlFor="file-upload" className="cursor-pointer bg-white py-2 px-3 border border-gray-300 rounded-md shadow-sm text-sm leading-4 font-medium text-gray-700 hover:bg-gray-50 focus-within:outline-none focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-indigo-500 flex items-center gap-2">
          <Paperclip className="h-4 w-4" />
          <span>Upload files</span>
          <input
            id="file-upload"
            name="file-upload"
            type="file"
            className="sr-only"
            multiple
            accept={allowedExtensions.join(',')}
            onChange={handleFileChange}
          />
        </label>
      </div>
      {files.length > 0 && (
        <ul className="mt-3 space-y-1">
          {files.map((file, index) => (
            <li key={index} className="text-sm text-gray-500 flex items-center justify-between py-1">
              <div className="flex items-center min-w-0">
                <Paperclip className="h-3 w-3 mr-2 text-gray-500 flex-shrink-0" />
                <span className="truncate">
                  {file.name} <span className="text-gray-600 text-xs ml-1">({formatBytes(file.size)})</span>
                </span>
              </div>
              <button
                type="button"
                onClick={() => removeFile(index)}
                className="ml-2 text-gray-500 hover:text-gray-700 focus:outline-none focus:text-gray-800 focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-indigo-500 rounded-sm"
                aria-label={`Remove ${file.name}`}
                title={`Remove ${file.name}`}
              >
                <X className="h-4 w-4" />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
