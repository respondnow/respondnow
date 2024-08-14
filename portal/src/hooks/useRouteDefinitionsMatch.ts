import React from 'react';
import { mapValues } from 'lodash-es';
import { type UseRouteDefinitionsProps, paths, normalizePath } from '@routes/RouteDefinitions';
import { useAppStore } from './useAppStore';

export function useRouteDefinitionsMatch(): UseRouteDefinitionsProps {
  const { matchPath } = useAppStore();

  return React.useMemo(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    // () => mapValues(paths, route => (params?: any) => normalizePath(`${matchPath}/${route(params)}`)),
    () => mapValues(paths, route => () => normalizePath(`${matchPath}/${route()}`)),

    [matchPath]
  );
}
