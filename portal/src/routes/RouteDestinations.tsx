import React from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import LoginController from '@controllers/Login';
import { GenericErrorHandler } from '@errors';
import GettingStartedController from '@controllers/GettingStarted';
import IncidentsController from '@controllers/Incidents';
import IncidentDetailsController from '@controllers/IncidentDetails';
import { paths } from './RouteDefinitions';
import AuthenticatedRoute from './AuthenticatedRoute';
import UnauthenticatedRoute from './UnauthenticatedRoute';

export function Routes(): React.ReactElement {
  const incidentId = ':incidentId';

  return (
    <Switch>
      <Redirect exact from={paths.toRoot()} to={paths.toLogin()} />
      <AuthenticatedRoute exact path={paths.toGetStarted()} component={GettingStartedController} />
      <AuthenticatedRoute exact path={paths.toIncidentDashboard()} component={IncidentsController} />
      <AuthenticatedRoute exact path={paths.toIncidentDetails({ incidentId })} component={IncidentDetailsController} />
      {/* TEMP */}
      <AuthenticatedRoute exact path={paths.toIncidentDetailsDummy()} component={IncidentDetailsController} />
      {/* REMOVE TEMP */}
      <UnauthenticatedRoute exact path={paths.toLogin()} component={LoginController} />
      <Route path="*" component={GenericErrorHandler} />
    </Switch>
  );
}
