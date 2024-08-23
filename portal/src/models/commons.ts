export interface SortInput<T> {
  field: T;
  ascending?: boolean;
}

export enum IncidentsSortType {
  NAME = 'NAME',
  DURATION = 'DURATION'
}
