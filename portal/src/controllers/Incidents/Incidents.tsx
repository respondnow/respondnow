import React from 'react';
import { useToaster } from '@harnessio/uicore';
import IncidentsView from '@views/Incidents';
import { getScope } from '@utils';
import { useListIncidentsQuery } from '@services/server/hooks/useListIncidentsQuery';

const IncidentsController: React.FC = () => {
  const scope = getScope();
  const { showError } = useToaster();

  const { data: incidentList, isLoading: incidentListLoading } = useListIncidentsQuery(
    {
      queryParams: scope
    },
    {
      onError: error => {
        showError(error.message);
      }
    }
  );

  return <IncidentsView incidents={incidentList} loading={incidentListLoading} />;
};

export default IncidentsController;
