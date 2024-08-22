import React from 'react';
import { Redirect, Route, RouteProps } from 'react-router-dom';
import { getUserDetails, isUserAuthenticated } from '@utils';
import { paths } from './RouteDefinitions';

function UnauthenticatedRoute(props: RouteProps): React.ReactElement {
  const { accessToken: token } = getUserDetails();
  const isUserLoggedIn = token && isUserAuthenticated();

  if (!isUserLoggedIn) {
    return <Route {...props} />;
  } else {
    return <Redirect to={paths.toIncidentDashboard()} />;
  }
}

export default UnauthenticatedRoute;
