import React from 'react';
import { Route, Switch } from 'react-router-dom';
import { useRouteDefinitionsMatch } from '@hooks';
import DummyPage from '@views/DummyPage';
import DummyLogin from '@views/DummyLogin';

export function Routes(): React.ReactElement {
  const matchPath = useRouteDefinitionsMatch();

  return (
    <Switch>
      <Route exact path={matchPath.toRoot()} component={DummyPage} />
      <Route exact path={matchPath.toLogin()} component={DummyLogin} />
    </Switch>
  );
}
