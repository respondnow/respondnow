import { isEqual } from 'lodash-es';
import React from 'react';
import { TableV2 } from '@harnessio/uicore';
import { Column } from 'react-table';
import cx from 'classnames';
import { IncidentsTableProps } from '@interfaces';
import { IncidentIncident } from '@services/server';
import { useStrings } from '@strings';
import * as CellRenderer from './CellRenderer';
import css from '../CommonTableStyles.module.scss';

const IncidentListTable: React.FC<IncidentsTableProps> = props => {
  const { content } = props;
  const { getString } = useStrings();

  const columns: Column<IncidentIncident>[] = React.useMemo(() => {
    return [
      {
        Header: getString('incident'),
        id: 'name',
        Cell: CellRenderer.IncidentsName,
        width: '35%'
      },
      {
        Header: getString('reportedBy'),
        id: 'reportedBy',
        Cell: CellRenderer.IncidentReportedBy,
        width: '20%'
      },
      {
        Header: getString('status'),
        id: 'status',
        Cell: CellRenderer.IncidentStatus,
        width: '15%'
      },
      {
        Header: getString('duration'),
        id: 'duration',
        Cell: CellRenderer.IncidentDuration,
        width: '15%'
      },
      {
        Header: '',
        id: 'cta',
        Cell: CellRenderer.IncidentCTA,
        width: '15%'
      }
    ];
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return <TableV2<IncidentIncident> columns={columns} data={content} className={cx(css.paginationFix)} />;
};

const MemoisedIncidentListTable = React.memo(IncidentListTable, (prev, current) => {
  return isEqual(prev, current);
});

export default MemoisedIncidentListTable;
