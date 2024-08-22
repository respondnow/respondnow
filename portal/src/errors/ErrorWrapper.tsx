import React from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { Fallback } from './Fallback';

interface ErrorWrapperProps {
  children: React.ReactNode | undefined;
}

export function ErrorWrapper({ children }: ErrorWrapperProps): React.ReactElement {
  return children ? <ErrorBoundary FallbackComponent={Fallback}>{children}</ErrorBoundary> : <></>;
}
