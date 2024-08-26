import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Icon, IconName } from '@harnessio/icons';
import { Color, FontVariation } from '@harnessio/design-system';
import { Layout, Text } from '@harnessio/uicore';
import { Fallback } from '@errors';
import { IncidentStatus } from '@services/server';
import css from './StatusBadge.module.scss';

const StatusBadge: React.FC<{ status: IncidentStatus | undefined }> = ({ status }) => {
  const getSeverityProps = (): {
    foreGroundColor: string;
    icon: IconName;
  } => {
    switch (status) {
      case 'Acknowledged':
        return {
          foreGroundColor: Color.PRIMARY_7,
          icon: 'steps-spinner'
        };
      case 'Resolved':
        return {
          foreGroundColor: Color.GREEN_700,
          icon: 'tick'
        };
      case 'Identified':
        return {
          foreGroundColor: Color.PRIMARY_7,
          icon: 'steps-spinner'
        };
      case 'Investigating':
        return {
          foreGroundColor: Color.PRIMARY_7,
          icon: 'steps-spinner'
        };
      case 'Mitigated':
        return {
          foreGroundColor: Color.PRIMARY_7,
          icon: 'steps-spinner'
        };
      case 'Started':
        return {
          foreGroundColor: Color.PRIMARY_7,
          icon: 'steps-spinner'
        };
      default:
        return {
          foreGroundColor: Color.BLACK,
          icon: 'steps-spinner'
        };
    }
  };

  const { foreGroundColor, icon } = getSeverityProps();

  return (
    <Layout.Horizontal flex={{ align: 'center-center' }} background={Color.GREY_50} className={css.badgeContainer}>
      <Icon name={icon} size={14} color={foreGroundColor} />
      <Text font={{ variation: FontVariation.SMALL }} color={foreGroundColor}>
        {status || 'N/A'}
      </Text>
    </Layout.Horizontal>
  );
};

export default withErrorBoundary(StatusBadge, { FallbackComponent: Fallback });
