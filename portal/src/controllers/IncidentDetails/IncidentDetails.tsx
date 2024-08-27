import React from 'react';
import { useParams } from 'react-router-dom';
import IncidentDetailsView from '@views/IncidentDetails';
import { useGetIncidentQuery } from '@services/server';
import { IncidentDetailsPathProps } from '@routes/RouteDefinitions';
import { getScope, scopeExists } from '@utils';

const IncidentDetailsController: React.FC = () => {
  const scope = getScope();
  const { incidentId } = useParams<IncidentDetailsPathProps>();

  const { data: incidentData, isLoading: incidentDataLoading } = useGetIncidentQuery(
    {
      incidentIdentifier: incidentId,
      queryParams: scope
    },
    {
      enabled: !!incidentId && scopeExists()
    }
  );

  return <IncidentDetailsView incidentData={incidentData?.data} incidentDataLoading={incidentDataLoading} />;
};

export default IncidentDetailsController;
