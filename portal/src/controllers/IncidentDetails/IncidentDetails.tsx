import React from 'react';
import { useParams } from 'react-router-dom';
import IncidentDetailsView from '@views/IncidentDetails';
import { useGetIncidentQuery } from '@services/server';
import { IncidentDetailsPathProps } from '@routes/RouteDefinitions';

const IncidentDetailsController: React.FC = () => {
  const { incidentId } = useParams<IncidentDetailsPathProps>();

  const { data: incidentData, isLoading: incidentDataLoading } = useGetIncidentQuery(
    {
      incidentIdentifier: incidentId,
      queryParams: {
        accountIdentifier: 'default',
        projectIdentifier: 'default',
        orgIdentifier: 'default'
      }
    },
    {
      enabled: !!incidentId
    }
  );

  return <IncidentDetailsView incidentData={incidentData?.data} incidentDataLoading={incidentDataLoading} />;
};

export default IncidentDetailsController;
