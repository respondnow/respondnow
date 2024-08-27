import React from 'react';
import { Button, ButtonSize, ButtonVariation, Card, Container, Layout, Text } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import { withErrorBoundary } from 'react-error-boundary';
import { DefaultLayout } from '@layouts';
import getStartedSlack from '@images/getStartedSlack.svg';
import { useStrings } from '@strings';
import { Fallback } from '@errors';
import { useAppStore } from '@hooks';
import css from './GettingStarted.module.scss';

const GettingStartedView: React.FC = () => {
  const { getString } = useStrings();
  const { currentUserInfo } = useAppStore();

  return (
    <DefaultLayout scale="full-screen">
      <Layout.Vertical height="100%" flex={{ align: 'center-center' }} style={{ gap: '1rem', overflow: 'auto' }}>
        <Layout.Vertical width="100%" style={{ gap: '0.5rem' }} flex={{ align: 'center-center' }}>
          <Text font={{ variation: FontVariation.H6 }} color={Color.GREY_800}>
            {getString('welcomeText', { name: currentUserInfo.name })},
          </Text>
          <Text font={{ variation: FontVariation.H4 }} color={Color.GREY_800}>
            {getString('getStartedHeading')}
          </Text>
        </Layout.Vertical>
        <Card className={css.cardContainer}>
          <Layout.Vertical width="100%" height="100%">
            <Container flex={{ align: 'center-center' }} height={250} background={Color.PURPLE_100}>
              <img src={getStartedSlack} alt="Slack" height={200} />
            </Container>
            <Layout.Vertical
              width="100%"
              style={{ gap: '1rem', flexGrow: 1 }}
              padding="large"
              flex={{ align: 'center-center' }}
            >
              <Text font={{ variation: FontVariation.H4 }} color={Color.GREY_800}>
                {getString('slackConfigHeader')}
              </Text>
              <Text font={{ variation: FontVariation.BODY, align: 'center' }} color={Color.GREY_800}>
                {getString('slackConfigDescription')}
              </Text>
              <Button
                size={ButtonSize.MEDIUM}
                text={getString('slackButtonText')}
                width="100%"
                variation={ButtonVariation.PRIMARY}
              />
            </Layout.Vertical>
          </Layout.Vertical>
        </Card>
      </Layout.Vertical>
    </DefaultLayout>
  );
};

export default withErrorBoundary(GettingStartedView, { FallbackComponent: Fallback });
