/* eslint-disable */
// This code is autogenerated using @harnessio/oats-cli.
// Please do not modify this code directly.
import type { IncidentStatus } from '../schemas/IncidentStatus';
import type { UtilsUserDetails } from '../schemas/UtilsUserDetails';

export interface IncidentStage {
  createdAt?: number;
  duration?: number;
  stageId?: string;
  type?: IncidentStatus;
  updatedAt?: number;
  userDetails?: UtilsUserDetails;
}
