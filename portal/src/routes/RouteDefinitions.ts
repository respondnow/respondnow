export function normalizePath(url: string): string {
  return url.replace(/\/{2,}/g, '/');
}

export interface UseRouteDefinitionsProps {
  toRoot(): string;
  toLogin(): string;
  toPasswordReset(): string;
  toGetStarted(): string;
  toIncidentDashboard(): string;
  toIncidentDetails(params: { incidentId: string }): string;
  // TEMP
  toIncidentDetailsDummy(): string;
}

export const paths: UseRouteDefinitionsProps = {
  toRoot: () => '/',
  toLogin: () => '/login',
  toPasswordReset: () => '/settings/password-reset',
  toGetStarted: () => '/getting-started',
  toIncidentDashboard: () => '/incidents',
  toIncidentDetails: ({ incidentId }) => `/incidents/${incidentId}`,
  // TEMP
  toIncidentDetailsDummy: () => '/incident'
};

export interface IncidentDetailsPathProps {
  incidentId: string;
}
