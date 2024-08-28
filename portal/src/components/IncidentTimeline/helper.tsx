import React from 'react';
import { Icon } from '@harnessio/icons';
import { Color } from '@harnessio/design-system';
import { TimelineUtilReturn } from '@interfaces';
import { IncidentTimeline } from '@services/server';

export function getTimelinePropsBasedOnIncidentData({ type }: IncidentTimeline): TimelineUtilReturn | undefined {
  switch (type) {
    case 'addComment':
      return {
        icon: (
          <Icon
            name="chaos-litmuschaos"
            size={12}
            background={Color.GREY_100}
            style={{ borderRadius: '10rem', padding: '7px' }}
          />
        ),
        headerContent: 'Comment',
        bodyContent: 'Body'
      };
    case 'updateSeverity':
      return {
        icon: (
          <Icon
            name="chaos-litmuschaos"
            size={12}
            background={Color.GREY_100}
            style={{ borderRadius: '10rem', padding: '7px' }}
          />
        ),
        headerContent: 'Severity',
        bodyContent: 'Body'
      };
    case 'updateStatus':
      return {
        icon: (
          <Icon
            name="chaos-litmuschaos"
            size={12}
            background={Color.GREY_100}
            style={{ borderRadius: '10rem', padding: '7px' }}
          />
        ),
        headerContent: 'Status',
        bodyContent: 'Body'
      };
    default:
      return undefined;
  }
}
