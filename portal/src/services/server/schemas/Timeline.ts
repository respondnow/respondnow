/* eslint-disable */
// This code is autogenerated using @harnessio/oats-cli.
// Please do not modify this code directly.
import type { Slack } from '../schemas/Slack';
import type { UserDetails } from '../schemas/UserDetails';

export interface Timeline {
  additionalDetails?: { [key: string]: any };
  /**
   * @format int64
   */
  createdAt?: number;
  currentState?: string;
  id?: string;
  message?: string;
  previousState?: string;
  slack?: Slack;
  type: 'Comment' | 'Incident_Created' | 'Roles' | 'Severity' | 'Slack_Channel_Created' | 'Status' | 'Summary';
  /**
   * @format int64
   */
  updatedAt?: number;
  userDetails?: UserDetails;
}
