import React from 'react';
import { Button, ButtonVariation } from '@harnessio/uicore';
import { DefaultLayout } from '@layouts';
import { useStrings } from '@strings';
import { IncidentListResponseDto } from '@services/server/schemas/IncidentListResponseDto';

interface IncidentsViewProps {
  incidents: IncidentListResponseDto | undefined;
  loading: boolean;
}

const IncidentsView: React.FC<IncidentsViewProps> = props => {
  const { incidents, loading } = props;
  const { getString } = useStrings();

  return (
    <DefaultLayout
      loading={loading}
      title={`${getString('incidents')}(${incidents?.data?.pagination?.totalItems})`}
      toolbar={<Button variation={ButtonVariation.PRIMARY} text="Report Incident on Slack" />}
      subHeader={<div>Incidents Subheader</div>}
      footer={<div>Incidents Footer</div>}
    >
      <div>Incidents List</div>
    </DefaultLayout>
  );
};

export default IncidentsView;
