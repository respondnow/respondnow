import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { StringsContextProvider } from '@strings';
import { AppStoreProvider, ReactQueryProvider } from '@context';
import { Routes } from '@routes/RouteDestinations';
import strings from 'strings/strings.en.yaml';
import '../styles/global.scss';

ReactDOM.render(
  <BrowserRouter>
    <AppStoreProvider>
      <StringsContextProvider data={strings}>
        <ReactQueryProvider>
          <Routes />
        </ReactQueryProvider>
      </StringsContextProvider>
    </AppStoreProvider>
  </BrowserRouter>,
  document.getElementById('react-root')
);
