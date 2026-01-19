import { AuthProvider, useAuth } from './context/AuthContext';
import './App.css';

function Dashboard() {
  const { user, currentOrg, login } = useAuth();

  if (!user) {
    return (
      <div className="login-container">
        <h1>OpsDeck</h1>
        <button onClick={login} className="login-button">
          Login with Google
        </button>
      </div>
    );
  }

  return (
    <div className="dashboard-container">
      <h1>Welcome, {user.name}!</h1>
      {currentOrg ? (
        <p>Current Workspace: <strong>{currentOrg.name}</strong></p>
      ) : (
        <p>No workspace selected.</p>
      )}
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <MainContent />
    </AuthProvider>
  );
}

function MainContent() {
  const { isLoading } = useAuth();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return <Dashboard />;
}

export default App;
