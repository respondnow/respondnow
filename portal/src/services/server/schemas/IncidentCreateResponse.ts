/* eslint-disable */
// This code is autogenerated using @harnessio/oats-cli.
// Please do not modify this code directly.
import type { IncidentAttachment } from '../schemas/IncidentAttachment';
import type { IncidentChannel } from '../schemas/IncidentChannel';
import type { IncidentConference } from '../schemas/IncidentConference';
import type { UtilsUserDetails } from '../schemas/UtilsUserDetails';
import type { IncidentEnvironment } from '../schemas/IncidentEnvironment';
import type { IncidentFunctionality } from '../schemas/IncidentFunctionality';
import type { IncidentIncidentChannel } from '../schemas/IncidentIncidentChannel';
import type { IncidentRole } from '../schemas/IncidentRole';
import type { IncidentService } from '../schemas/IncidentService';
import type { IncidentSeverity } from '../schemas/IncidentSeverity';
import type { IncidentStage } from '../schemas/IncidentStage';
import type { IncidentStatus } from '../schemas/IncidentStatus';
import type { IncidentTimeline } from '../schemas/IncidentTimeline';
import type { IncidentType } from '../schemas/IncidentType';

export interface IncidentCreateResponse {
  accountIdentifier?: string;
  active: boolean;
  attachments?: IncidentAttachment[];
  channels?: IncidentChannel[];
  conferenceDetails?: IncidentConference[];
  correlationID?: string;
  createdAt?: number;
  createdBy?: UtilsUserDetails;
  description?: string;
  environments?: IncidentEnvironment[];
  functionalities?: IncidentFunctionality[];
  id: string;
  identifier: string;
  incidentChannel?: IncidentIncidentChannel;
  name: string;
  orgIdentifier?: string;
  projectIdentifier?: string;
  removed?: boolean;
  removedAt?: number;
  roles?: IncidentRole[];
  services?: IncidentService[];
  severity: IncidentSeverity;
  stages?: IncidentStage[];
  status: IncidentStatus;
  summary: string;
  tags?: string[];
  timelines?: IncidentTimeline[];
  type?: IncidentType;
  updatedAt?: number;
  updatedBy?: UtilsUserDetails;
}
