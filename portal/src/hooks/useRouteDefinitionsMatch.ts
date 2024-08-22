import React from 'react';
import { mapValues } from 'lodash-es';
import { type UseRouteDefinitionsProps, paths, normalizePath } from '@routes/RouteDefinitions';

export function useRouteDefinitionsMatch(): UseRouteDefinitionsProps {
  return React.useMemo(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    () => mapValues(paths, route => () => normalizePath(`/${route()}`)),
    []
  );
}
