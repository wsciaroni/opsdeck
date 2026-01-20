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
  switchOrganization: (orgID: string) => void;
  refreshOrganizations: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [currentOrg, setCurrentOrg] = useState<Organization | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const fetchMe = async () => {
    try {
      const response = await client.get('/me');
      const data = response.data;
      setUser(data.user);
      const orgs = data.organizations || [];
      setOrganizations(orgs);

      // Determine which org to select
      if (orgs.length > 0) {
        const lastOrgID = localStorage.getItem('last_org_id');
        if (lastOrgID) {
          const found = orgs.find((o: Organization) => o.id === lastOrgID);
          if (found) {
            setCurrentOrg(found);
            return;
          }
        }
        // Default to first if no valid stored org or not found
        // Only change currentOrg if it's currently null or not in the new list (which implies a refresh)
        // Actually, if we are refreshing, we might want to keep the current one if it still exists.
        // But for initial load (when currentOrg is null), this works.
        // If we call refreshOrganizations(), we should try to keep the current one.
        setCurrentOrg((prev) => {
           if (prev && orgs.find((o: Organization) => o.id === prev.id)) {
             return prev;
           }
           return orgs[0];
        });
      } else {
        setCurrentOrg(null);
      }
    } catch (error) {
      console.error("Failed to fetch user", error);
      setUser(null);
    }
  };

  useEffect(() => {
    fetchMe().finally(() => setIsLoading(false));
  }, []);

  const switchOrganization = (orgID: string) => {
    const org = organizations.find((o) => o.id === orgID);
    if (org) {
      setCurrentOrg(org);
      localStorage.setItem('last_org_id', orgID);
    }
  };

  const refreshOrganizations = async () => {
    await fetchMe();
  };

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
    <AuthContext.Provider value={{ user, organizations, currentOrg, isLoading, login, logout, switchOrganization, refreshOrganizations }}>
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
