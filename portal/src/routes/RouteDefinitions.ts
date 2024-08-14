export function normalizePath(url: string): string {
  return url.replace(/\/{2,}/g, '/');
}

export interface UseRouteDefinitionsProps {
  toRoot(): string;
  toLogin(): string;
}

export const paths: UseRouteDefinitionsProps = {
  toRoot: () => '/',
  toLogin: () => '/login'
};
