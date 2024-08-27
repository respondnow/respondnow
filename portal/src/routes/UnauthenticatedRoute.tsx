import React from 'react';
import { Redirect, Route, RouteProps } from 'react-router-dom';
import { isUserAuthenticated } from '@utils';
import { paths } from './RouteDefinitions';

function UnauthenticatedRoute(props: RouteProps): React.ReactElement {
  const isUserLoggedIn = isUserAuthenticated();

  if (!isUserLoggedIn) {
    return <Route {...props} />;
  } else {
    return <Redirect to={paths.toIncidentDashboard()} />;
  }
}

export default UnauthenticatedRoute;
