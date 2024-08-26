import React from 'react';
// import { useParams } from 'react-router-dom';
import IncidentDetailsView from '@views/IncidentDetails';
import { useGetIncidentQuery } from '@services/server';
// import { IncidentDetailsPathProps } from '@routes/RouteDefinitions';

const IncidentDetailsController: React.FC = () => {
  // const { incidentId } = useParams<IncidentDetailsPathProps>();

  const { data: incidentData, isLoading: incidentDataLoading } = useGetIncidentQuery({
    incidentIdentifier: 'Access Control Is Up-d7364861-4917-4393-ae34-df45c4b44a14',
    queryParams: {
      accountIdentifier: 'default',
      projectIdentifier: 'default',
      orgIdentifier: 'default'
    }
  });

  return <IncidentDetailsView incidentData={incidentData?.data} incidentDataLoading={incidentDataLoading} />;
};

export default IncidentDetailsController;
