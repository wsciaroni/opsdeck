import { type ReactNode } from 'react';
import { useAuth } from '../context/AuthContext';
import ProfileDropdown from './ProfileDropdown';

export default function Layout({ children }: { children: ReactNode }) {
  const { user } = useAuth();

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white border-b border-gray-200 relative z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex">
              <div className="flex-shrink-0 flex items-center">
                <span className="text-xl font-bold text-indigo-600">OpsDeck</span>
              </div>
            </div>
            <div className="flex items-center">
              {user && (
                <ProfileDropdown />
              )}
            </div>
          </div>
        </div>
      </nav>
      <main>
        {children}
      </main>
    </div>
  );
}
