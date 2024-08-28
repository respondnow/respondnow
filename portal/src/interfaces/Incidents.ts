import { PaginationProps } from '@harnessio/uicore';
import React from 'react';
import { IncidentIncident } from '@services/server';

export interface IncidentsTableProps {
  content: IncidentIncident[];
  pagination?: PaginationProps;
  isLoading?: boolean;
}

export interface TimelineUtilReturn {
  icon: React.ReactNode;
  headerContent: React.ReactNode;
  bodyContent: React.ReactNode;
}
