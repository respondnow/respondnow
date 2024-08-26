import { PaginationProps } from '@harnessio/uicore';
import { IncidentIncident } from '@services/server';

export interface IncidentsTableProps {
  content: IncidentIncident[];
  pagination?: PaginationProps;
  isLoading?: boolean;
}
