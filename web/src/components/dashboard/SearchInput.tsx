import { useState, useEffect, memo } from 'react';
import { Search, X } from 'lucide-react';

interface SearchInputProps {
  onSearch: (query: string) => void;
  placeholder?: string;
  initialValue?: string;
}

const SearchInput = memo(function SearchInput({ onSearch, placeholder = "Search...", initialValue = "" }: SearchInputProps) {
  const [inputValue, setInputValue] = useState(initialValue);

  useEffect(() => {
    const timer = setTimeout(() => {
      onSearch(inputValue);
    }, 300);
    return () => clearTimeout(timer);
  }, [inputValue, onSearch]);

  return (
    <div className="relative rounded-md shadow-sm w-full sm:w-64">
      <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
        <Search className="h-4 w-4 text-gray-400" aria-hidden="true" />
      </div>
      <label htmlFor="search" className="sr-only">Search</label>
      <input
        type="text"
        name="search"
        id="search"
        className="block w-full rounded-md border-gray-300 pl-10 pr-10 focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm py-2 border"
        placeholder={placeholder}
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
      />
      {inputValue && (
        <div className="absolute inset-y-0 right-0 flex items-center pr-3">
          <button
            type="button"
            onClick={() => setInputValue('')}
            className="text-gray-400 hover:text-gray-500 focus:outline-none focus:text-gray-500"
            aria-label="Clear search"
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </button>
        </div>
      )}
    </div>
  );
});

export default SearchInput;
