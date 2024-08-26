import React from 'react';
import { Button, ButtonVariation, FlexExpander, Layout, Pagination } from '@harnessio/uicore';
import { DefaultLayout } from '@layouts';
import { useStrings } from '@strings';
import { IncidentsTableProps } from '@interfaces';
import MemoisedIncidentListTable from '@tables/Incidents/ListIncidentsTable';

interface IncidentsViewProps {
  tableData: IncidentsTableProps;
  incidentsSearchBar: React.JSX.Element;
  incidentsStatusFilter: React.JSX.Element;
  incidentsSeverityFilter: React.JSX.Element;
  resetFiltersButton: React.JSX.Element;
  areFiltersSet: boolean;
}

const IncidentsView: React.FC<IncidentsViewProps> = props => {
  const {
    tableData,
    incidentsSearchBar,
    incidentsStatusFilter,
    incidentsSeverityFilter,
    areFiltersSet,
    resetFiltersButton
  } = props;
  const { getString } = useStrings();

  const isDataPresent = !!tableData.content.length;

  const reportIncidentButton = <Button variation={ButtonVariation.PRIMARY} text="Report Incident on Slack" />;

  return (
    <DefaultLayout
      loading={tableData.isLoading}
      title={getString('incidents')}
      toolbar={reportIncidentButton}
      subHeader={
        <>
          <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} spacing="medium">
            {incidentsStatusFilter}
            {incidentsSeverityFilter}
          </Layout.Horizontal>
          <FlexExpander />
          {incidentsSearchBar}
        </>
      }
      footer={tableData.pagination && <Pagination {...tableData.pagination} />}
      noData={!isDataPresent}
      noDataProps={{
        height: 350,
        ...(areFiltersSet
          ? {
              title: getString('noIncidentsFoundMatchingFilters'),
              subtitle: getString('noIncidentsFoundMatchingFiltersDescription'),
              ctaButton: resetFiltersButton
            }
          : {
              title: getString('noIncidentsFound'),
              subtitle: getString('noIncidentsFoundDescription'),
              ctaButton: reportIncidentButton
            })
      }}
    >
      <MemoisedIncidentListTable content={tableData.content} />
    </DefaultLayout>
  );
};

export default IncidentsView;
