import React from 'react';
import { mapValues } from 'lodash-es';
import { type UseRouteDefinitionsProps, paths, normalizePath } from '@routes/RouteDefinitions';
import { useAppStore } from './useAppStore';

export function useRouteWithBaseUrl(): UseRouteDefinitionsProps {
  const { renderUrl } = useAppStore();

  return React.useMemo(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    // () => mapValues(paths, route => (params?: any) => normalizePath(`${renderUrl}/${route(params)}`)),
    () => mapValues(paths, route => () => normalizePath(`${renderUrl}/${route()}`)),
    [renderUrl]
  );
}
