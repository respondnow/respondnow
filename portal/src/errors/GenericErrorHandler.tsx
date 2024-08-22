import React from 'react';
import { Link, useHistory } from 'react-router-dom';
import { Layout, Heading, Text, Container } from '@harnessio/uicore';
import { useStrings } from '@strings';
import respondLogo from '@images/respondNow.svg';
import { paths } from '@routes/RouteDefinitions';

interface ErrorType {
  errorMessage?: string;
  errStatusCode?: number | string;
  allowUserToGoBack?: boolean;
}

export function GenericErrorHandler({ errorMessage, errStatusCode, allowUserToGoBack }: ErrorType): JSX.Element {
  const history = useHistory();
  const { getString } = useStrings();
  return (
    <Container height="var(--page-min-height)" flex={{ align: 'center-center' }}>
      <Layout.Vertical spacing="large" flex={{ align: 'center-center' }}>
        <Heading>{errStatusCode || 404}</Heading>
        <Text>{errorMessage || getString('404Error')}</Text>
        {allowUserToGoBack ? (
          <Text onClick={() => history.goBack()}>
            <a>{getString('goBack')}</a>
          </Text>
        ) : (
          <Link to={paths.toRoot()}>{getString('goToHome')}</Link>
        )}
        <img height={40} src={respondLogo} alt="RespondNow" />
      </Layout.Vertical>
    </Container>
  );
}
