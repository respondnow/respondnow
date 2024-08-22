import React from 'react';
import { Redirect, Route, RouteProps } from 'react-router-dom';
import { getUserDetails, isUserAuthenticated } from '@utils';
import { paths } from './RouteDefinitions';

function AuthenticatedRoute(props: RouteProps): React.ReactElement {
  const { accessToken: token } = getUserDetails();
  const showAuthorizedRoutes = token && isUserAuthenticated();

  if (showAuthorizedRoutes) {
    return <Route {...props} />;
  } else {
    return <Redirect to={paths.toLogin()} />;
  }
}

export default AuthenticatedRoute;
