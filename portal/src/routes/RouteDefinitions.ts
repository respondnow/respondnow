export function normalizePath(url: string): string {
  return url.replace(/\/{2,}/g, '/');
}

export interface UseRouteDefinitionsProps {
  toRoot(): string;
  toLogin(): string;
  toPasswordReset(): string;
  toGetStarted(): string;
  toIncidentDashboard(): string;
}

export const paths: UseRouteDefinitionsProps = {
  toRoot: () => '/',
  toLogin: () => '/login',
  toPasswordReset: () => '/settings/password-reset',
  toGetStarted: () => '/get-started',
  toIncidentDashboard: () => '/incidents'
};
