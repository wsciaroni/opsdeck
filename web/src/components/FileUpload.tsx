import { Paperclip, X } from 'lucide-react';
import toast from 'react-hot-toast';
import { formatBytes } from '../utils';
import { useState, useCallback } from 'react';
import clsx from 'clsx';

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
  const [isDragging, setIsDragging] = useState(false);

  const processFiles = useCallback((selectedFiles: File[]) => {
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
  }, [files, maxSize, allowedExtensions, onFilesChange]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      processFiles(Array.from(e.target.files));
      e.target.value = '';
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    if (e.currentTarget.contains(e.relatedTarget as Node)) {
      return;
    }
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      processFiles(Array.from(e.dataTransfer.files));
    }
  };

  const removeFile = (indexToRemove: number) => {
    onFilesChange(files.filter((_, index) => index !== indexToRemove));
  };

  return (
    <div>
      <div id="attachments-label" className="block text-sm font-medium text-gray-700 mb-1 flex justify-between items-center">
        <span>Attachments</span>
        {files.length > 1 && (
          <button
            type="button"
            onClick={() => onFilesChange([])}
            className="text-xs text-red-600 hover:text-red-800 transition-colors"
          >
            Clear all
          </button>
        )}
      </div>

      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        className={clsx(
          "relative flex justify-center rounded-md border-2 border-dashed px-6 pt-5 pb-6 transition-colors",
          isDragging ? "border-indigo-500 bg-indigo-50" : "border-gray-300 hover:border-gray-400"
        )}
      >
        <div className="space-y-1 text-center">
          <Paperclip className="mx-auto h-8 w-8 text-gray-400" />
          <div className="flex text-sm text-gray-600 justify-center">
            <label
              htmlFor="file-upload"
              className="relative cursor-pointer rounded-md bg-white font-medium text-indigo-600 focus-within:outline-none focus-within:ring-2 focus-within:ring-indigo-500 focus-within:ring-offset-2 hover:text-indigo-500"
            >
              <span>Upload a file</span>
              <input
                id="file-upload"
                name="file-upload"
                type="file"
                className="sr-only"
                multiple
                accept={allowedExtensions.join(',')}
                onChange={handleFileChange}
                aria-labelledby="attachments-label"
              />
            </label>
            <p className="pl-1">or drag and drop</p>
          </div>
          <p className="text-xs text-gray-500">
            {allowedExtensions.join(', ')} up to {formatBytes(maxSize)}
          </p>
        </div>
      </div>

      {files.length > 0 && (
        <ul className="mt-3 space-y-1">
          {files.map((file, index) => (
            <li key={index} className="text-sm text-gray-500 flex items-center justify-between py-1 bg-gray-50 px-2 rounded">
              <div className="flex items-center min-w-0">
                <Paperclip className="h-3 w-3 mr-2 text-gray-500 flex-shrink-0" />
                <span className="truncate">
                  {file.name} <span className="text-gray-600 text-xs ml-1">({formatBytes(file.size)})</span>
                </span>
              </div>
              <button
                type="button"
                onClick={() => removeFile(index)}
                className="ml-2 text-gray-500 hover:text-red-600 focus:outline-none focus:text-red-700 focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-red-500 rounded-sm"
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
