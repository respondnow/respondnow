import React from 'react';
import { Icon } from '@harnessio/icons';
import { Color, FontVariation } from '@harnessio/design-system';
import { Avatar, Button, Layout, Text } from '@harnessio/uicore';
import { TimelineUtilReturn } from '@interfaces';
import { useStrings } from '@strings';
import SlackIcon from '@images/slack.svg';
import { generateSlackChannelLink } from '@utils';
import SeverityBadge from '@components/SeverityBadge';
import StatusBadge from '@components/StatusBadge';
import { Incident, Timeline } from '@services/server';
import css from './IncidentTimeline.module.scss';

interface IncidentTimelineHelperProps {
  incident: Incident | undefined;
  timeline: Timeline;
}

export function getTimelinePropsBasedOnIncidentData(
  props: IncidentTimelineHelperProps
): TimelineUtilReturn | undefined {
  const { timeline, incident } = props;
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const { getString } = useStrings();

  const userName = timeline.userDetails?.name || timeline.userDetails?.userName;

  if (!timeline || !incident) {
    return undefined;
  }

  const SlackIconRenderer = (
    <div className={css.slackIcon}>
      <img src={SlackIcon} height={14} width={14} alt="Slack" />
    </div>
  );

  switch (timeline.type) {
    case 'Incident_Created':
      return {
        icon: SlackIconRenderer,
        headerContent: (
          <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} style={{ gap: '0.2rem' }}>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.PRIMARY_7}>
              {userName}
            </Text>
            <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
              {getString('reportedAn')}
            </Text>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
              {incident.severity.split(' - ')[0]}
            </Text>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
              {getString('incident').toLowerCase()}
            </Text>
          </Layout.Horizontal>
        ),
        bodyContent: (
          <Layout.Vertical style={{ gap: '0.25rem' }}>
            <Layout.Horizontal
              flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
              style={{ gap: '0.2rem' }}
            >
              <Text font={{ variation: FontVariation.SMALL_SEMI }} color={Color.GREY_800}>
                {getString('title')}:
              </Text>
              <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500}>
                &ldquo;{incident.name}&ldquo;
              </Text>
            </Layout.Horizontal>
            <Layout.Horizontal
              flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
              style={{ gap: '0.2rem' }}
            >
              <Text font={{ variation: FontVariation.SMALL_SEMI }} color={Color.GREY_800}>
                {getString('summary')}:
              </Text>
              <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500}>
                &ldquo;{incident.summary}&ldquo;
              </Text>
            </Layout.Horizontal>
          </Layout.Vertical>
        )
      };
    case 'Slack_Channel_Created':
      return {
        icon: SlackIconRenderer,
        headerContent: (
          <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
            {getString('incidentChannelCreated')}
          </Text>
        ),
        bodyContent: (
          <Layout.Vertical style={{ gap: '0.25rem' }}>
            <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} style={{ gap: '0.2rem' }}>
              <Text font={{ variation: FontVariation.SMALL_SEMI }} color={Color.GREY_800}>
                {getString('addedMemberAutomatically')}:
              </Text>
              <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500}>
                {userName}
              </Text>
            </Layout.Horizontal>
            <Button
              noStyling
              width={210}
              className={css.ctaButton}
              disabled={!incident.incidentChannel?.slack?.teamDomain || !incident.channels?.[0].id}
              onClick={() => {
                window.open(
                  generateSlackChannelLink(
                    incident.incidentChannel?.slack?.teamDomain || '',
                    incident.channels?.[0].id || ''
                  ),
                  '_blank'
                );
              }}
            >
              <Layout.Horizontal
                flex={{ alignItems: 'center', justifyContent: 'space-between' }}
                style={{ gap: '0.25rem' }}
              >
                <img src={SlackIcon} alt="Slack" height={12} />
                <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500}>
                  {`${timeline.slack?.channelName?.slice(0, 22)}...`}
                </Text>
                <Icon name="link" size={10} color={Color.GREY_800} />
              </Layout.Horizontal>
            </Button>
          </Layout.Vertical>
        )
      };
    case 'Comment':
      return {
        icon: SlackIconRenderer,
        headerContent: (
          <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} style={{ gap: '0.2rem' }}>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.PRIMARY_7}>
              {userName}
            </Text>
            <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
              {getString('commentedIn')}
            </Text>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
              {getString('slackChannel')}
            </Text>
          </Layout.Horizontal>
        ),
        bodyContent: (
          <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500}>
            &ldquo;{timeline.currentState}&ldquo;
          </Text>
        )
      };
    case 'Severity':
      return {
        icon: SlackIconRenderer,
        headerContent: (
          <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} style={{ gap: '0.2rem' }}>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.PRIMARY_7}>
              {userName}
            </Text>
            <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
              {getString('updated').toLowerCase()}
            </Text>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
              {getString('incident')} {getString('severity')}
            </Text>
          </Layout.Horizontal>
        ),
        bodyContent: (
          <Layout.Horizontal flex={{ align: 'center-center', justifyContent: 'flex-start' }} style={{ gap: '0.25rem' }}>
            <SeverityBadge severity={timeline.previousState as Incident['severity']} />
            <Icon name="arrow-right" size={13} color={Color.GREY_500} />
            <SeverityBadge severity={timeline.currentState as Incident['severity']} />
          </Layout.Horizontal>
        )
      };
    case 'Status':
      return {
        icon: SlackIconRenderer,
        headerContent: (
          <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} style={{ gap: '0.2rem' }}>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.PRIMARY_7}>
              {userName}
            </Text>
            <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
              {getString('updated').toLowerCase()}
            </Text>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
              {getString('incident')} {getString('status')}
            </Text>
          </Layout.Horizontal>
        ),
        bodyContent: (
          <Layout.Horizontal flex={{ align: 'center-center', justifyContent: 'flex-start' }} style={{ gap: '0.25rem' }}>
            <StatusBadge status={timeline.previousState as Incident['status']} />
            <Icon name="arrow-right" size={13} color={Color.GREY_500} />
            <StatusBadge status={timeline.currentState as Incident['status']} />
          </Layout.Horizontal>
        )
      };
    case 'Summary':
      return {
        icon: SlackIconRenderer,
        headerContent: (
          <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} style={{ gap: '0.2rem' }}>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.PRIMARY_7}>
              {userName}
            </Text>
            <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
              {getString('updated').toLowerCase()}
            </Text>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
              {getString('incident')} {getString('summary')}
            </Text>
          </Layout.Horizontal>
        ),
        bodyContent: (
          <Layout.Vertical style={{ gap: '0.25rem' }}>
            <Layout.Horizontal
              flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
              style={{ gap: '0.2rem' }}
            >
              <Text font={{ variation: FontVariation.SMALL_SEMI }} color={Color.GREY_800}>
                {getString('from')}:
              </Text>
              <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500}>
                &ldquo;{timeline.previousState}&ldquo;
              </Text>
            </Layout.Horizontal>
            <Layout.Horizontal
              flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
              style={{ gap: '0.2rem' }}
            >
              <Text font={{ variation: FontVariation.SMALL_SEMI }} color={Color.GREY_800}>
                {getString('to')}:
              </Text>
              <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500}>
                &ldquo;{timeline.currentState}&ldquo;
              </Text>
            </Layout.Horizontal>
          </Layout.Vertical>
        )
      };
    case 'Roles':
      return {
        icon: SlackIconRenderer,
        headerContent: (
          <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} style={{ gap: '0.2rem' }}>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.PRIMARY_7}>
              {userName}
            </Text>
            <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
              {getString('updated').toLowerCase()}
            </Text>
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
              {getString('incident')} {getString('keyMembers')}
            </Text>
          </Layout.Horizontal>
        ),
        bodyContent: (
          <Layout.Horizontal style={{ gap: '0.5rem' }} flex={{ alignItems: 'center', justifyContent: 'flex-start' }}>
            <Layout.Vertical style={{ gap: '0.25rem' }}>
              {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
              {(timeline.additionalDetails as any).previousState.map((member: any) => (
                <Layout.Horizontal
                  key={member?.userDetails?.username}
                  flex={{ alignItems: 'center', justifyContent: 'flex-start' }}
                >
                  <Avatar hoverCard={false} src={SlackIcon} size="small" />
                  <Layout.Vertical>
                    <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
                      {member?.userDetails?.name || member?.userDetails?.username}
                    </Text>
                    <Text font={{ variation: FontVariation.TINY, italic: true }} color={Color.GREY_500}>
                      {member?.roleType}
                    </Text>
                  </Layout.Vertical>
                </Layout.Horizontal>
              ))}
            </Layout.Vertical>
            <Icon name="arrow-right" size={13} color={Color.GREY_500} />
            <Layout.Vertical style={{ gap: '0.25rem' }}>
              {/* eslint-disable-next-line @typescript-eslint/no-explicit-any */}
              {(timeline.additionalDetails as any).currentState.map((member: any) => (
                <Layout.Horizontal
                  key={member?.userDetails?.username}
                  flex={{ alignItems: 'center', justifyContent: 'flex-start' }}
                >
                  <Avatar hoverCard={false} src={SlackIcon} size="small" />
                  <Layout.Vertical>
                    <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
                      {member?.userDetails?.name || member?.userDetails?.username}
                    </Text>
                    <Text font={{ variation: FontVariation.TINY, italic: true }} color={Color.GREY_500}>
                      {member?.roleType}
                    </Text>
                  </Layout.Vertical>
                </Layout.Horizontal>
              ))}
            </Layout.Vertical>
          </Layout.Horizontal>
        )
      };
    default:
      return undefined;
  }
}
