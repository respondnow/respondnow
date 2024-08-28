import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Container, Layout, Text } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import moment from 'moment';
import { Fallback } from '@errors';
import { IncidentIncident, IncidentTimeline } from '@services/server';
import { getTimelinePropsBasedOnIncidentData } from './helper';
import css from './IncidentTimeline.module.scss';

interface IncidentTimelineProps {
  incident: IncidentIncident | undefined;
}

const IncidentTimeline: React.FC<IncidentTimelineProps> = props => {
  const { incident } = props;

  if (!incident || !incident.timelines) {
    return null;
  }

  const timelines = incident.timelines;

  return (
    <Layout.Vertical width="100%">
      {timelines.map((item, index) => {
        const timelineProps = getTimelinePropsBasedOnIncidentData({ incident, timeline: item });
        return (
          <Container key={item.id} width="100%" height="100%" className={css.timelineRowContainer}>
            <Container height="100%" padding={{ top: 'xsmall', bottom: 'large' }}>
              {item.createdAt && (
                <Text font={{ variation: FontVariation.SMALL, align: 'right' }} color={Color.GREY_800}>
                  {moment(item.createdAt * 1000).format('MMM D, YYYY')}
                </Text>
              )}
              {item.createdAt && (
                <Text font={{ variation: FontVariation.SMALL, align: 'right' }} color={Color.GREY_500}>
                  {moment(item.createdAt * 1000).format('h:mm A')}
                </Text>
              )}
            </Container>
            <Container height="100%" className={css.dividerContainer}>
              {timelineProps?.icon}
              {!(index === timelines.length - 1) && <div className={css.splitBar} />}
            </Container>
            <Layout.Vertical height="100%" padding={{ top: 'xsmall', bottom: 'large' }} style={{ gap: '0.5rem' }}>
              {timelineProps?.headerContent}
              {timelineProps?.bodyContent}
            </Layout.Vertical>
          </Container>
        );
      })}
    </Layout.Vertical>
  );
};

export default withErrorBoundary(IncidentTimeline, { FallbackComponent: Fallback });
