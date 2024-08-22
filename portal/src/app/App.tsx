import React from 'react';
import { StringsContextProvider } from '@strings';
import { AppStoreProvider, ReactQueryProvider } from '@context';
import { RoutesWithAuthentication, RoutesWithoutAuthentication } from '@routes/RouteDestinations';
import strings from 'strings/strings.en.yaml';

export function AppWithAuthentication(): React.ReactElement {
  return (
    <AppStoreProvider scope={{}} updateAppStore={() => void 0}>
      <StringsContextProvider data={strings}>
        <ReactQueryProvider>
          <RoutesWithAuthentication />
        </ReactQueryProvider>
      </StringsContextProvider>
    </AppStoreProvider>
  );
}

export function AppWithoutAuthentication(): React.ReactElement {
  return (
    <StringsContextProvider data={strings}>
      <ReactQueryProvider>
        <RoutesWithoutAuthentication />
      </ReactQueryProvider>
    </StringsContextProvider>
  );
}
