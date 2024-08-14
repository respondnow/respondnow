import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
// eslint-disable-next-line import/no-cycle
import { App } from 'app/App';
import './styles/global.scss';

ReactDOM.render(
  <BrowserRouter>
    <Switch>
      <Route path="/">
        <App />
      </Route>
    </Switch>
  </BrowserRouter>,
  document.getElementById('react-root')
);
