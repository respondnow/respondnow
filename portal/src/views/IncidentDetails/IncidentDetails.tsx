import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Button, ButtonVariation, Layout } from '@harnessio/uicore';
import { Color } from '@harnessio/design-system';
import { Fallback } from '@errors';
import { DefaultLayout } from '@layouts';
import { Incident } from '@services/server';
import { generateSlackChannelLink } from '@utils';
import SlackIcon from '@images/slack-mono.svg';
import DetailsSection from './sections/DetailsSection';
import TimelineSection from './sections/Timeline';

interface IncidentDetailsViewProps {
  incidentData: Incident | undefined;
  incidentDataLoading: boolean;
}

const IncidentDetailsView: React.FC<IncidentDetailsViewProps> = props => {
  const { incidentData, incidentDataLoading } = props;

  const isIncidentPresent = !!incidentData;

  return (
    <DefaultLayout
      title={incidentData?.name || 'Incident Name'}
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
        <DetailsSection incidentData={incidentData} />
        <TimelineSection incidentData={incidentData} />
      </Layout.Horizontal>
    </DefaultLayout>
  );
};

export default withErrorBoundary(IncidentDetailsView, { FallbackComponent: Fallback });
