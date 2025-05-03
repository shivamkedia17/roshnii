import { ErrorBoundary } from "./common/ErrorBoundary";
import { useAuth } from "../context/AuthContext";
import { Login } from "./auth/Login";
import { MainLayout } from "./layout/MainLayout";
import "@/css/App.css";

export default function App() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <ErrorBoundary>
      {isAuthenticated ? <MainLayout /> : <Login />}
    </ErrorBoundary>
  );
}
