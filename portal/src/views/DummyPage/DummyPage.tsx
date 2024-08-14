import React from 'react';
import { Button, Container, Layout, Text } from '@harnessio/uicore';
import { Color } from '@harnessio/design-system';
import { useHistory } from 'react-router-dom';
import { useRouteWithBaseUrl } from '@hooks';
import { useStrings } from '@strings';

const DummyPage: React.FC = () => {
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
          {getString('respondNow')} {getString('home')}
        </Text>
        <Button icon="play" text={getString('login')} onClick={() => history.push(paths.toLogin())} />
      </Container>
    </Layout.Vertical>
  );
};

export default DummyPage;
