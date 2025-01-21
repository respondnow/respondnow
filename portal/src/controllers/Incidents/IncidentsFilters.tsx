/*
 * Copyright 2022 Harness Inc. All rights reserved.
 * Use of this source code is governed by the PolyForm Shield 1.0.0 license
 * that can be found in the licenses directory at the root of this repository, also available at
 * https://polyformproject.org/wp-content/uploads/2020/06/PolyForm-Shield-1.0.0.txt.
 */

import React from 'react';
import {
  Button,
  ButtonVariation,
  ExpandingSearchInputHandle,
  ExpandingSearchInput,
  Layout,
  Text,
  SelectOption,
  DropDown
} from '@harnessio/uicore';
import { Color } from '@harnessio/design-system';
import { Icon } from '@harnessio/icons';
import { useStrings } from '@strings';
import { IncidentsFilter, IncidentsFilterAction, IncidentsFilterActionKind } from 'hooks';
import { IncidentsSortType } from 'models';
import { Incident } from '@services/server';

export interface FilterProps {
  state: IncidentsFilter;
  resetPage: () => void;
  dispatch: React.Dispatch<IncidentsFilterAction>;
}

export const SortNameHeader = ({ state, dispatch }: FilterProps): React.ReactElement => {
  const { getString } = useStrings();

  return (
    <Layout.Horizontal
      flex={{ alignItems: 'center', justifyContent: 'flex-start' }}
      style={{ gap: '0.25rem', cursor: 'pointer' }}
      onClick={() => {
        dispatch({
          type: IncidentsFilterActionKind.CHANGE_SORT_TYPE,
          payload: {
            sortType: {
              field: IncidentsSortType.NAME,
              ascending: state.sortType?.field === IncidentsSortType.NAME ? !state.sortType?.ascending : false
            }
          }
        });
      }}
    >
      <Text color={Color.GREY_900}>{getString('incidents').toUpperCase()}</Text>
      {state.sortType?.field === IncidentsSortType.NAME && (
        <Icon
          name={state.sortType.ascending ? 'main-chevron-up' : 'main-chevron-down'}
          size={8}
          color={Color.GREY_900}
        />
      )}
    </Layout.Horizontal>
  );
};

export const IncidentsSearchBar = ({ state, dispatch, resetPage }: FilterProps): React.ReactElement => {
  const { getString } = useStrings();
  const ref = React.useRef<ExpandingSearchInputHandle | undefined>();
  React.useEffect(() => {
    if (state.incidentName === '' && ref.current) {
      ref.current.clear();
    }
  }, [state.incidentName]);
  return (
    <ExpandingSearchInput
      ref={ref}
      width={300}
      alwaysExpanded
      placeholder={getString('searchForAnIncident')}
      throttle={500}
      autoFocus={false}
      onChange={incidentName => {
        if (!(state.incidentName === incidentName)) {
          resetPage();
          dispatch({
            type: IncidentsFilterActionKind.CHANGE_INCIDENTS_NAME,
            payload: {
              incidentName
            }
          });
        }
      }}
    />
  );
};

export const IncidentsStatusFilter = ({ state, dispatch, resetPage }: FilterProps): React.ReactElement => {
  const { getString } = useStrings();
  const dropdownItems: SelectOption[] = [
    {
      label: 'Acknowledged',
      value: 'Acknowledged'
    },
    {
      label: 'Identified',
      value: 'Identified'
    },
    {
      label: 'Investigating',
      value: 'Investigating'
    },
    {
      label: 'Mitigated',
      value: 'Mitigated'
    },
    {
      label: 'Resolved',
      value: 'Resolved'
    },
    {
      label: 'Started',
      value: 'Started'
    }
  ];

  const handleChange = (incidentStatus: Incident['status']): void => {
    resetPage();
    dispatch({
      type: IncidentsFilterActionKind.CHANGE_INCIDENTS_STATUS,
      payload: {
        incidentStatus
      }
    });
  };

  return (
    <DropDown
      addClearBtn
      filterable={false}
      items={dropdownItems}
      onChange={value => handleChange(String(value.value) as Incident['status'])}
      value={state.incidentStatus}
      placeholder={getString('status')}
      width={250}
    />
  );
};

export const IncidentsSeverityFilter = ({ state, dispatch, resetPage }: FilterProps): React.ReactElement => {
  const { getString } = useStrings();
  const dropdownItems: SelectOption[] = [
    {
      label: 'SEV0 - Critical, High Impact',
      value: 'SEV0 - Critical, High Impact'
    },
    {
      label: 'SEV1 - Major, Significant Impact',
      value: 'SEV1 - Major, Significant Impact'
    },
    {
      label: 'SEV2 - Minor, Low Impact',
      value: 'SEV2 - Minor, Low Impact'
    }
  ];

  const handleChange = (incidentSeverity: Incident['severity']): void => {
    resetPage();
    dispatch({
      type: IncidentsFilterActionKind.CHANGE_INCIDENTS_SEVERITY,
      payload: {
        incidentSeverity
      }
    });
  };

  return (
    <DropDown
      addClearBtn
      filterable={false}
      items={dropdownItems}
      onChange={value => handleChange(String(value.value) as Incident['severity'])}
      value={state.incidentSeverity}
      placeholder={getString('severity')}
      width={250}
    />
  );
};

interface ResetButtonProps extends FilterProps {
  minimal?: boolean;
}

export const ResetFilterButton = ({ dispatch, resetPage, minimal = true }: ResetButtonProps): React.ReactElement => {
  const { getString } = useStrings();
  return (
    <Button
      icon="reset"
      variation={minimal ? ButtonVariation.LINK : ButtonVariation.SECONDARY}
      onClick={() => {
        resetPage();
        dispatch({
          type: IncidentsFilterActionKind.RESET_FILTERS,
          payload: {}
        });
      }}
      text={getString('resetFilters')}
    />
  );
};
