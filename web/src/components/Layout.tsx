import { useState, type ReactNode } from 'react';
import { useAuth } from '../context/AuthContext';
import ProfileDropdown from './ProfileDropdown';
import { Link, useLocation } from 'react-router-dom';
import clsx from 'clsx';
import { Menu, X } from 'lucide-react';

export default function Layout({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const location = useLocation();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const isCurrent = (path: string) => location.pathname === path;

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white border-b border-gray-200 relative z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex">
              <div className="flex-shrink-0 flex items-center">
                <Link to="/" className="text-xl font-bold text-indigo-600">OpsDeck</Link>
              </div>
              {user && (
                <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                  <Link
                    to="/"
                    className={clsx(
                      isCurrent('/')
                        ? 'border-indigo-500 text-gray-900'
                        : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                      'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium'
                    )}
                  >
                    Dashboard
                  </Link>
                  <Link
                    to="/scheduled-tasks"
                    className={clsx(
                      isCurrent('/scheduled-tasks')
                        ? 'border-indigo-500 text-gray-900'
                        : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                      'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium'
                    )}
                  >
                    Scheduled Tasks
                  </Link>
                </div>
              )}
            </div>
            <div className="flex items-center">
              {user && (
                <>
                  <ProfileDropdown />
                  <div className="flex items-center sm:hidden ml-4">
                    <button
                      onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                      className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500"
                    >
                      <span className="sr-only">Open main menu</span>
                      {isMobileMenuOpen ? (
                        <X className="block h-6 w-6" aria-hidden="true" />
                      ) : (
                        <Menu className="block h-6 w-6" aria-hidden="true" />
                      )}
                    </button>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>

        {/* Mobile menu, show/hide based on menu state. */}
        {user && isMobileMenuOpen && (
          <div className="sm:hidden bg-white border-b border-gray-200">
            <div className="pt-2 pb-3 space-y-1">
              <Link
                to="/"
                className={clsx(
                  isCurrent('/')
                    ? 'bg-indigo-50 border-indigo-500 text-indigo-700'
                    : 'border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700',
                  'block pl-3 pr-4 py-2 border-l-4 text-base font-medium'
                )}
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Dashboard
              </Link>
              <Link
                to="/scheduled-tasks"
                className={clsx(
                  isCurrent('/scheduled-tasks')
                    ? 'bg-indigo-50 border-indigo-500 text-indigo-700'
                    : 'border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700',
                  'block pl-3 pr-4 py-2 border-l-4 text-base font-medium'
                )}
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Scheduled Tasks
              </Link>
            </div>
          </div>
        )}
      </nav>
      <main>
        {children}
      </main>
    </div>
  );
}
