import React from 'react';
import { Icon } from '@harnessio/icons';
import { Color, FontVariation } from '@harnessio/design-system';
import { Button, Layout, Text } from '@harnessio/uicore';
import { TimelineUtilReturn } from '@interfaces';
import { IncidentIncident, IncidentSeverity, IncidentStatus, IncidentTimeline } from '@services/server';
import { useStrings } from '@strings';
import SlackIcon from '@images/slack.svg';
import { generateSlackChannelLink } from '@utils';
import SeverityBadge from '@components/SeverityBadge';
import StatusBadge from '@components/StatusBadge';
import css from './IncidentTimeline.module.scss';

interface IncidentTimelineHelperProps {
  incident: IncidentIncident;
  timeline: IncidentTimeline;
}

export function getTimelinePropsBasedOnIncidentData(
  props: IncidentTimelineHelperProps
): TimelineUtilReturn | undefined {
  const { timeline, incident } = props;
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const { getString } = useStrings();

  const userName = timeline.userDetails?.name || timeline.userDetails?.userName;

  if (!timeline) {
    return undefined;
  }

  switch (timeline.type) {
    case 'incidentCreated':
      return {
        icon: (
          <Icon
            name="chaos-litmuschaos"
            size={12}
            background={Color.GREY_100}
            style={{ borderRadius: '10rem', padding: '7px' }}
          />
        ),
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
    case 'slackChannelCreated':
      return {
        icon: (
          <Icon
            name="chaos-litmuschaos"
            size={12}
            background={Color.GREY_100}
            style={{ borderRadius: '10rem', padding: '7px' }}
          />
        ),
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
    case 'comment':
      return {
        icon: (
          <Icon
            name="chaos-litmuschaos"
            size={12}
            background={Color.GREY_100}
            style={{ borderRadius: '10rem', padding: '7px' }}
          />
        ),
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
            &ldquo;{timeline.message}&ldquo;
          </Text>
        )
      };
    case 'severity':
      return {
        icon: (
          <Icon
            name="chaos-litmuschaos"
            size={12}
            background={Color.GREY_100}
            style={{ borderRadius: '10rem', padding: '7px' }}
          />
        ),
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
            <SeverityBadge severity={timeline.previousState as IncidentSeverity} />
            <Icon name="arrow-right" size={13} color={Color.GREY_500} />
            <SeverityBadge severity={timeline.currentState as IncidentSeverity} />
          </Layout.Horizontal>
        )
      };
    case 'status':
      return {
        icon: (
          <Icon
            name="chaos-litmuschaos"
            size={12}
            background={Color.GREY_100}
            style={{ borderRadius: '10rem', padding: '7px' }}
          />
        ),
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
            <StatusBadge status={timeline.previousState as IncidentStatus} />
            <Icon name="arrow-right" size={13} color={Color.GREY_500} />
            <StatusBadge status={timeline.currentState as IncidentStatus} />
          </Layout.Horizontal>
        )
      };
    default:
      return undefined;
  }
}
