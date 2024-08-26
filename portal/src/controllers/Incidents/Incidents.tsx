import React from 'react';
import { useToaster } from '@harnessio/uicore';
import { isEqual } from 'lodash-es';
import IncidentsView from '@views/Incidents';
import { getScope } from '@utils';
import { useListIncidentsQuery } from '@services/server/hooks/useListIncidentsQuery';
import { initialIncidenrsFilterState, useIncidentsFilter, usePagination } from '@hooks';
import { IncidentsTableProps } from '@interfaces';
import {
  FilterProps,
  IncidentsSearchBar,
  IncidentsSeverityFilter,
  IncidentsStatusFilter,
  ResetFilterButton
} from './IncidentsFilters';

const IncidentsController: React.FC = () => {
  const scope = getScope();
  const { showError } = useToaster();
  // Filter props
  const { page, limit, setPage, setLimit, pageSizeOptions } = usePagination([10, 20], { page: 0, limit: 10 }, true);
  const { state, dispatch } = useIncidentsFilter();
  const resetPage = (): void => {
    setPage(0);
  };

  const filterProps: FilterProps = {
    state,
    dispatch,
    resetPage
  };

  const { data: incidentList, isLoading: incidentListLoading } = useListIncidentsQuery(
    {
      queryParams: {
        accountIdentifier: 'default',
        projectIdentifier: 'default',
        orgIdentifier: 'default',
        page,
        all: false,
        limit,
        search: state.incidentName,
        status: state.incidentStatus,
        severity: state.incidentSeverity
      }
    },
    {
      onError: error => {
        showError(error.message);
      },
      refetchInterval: 20000
    }
  );

  const tableData: IncidentsTableProps = {
    content: incidentList?.data?.content ?? [],
    pagination: {
      itemCount: incidentList?.data?.pagination?.totalItems || 0,
      pageCount: incidentList?.data?.pagination?.totalPages || 0,
      pageIndex: incidentList?.data?.pagination?.index || 0,
      pageSize: incidentList?.data?.pagination?.limit || 10,
      pageSizeOptions: pageSizeOptions,
      gotoPage: setPage,
      onPageSizeChange: setLimit
    },
    isLoading: incidentListLoading
  };

  const areFiltersSet = !isEqual(state, initialIncidenrsFilterState);

  return (
    <IncidentsView
      tableData={tableData}
      incidentsSearchBar={<IncidentsSearchBar {...filterProps} />}
      incidentsStatusFilter={<IncidentsStatusFilter {...filterProps} />}
      incidentsSeverityFilter={<IncidentsSeverityFilter {...filterProps} />}
      resetFiltersButton={<ResetFilterButton {...filterProps} />}
      areFiltersSet={areFiltersSet}
    />
  );
};

export default IncidentsController;
