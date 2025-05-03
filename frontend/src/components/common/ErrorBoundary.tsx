import { Component, ErrorInfo, ReactNode } from "react";
import { SessionExpiredModal } from "../auth/SessionExpiredModal";

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
  isSessionExpired: boolean;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      isSessionExpired: false,
    };
  }

  static getDerivedStateFromError(error: Error): State {
    // Check if the error is related to authentication
    const isAuthError =
      error.message.includes("401") ||
      error.message.includes("unauthorized") ||
      error.message.includes("expired");

    return {
      hasError: true,
      error,
      isSessionExpired: isAuthError,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    console.error("ErrorBoundary caught an error:", error, errorInfo);
  }

  handleSessionRefresh = () => {
    this.setState({ hasError: false, error: null, isSessionExpired: false });
  };

  render() {
    if (this.state.hasError) {
      if (this.state.isSessionExpired) {
        return <SessionExpiredModal onClose={this.handleSessionRefresh} />;
      }

      // General error UI
      return (
        <div className="error-container">
          <h2>Something went wrong</h2>
          <p>{this.state.error?.message || "An unexpected error occurred"}</p>
          <button onClick={() => this.setState({ hasError: false })}>
            Try Again
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}
