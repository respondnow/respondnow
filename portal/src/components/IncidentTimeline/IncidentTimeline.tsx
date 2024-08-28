import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Container, Layout, Text } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import moment from 'moment';
import { Fallback } from '@errors';
import { IncidentTimeline } from '@services/server';
import { getTimelinePropsBasedOnIncidentData } from './helper';
import css from './IncidentTimeline.module.scss';

interface IncidentTimelineProps {
  timeline: IncidentTimeline[] | undefined;
}

const IncidentTimeline: React.FC<IncidentTimelineProps> = props => {
  const { timeline } = props;

  if (!timeline || !timeline.length) {
    return <></>;
  }

  return (
    <Layout.Vertical width="100%">
      {timeline.map((item, index) => {
        const timelineProps = getTimelinePropsBasedOnIncidentData(item);
        return (
          <Container key={item.id} width="100%" height="100%" className={css.timelineRowContainer}>
            <Container height="100%" padding={{ top: 'xsmall' }}>
              <Text font={{ variation: FontVariation.SMALL, align: 'right' }} color={Color.GREY_800}>
                {moment(item.createdAt).format('MMM D, YYYY')}
              </Text>
              <Text font={{ variation: FontVariation.SMALL, align: 'right' }} color={Color.GREY_500}>
                {moment(item.createdAt).format('h:mm A')}
              </Text>
            </Container>
            <Container height="100%" className={css.dividerContainer}>
              {timelineProps?.icon}
              {!(index === timeline.length - 1) && <div className={css.splitBar} />}
            </Container>
            <Container height="100%" padding={{ top: 'xsmall' }}>
              <Text font={{ variation: FontVariation.SMALL_SEMI }} color={Color.GREY_800}>
                {timelineProps?.headerContent}
              </Text>
              <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500}>
                {timelineProps?.bodyContent}
              </Text>
            </Container>
          </Container>
        );
      })}
    </Layout.Vertical>
  );
};

export default withErrorBoundary(IncidentTimeline, { FallbackComponent: Fallback });
