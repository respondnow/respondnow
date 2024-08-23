import React from 'react';
import { StringsContextProvider } from '@strings';
import { AppStoreProvider, ReactQueryProvider } from '@context';
import { Routes } from '@routes/RouteDestinations';
import strings from 'strings/strings.en.yaml';

export function App(): React.ReactElement {
  return (
    <AppStoreProvider>
      <StringsContextProvider data={strings}>
        <ReactQueryProvider>
          <Routes />
        </ReactQueryProvider>
      </StringsContextProvider>
    </AppStoreProvider>
  );
}
