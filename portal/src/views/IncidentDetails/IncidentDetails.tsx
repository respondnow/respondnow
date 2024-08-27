import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Button, ButtonVariation, Card, Layout } from '@harnessio/uicore';
import { Color } from '@harnessio/design-system';
import { Fallback } from '@errors';
import { DefaultLayout } from '@layouts';
import { IncidentIncident } from '@services/server';
import { generateSlackChannelLink } from '@utils';
import SlackIcon from '@images/slack.svg';
import DetailsSection from './sections/DetailsSection';
import TimelineSection from './sections/Timeline';
import css from './IncidentDetails.module.scss';

interface IncidentDetailsViewProps {
  incidentData: IncidentIncident | undefined;
  incidentDataLoading: boolean;
}

const IncidentDetailsView: React.FC<IncidentDetailsViewProps> = props => {
  const { incidentData, incidentDataLoading } = props;

  const isIncidentPresent = !!incidentData;

  return (
    <DefaultLayout
      title={incidentData?.name || 'Incident Details'}
      loading={incidentDataLoading}
      noData={!isIncidentPresent}
      noDataProps={{
        title: 'No Incident Found',
        subtitle: 'The incident you are looking for does not exist or has been deleted.'
      }}
      toolbar={
        <Button
          onClick={() => {
            window.open(
              generateSlackChannelLink(
                incidentData?.incidentChannel?.slack?.teamDomain || '',
                incidentData?.channels?.[0].id || ''
              ),
              '_blank'
            );
          }}
          variation={ButtonVariation.SECONDARY}
          text="View Channel"
          icon={<img src={SlackIcon} height={16} />}
          disabled={!incidentData?.incidentChannel?.slack?.teamDomain || !incidentData?.channels?.[0].id}
        />
      }
    >
      <Layout.Horizontal height="100%" spacing="large" background={Color.PRIMARY_BG}>
        <Card className={css.detailsCardContainer}>
          <DetailsSection incidentData={incidentData} />
        </Card>
        <Card className={css.timelineCardContainer}>
          <TimelineSection incidentData={incidentData} />
        </Card>
      </Layout.Horizontal>
    </DefaultLayout>
  );
};

export default withErrorBoundary(IncidentDetailsView, { FallbackComponent: Fallback });
