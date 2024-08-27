import React from 'react';
import { Redirect, Route, RouteProps } from 'react-router-dom';
import { getTokenFromLocalStorage, isUserAuthenticated, updateUserAndScopeFromAPI } from '@utils';
import { initialAppContext } from '@context';
import { useAppStore } from '@hooks';
import { paths } from './RouteDefinitions';

function AuthenticatedRoute(props: RouteProps): React.ReactElement {
  const showAuthorizedRoutes = isUserAuthenticated();
  const appStore = useAppStore();
  const { token, isInitialLogin } = getTokenFromLocalStorage();

  React.useEffect(() => {
    if (appStore.currentUserInfo === initialAppContext.currentUserInfo) {
      updateUserAndScopeFromAPI(appStore.updateAppStore, token, isInitialLogin);
    }
  }, [appStore, token, isInitialLogin]);

  if (showAuthorizedRoutes) {
    return <Route {...props} />;
  } else {
    return <Redirect to={paths.toLogin()} />;
  }
}

export default AuthenticatedRoute;
