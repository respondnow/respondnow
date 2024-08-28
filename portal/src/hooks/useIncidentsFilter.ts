import React, { Dispatch } from 'react';
import { IncidentsSortType, SortInput } from 'models';
import { IncidentSeverity, IncidentStatus } from '@services/server';

export enum IncidentsFilterActionKind {
  CHANGE_INCIDENTS_NAME = 'CHANGE_INCIDENTS_NAME',
  CHANGE_INCIDENTS_STATUS = 'CHANGE_INCIDENTS_STATUS',
  CHANGE_INCIDENTS_SEVERITY = 'CHANGE_INCIDENTS_SEVERITY',
  CHANGE_INCIDENTS_TAGS = 'CHANGE_INCIDENTS_TAGS',
  CHANGE_INCIDENTS_TIMEFRAME = 'CHANGE_INCIDENTS_TIMEFRAME',
  CHANGE_SORT_TYPE = 'CHANGE_SORT_TYPE',
  RESET_FILTERS = 'RESET_FILTERS'
}

export interface IncidentsFilter {
  incidentName?: string | undefined;
  incidentStatus?: IncidentStatus;
  incidentSeverity?: IncidentSeverity;
  incidentTags?: string[] | undefined;
  incidentTimeframe?: string | undefined;
  sortType?: SortInput<IncidentsSortType> | undefined;
}

export interface IncidentsFilterAction {
  type: IncidentsFilterActionKind;
  payload: IncidentsFilter;
}

interface ReducerReturn {
  state: IncidentsFilter;
  dispatch: Dispatch<IncidentsFilterAction>;
}

export const initialIncidenrsFilterState: IncidentsFilter = {
  incidentName: '',
  incidentStatus: undefined,
  incidentTags: undefined,
  incidentTimeframe: undefined,
  sortType: undefined
};

function reducer(state: IncidentsFilter, action: IncidentsFilterAction): ReducerReturn['state'] {
  switch (action.type) {
    case IncidentsFilterActionKind.CHANGE_INCIDENTS_NAME:
      return { ...state, incidentName: action.payload.incidentName };
    case IncidentsFilterActionKind.CHANGE_INCIDENTS_STATUS:
      return { ...state, incidentStatus: action.payload.incidentStatus };
    case IncidentsFilterActionKind.CHANGE_INCIDENTS_SEVERITY:
      return { ...state, incidentSeverity: action.payload.incidentSeverity };
    case IncidentsFilterActionKind.CHANGE_INCIDENTS_TAGS:
      return { ...state, incidentTags: action.payload.incidentTags };
    case IncidentsFilterActionKind.CHANGE_INCIDENTS_TIMEFRAME:
      return { ...state, incidentTimeframe: action.payload.incidentTimeframe };
    case IncidentsFilterActionKind.CHANGE_SORT_TYPE:
      return { ...state, sortType: action.payload.sortType };
    case IncidentsFilterActionKind.RESET_FILTERS:
      return { ...initialIncidenrsFilterState };
    default:
      throw new Error();
  }
}

export function useIncidentsFilter(): ReducerReturn {
  const [state, dispatch] = React.useReducer(reducer, initialIncidenrsFilterState);

  return { state, dispatch };
}
