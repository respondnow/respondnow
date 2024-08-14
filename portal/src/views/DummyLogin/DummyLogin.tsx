import React from 'react';
import { Button, Container, Layout, Text } from '@harnessio/uicore';
import { Color } from '@harnessio/design-system';
import { useHistory } from 'react-router-dom';
import { useRouteWithBaseUrl } from '@hooks';
import { useStrings } from '@strings';

const DummyLogin: React.FC = () => {
  const history = useHistory();
  const paths = useRouteWithBaseUrl();
  const { getString } = useStrings();
  return (
    <Layout.Vertical height="100%" width="100%" flex={{ align: 'center-center' }}>
      <Container
        intent="primary"
        padding="medium"
        font={{
          align: 'center'
        }}
        background={Color.PURPLE_100}
        border={{
          color: Color.PURPLE_500
        }}
      >
        <Text margin={{ bottom: 'large' }}>
          {getString('respondNow')} {getString('login')}
        </Text>
        <Button icon="arrow-left" text={getString('home')} onClick={() => history.push(paths.toRoot())} />
      </Container>
    </Layout.Vertical>
  );
};

export default DummyLogin;
