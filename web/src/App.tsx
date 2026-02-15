import { AuthProvider, useAuth } from './context/AuthContext';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Routes, Route, Navigate, Outlet } from 'react-router-dom';
import Layout from './components/Layout';
import DemoBanner from './components/DemoBanner';
import { type ReactNode, Suspense, lazy } from 'react';

// Lazy load pages for code splitting to improve initial load performance
const Dashboard = lazy(() => import('./pages/Dashboard'));
const ScheduledTasks = lazy(() => import('./pages/ScheduledTasks'));
const Reports = lazy(() => import('./pages/Reports'));
const TicketDetail = lazy(() => import('./pages/TicketDetail'));
const TeamSettings = lazy(() => import('./pages/TeamSettings'));
const PublicTicketSubmit = lazy(() => import('./pages/PublicTicketSubmit'));
const PublicDashboard = lazy(() => import('./pages/PublicDashboard'));
const PublicTicketDetail = lazy(() => import('./pages/PublicTicketDetail'));
const Profile = lazy(() => import('./pages/Profile'));
const Login = lazy(() => import('./pages/Login'));
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
        <DemoBanner />
        <Suspense fallback={<div className="flex justify-center items-center h-screen text-gray-500">Loading...</div>}>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/submit-ticket" element={<PublicTicketSubmit />} />
            <Route path="/public/:token" element={<PublicDashboard />} />
            <Route path="/public/:token/tickets/:ticketId" element={<PublicTicketDetail />} />
            <Route element={<RequireAuth><Layout><Outlet /></Layout></RequireAuth>}>
              <Route path="/" element={<Dashboard />} />
              <Route path="/scheduled-tasks" element={<ScheduledTasks />} />
              <Route path="/reports" element={<Reports />} />
              <Route path="/tickets/:id" element={<TicketDetail />} />
              <Route path="/organizations/:orgId/settings/team" element={<TeamSettings />} />
              <Route path="/profile" element={<Profile />} />
            </Route>
          </Routes>
        </Suspense>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
