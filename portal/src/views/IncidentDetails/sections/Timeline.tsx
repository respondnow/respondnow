import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Container, Layout, Text } from '@harnessio/uicore';
import { FontVariation } from '@harnessio/design-system';
import { Fallback } from '@errors';
import { useStrings } from '@strings';
import { IncidentIncident } from '@services/server';
import css from '../IncidentDetails.module.scss';

interface TimelineSectionProps {
  incidentData: IncidentIncident | undefined;
}

const TimelineSection: React.FC<TimelineSectionProps> = () => {
  const { getString } = useStrings();

  return (
    <Layout.Vertical
      flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
      width="100%"
      height="100%"
      padding="medium"
      className={css.timelineSectionContainer}
    >
      <Text font={{ variation: FontVariation.H5 }}>{getString('incidentTimeline')}</Text>
      <Container>TimeLine Goes Here</Container>
    </Layout.Vertical>
  );
};

export default withErrorBoundary(TimelineSection, { FallbackComponent: Fallback });
