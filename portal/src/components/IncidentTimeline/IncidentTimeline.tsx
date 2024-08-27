import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Container, Layout } from '@harnessio/uicore';
import { Icon } from '@harnessio/icons';
import { Color } from '@harnessio/design-system';
import { Fallback } from '@errors';
import { IncidentTimeline } from '@services/server';
import css from './IncidentTimeline.module.scss';

interface IncidentTimelineProps {
  timeline: IncidentTimeline[];
}

const IncidentTimeline: React.FC<IncidentTimelineProps> = props => {
  const { timeline } = props;

  return (
    <Layout.Vertical width="100%" height={150}>
      {timeline.map(item => (
        <Container key={item.id} width="100%" height="100%" className={css.timelineRowContainer}>
          <Container height="100%" padding={{ top: 'xsmall' }}>
            TimeFrame
          </Container>
          <Container height="100%" className={css.dividerContainer}>
            <Icon
              name="Account"
              size={15}
              background={Color.GREY_100}
              style={{ borderRadius: '10rem', padding: '6px' }}
            />
            <div />
          </Container>
          <Container height="100%" padding={{ top: 'xsmall' }}>
            Stage Details
          </Container>
        </Container>
      ))}
    </Layout.Vertical>
  );
};

export default withErrorBoundary(IncidentTimeline, { FallbackComponent: Fallback });
