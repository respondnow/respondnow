import React from 'react';
import {
  Button,
  ButtonSize,
  ButtonVariation,
  Card,
  Container,
  Layout,
  Stepper,
  Text,
  useToggleOpen
} from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import { withErrorBoundary } from 'react-error-boundary';
import { Dialog } from '@blueprintjs/core';
import { useHistory } from 'react-router-dom';
import { DefaultLayout } from '@layouts';
import getStartedSlack from '@images/getStartedSlack.svg';
import { useStrings } from '@strings';
import { Fallback } from '@errors';
import { useAppStore } from '@hooks';
import { paths } from '@routes/RouteDefinitions';
import css from './GettingStarted.module.scss';

const GettingStartedView: React.FC = () => {
  const { getString } = useStrings();
  const { currentUserInfo } = useAppStore();
  const { open, close, isOpen } = useToggleOpen();
  const history = useHistory();

  return (
    <DefaultLayout
      scale="full-screen"
      popovers={
        <Dialog isOpen={isOpen} onClose={close} style={{ paddingBottom: 0, height: 430, width: 600 }}>
          <Container height="100%" width="100%" padding="large">
            <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'space-between' }}>
              <Text font={{ variation: FontVariation.H4 }} color={Color.GREY_800}>
                {getString('gettingStarted')}
              </Text>
              <Button variation={ButtonVariation.ICON} icon="cross" size={ButtonSize.SMALL} onClick={close} />
            </Layout.Horizontal>
            <Container className={css.stepperContainer} margin={{ top: 'xxlarge' }}>
              <Stepper
                id="gettingStarted"
                stepList={[
                  {
                    id: 'createSlackApp',
                    title: getString('gettingStartedCreateSlackAppHeader'),
                    panel: (
                      <Text font={{ variation: FontVariation.BODY }}>
                        {getString('gettingStartedCreateSlackAppContent')}
                      </Text>
                    )
                  },
                  {
                    id: 'installSlackApp',
                    title: getString('gettingStartedInstallSlackAppHeader'),
                    panel: (
                      <Text font={{ variation: FontVariation.BODY }}>
                        {getString('gettingStartedInstallSlackAppContent')}
                      </Text>
                    )
                  },
                  {
                    id: 'configureIncidentChannel',
                    title: getString('gettingStartedConfigureIncidentChannelHeader'),
                    panel: (
                      <Layout.Vertical
                        style={{ gap: '1rem' }}
                        flex={{ alignItems: 'center', justifyContent: 'flex-start' }}
                      >
                        <Text font={{ variation: FontVariation.BODY }}>
                          {getString('gettingStartedConfigureIncidentChannelContent')}
                        </Text>
                        <Button
                          text={
                            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.WHITE}>
                              {getString('goToIncidentDashboard')}
                            </Text>
                          }
                          variation={ButtonVariation.PRIMARY}
                          size={ButtonSize.SMALL}
                          onClick={() => history.push(paths.toIncidentDashboard())}
                        />
                      </Layout.Vertical>
                    )
                  }
                ]}
                isStepValid={() => true}
              />
            </Container>
          </Container>
        </Dialog>
      }
    >
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
                onClick={open}
              />
            </Layout.Vertical>
          </Layout.Vertical>
        </Card>
      </Layout.Vertical>
    </DefaultLayout>
  );
};

export default withErrorBoundary(GettingStartedView, { FallbackComponent: Fallback });
