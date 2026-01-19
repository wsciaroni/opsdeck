import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import { client } from '../api/client';
import type { User, Organization } from '../types';

interface AuthContextType {
  user: User | null;
  organizations: Organization[];
  currentOrg: Organization | null;
  isLoading: boolean;
  login: () => void;
  logout: () => void;
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

  const logout = () => {
    // In a real app, we might call an API endpoint to clear the cookie.
    // For now, we'll just clear state and maybe redirect.
    // Assuming there is a backend logout endpoint would be better, but the task didn't specify one.
    // However, simply redirecting to home or clearing state is what we can do.
    // Since session is cookie based, we really should have a logout endpoint.
    // But per task requirements: "login() function: window.location.href = '/auth/login'."
    // Logout wasn't strictly detailed but "logout" was listed in context.
    setUser(null);
    setOrganizations([]);
    setCurrentOrg(null);
    // Optionally: window.location.href = '/auth/logout'; // if it existed
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
