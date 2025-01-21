import { PaginationProps } from '@harnessio/uicore';
import React from 'react';
import { Incident } from '@services/server';

export interface IncidentsTableProps {
  content: Incident[];
  pagination?: PaginationProps;
  isLoading?: boolean;
}

export interface TimelineUtilReturn {
  icon: React.ReactNode;
  headerContent: React.ReactNode;
  bodyContent: React.ReactNode;
}
