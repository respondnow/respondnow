import React from 'react';
import { mapValues } from 'lodash-es';
import { type UseRouteDefinitionsProps, paths, normalizePath } from '@routes/RouteDefinitions';

export function useRouteWithBaseUrl(): UseRouteDefinitionsProps {
  return React.useMemo(() => mapValues(paths, route => () => normalizePath(`/${route()}`)), []);
}
