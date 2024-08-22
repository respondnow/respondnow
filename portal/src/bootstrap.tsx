import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { AppWithAuthentication, AppWithoutAuthentication } from 'app/App';
// eslint-disable-next-line import/no-cycle
import './styles/global.scss';

ReactDOM.render(
  <BrowserRouter>
    <Switch>
      <Route path={'/account/:accountID'}>
        <AppWithAuthentication />
      </Route>
      <Route path="/">
        <AppWithoutAuthentication />
      </Route>
    </Switch>
  </BrowserRouter>,
  document.getElementById('react-root')
);
