import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import { client } from '../api/client';
import { logout as apiLogout } from '../api/auth';
import type { User, Organization } from '../types';

interface AuthContextType {
  user: User | null;
  organizations: Organization[];
  currentOrg: Organization | null;
  isLoading: boolean;
  login: () => void;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [currentOrg, setCurrentOrg] = useState<Organization | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function fetchMe() {
      try {
        const response = await client.get('/me');
        const data = response.data;
        setUser(data.user);
        setOrganizations(data.organizations || []);
        if (data.organizations && data.organizations.length > 0) {
          setCurrentOrg(data.organizations[0]);
        }
      } catch (error) {
        // If 401, we are just not logged in.
        // For other errors, we might want to log them or show a message,
        // but for now we'll just assume not logged in / user is null.
        console.error("Failed to fetch user", error);
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    }

    fetchMe();
  }, []);

  const login = () => {
    window.location.href = '/auth/login';
  };

  const logout = async () => {
    try {
      await apiLogout();
    } catch (error) {
      console.error('Logout failed', error);
    } finally {
      setUser(null);
      setOrganizations([]);
      setCurrentOrg(null);
      // We don't need to force reload, the state change will show the login screen.
      // But if we want to be sure all cookies are cleared in the browser view:
      // window.location.href = '/';
      // However, since we are SPA, state clearing is enough for UI.
      // If the backend cookie is gone, next refresh will also show login screen.
    }
  };

  return (
    <AuthContext.Provider value={{ user, organizations, currentOrg, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
