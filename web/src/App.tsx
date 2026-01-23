import { AuthProvider, useAuth } from './context/AuthContext';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Routes, Route, Navigate, Outlet } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import TicketDetail from './pages/TicketDetail';
import TeamSettings from './pages/TeamSettings';
import PublicTicketSubmit from './pages/PublicTicketSubmit';
import PublicDashboard from './pages/PublicDashboard';
import PublicTicketDetail from './pages/PublicTicketDetail';
import Profile from './pages/Profile';
import Login from './pages/Login';
import Layout from './components/Layout';
import { type ReactNode } from 'react';
import { Toaster } from 'react-hot-toast';
import './App.css';

const queryClient = new QueryClient();

function RequireAuth({ children }: { children?: ReactNode }) {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return <div className="flex justify-center items-center h-screen">Loading...</div>;
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return children || <Outlet />;
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Toaster position="bottom-right" reverseOrder={false} />
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/submit-ticket" element={<PublicTicketSubmit />} />
          <Route path="/public/:token" element={<PublicDashboard />} />
          <Route path="/public/:token/tickets/:ticketId" element={<PublicTicketDetail />} />
          <Route element={<RequireAuth><Layout><Outlet /></Layout></RequireAuth>}>
            <Route path="/" element={<Dashboard />} />
            <Route path="/tickets/:id" element={<TicketDetail />} />
            <Route path="/organizations/:orgId/settings/team" element={<TeamSettings />} />
            <Route path="/profile" element={<Profile />} />
          </Route>
        </Routes>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
