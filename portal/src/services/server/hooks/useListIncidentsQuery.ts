/* eslint-disable */
// This code is autogenerated using @harnessio/oats-cli.
// Please do not modify this code directly.
import { useQuery, UseQueryOptions } from '@tanstack/react-query';

import type { IncidentListResponseDto } from '../schemas/IncidentListResponseDto';
import type { UtilsDefaultResponseDto } from '../schemas/UtilsDefaultResponseDto';
import { fetcher, FetcherOptions } from '@services/fetcher';

export interface ListIncidentsQueryQueryParams {
  accountIdentifier: string;
  orgIdentifier?: string;
  projectIdentifier?: string;
  type?: 'Availability' | 'Latency' | 'Other' | 'Security';
  severity?: 'SEV0 - Critical, High Impact' | 'SEV1 - Major, Significant Impact' | 'SEV2 - Minor, Low Impact';
  status?: 'Acknowledged' | 'Identified' | 'Investigating' | 'Mitigated' | 'Resolved' | 'Started';
  active?: boolean;
  incidentChannelType?: 'slack';
  search?: string;
  page?: number;
  limit?: number;
  correlationId?: string;
  all?: boolean;
}

export type ListIncidentsOkResponse = IncidentListResponseDto;

export type ListIncidentsErrorResponse = UtilsDefaultResponseDto;

export interface ListIncidentsProps extends Omit<FetcherOptions<ListIncidentsQueryQueryParams, unknown>, 'url'> {
  queryParams: ListIncidentsQueryQueryParams;
}

export function listIncidents(props: ListIncidentsProps): Promise<ListIncidentsOkResponse> {
  return fetcher<ListIncidentsOkResponse, ListIncidentsQueryQueryParams, unknown>({
    url: `/api/incident/list`,
    method: 'GET',
    ...props
  });
}

/**
 * List incidents
 */
export function useListIncidentsQuery(
  props: ListIncidentsProps,
  options?: Omit<UseQueryOptions<ListIncidentsOkResponse, ListIncidentsErrorResponse>, 'queryKey' | 'queryFn'>
) {
  return useQuery<ListIncidentsOkResponse, ListIncidentsErrorResponse>(
    ['ListIncidents', props.queryParams],
    ({ signal }) => listIncidents({ ...props, signal }),
    options
  );
}
