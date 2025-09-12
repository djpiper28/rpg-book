import type React from "react";
import { type ReactElement } from "react";
import { Component } from "react";
import { getLogger } from "../../lib/grpcClient/client";
import { H1 } from "../typography/H1";

interface Props {
  children: ReactElement;
  onError: (error: Error) => void;
}

interface State {
  error: Error | undefined;
  hasError: boolean;
}

class ErrorBoundary extends Component<Props, State> {
  static getDerivedStateFromError(error: Error): State {
    return { error, hasError: true };
  }
  constructor(props: Props) {
    super(props);
    this.state = { error: undefined, hasError: false };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
    this.props.onError(error);

    getLogger().error("Uncaught error", {
      error: error.toString(),
      errorInfo: JSON.stringify(errorInfo),
    });
  }

  render(): React.ReactNode {
    if (this.state.hasError) {
      return (
        <div>
          <H1>Something went wrong.</H1>
          <p>Please report this error to the developers.</p>
          <pre>{this.state.error?.toString()}</pre>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
