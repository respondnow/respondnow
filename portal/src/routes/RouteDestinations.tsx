import React from 'react';
import { Redirect, Route, Switch, useHistory } from 'react-router-dom';
import { useRouteDefinitionsMatch } from '@hooks';
import LoginController from '@controllers/Login';
import { GenericErrorHandler } from '@errors';
import { getUserDetails, isUserAuthenticated, useLogout } from '@utils';

export function RoutesWithAuthentication(): React.ReactElement {
  const paths = useRouteDefinitionsMatch();
  const history = useHistory();
  const { accessToken: token, isInitialLogin } = getUserDetails();

  const { forceLogout } = useLogout();
  React.useEffect(() => {
    if (!token || !isUserAuthenticated()) {
      forceLogout();
    }
    if (isInitialLogin) {
      history.push(`/settings/password-reset`);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token, isInitialLogin]);

  return (
    <Switch>
      <Redirect exact from={paths.toRoot()} to={paths.toIncidentDashboard()} />
      {/* Auth Routes */}
      {/* <Route exact path={paths.toPasswordReset()} component={Dummy} /> */}
    </Switch>
  );
}

export function RoutesWithoutAuthentication(): React.ReactElement {
  const paths = useRouteDefinitionsMatch();

  return (
    <Switch>
      <Redirect exact from={paths.toRoot()} to={paths.toLogin()} />
      {/* Auth Routes */}
      <Route exact path={paths.toLogin()} component={LoginController} />
      <Route path="*" component={GenericErrorHandler} />
    </Switch>
  );
}
