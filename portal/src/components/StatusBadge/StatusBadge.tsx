import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Icon, IconName } from '@harnessio/icons';
import { Color, FontVariation } from '@harnessio/design-system';
import { Layout, Text } from '@harnessio/uicore';
import { Fallback } from '@errors';
import { Incident } from '@services/server';
import css from './StatusBadge.module.scss';

const StatusBadge: React.FC<{ status: Incident['status'] | undefined }> = ({ status }) => {
  const getSeverityProps = (): {
    foregroundColor: string;
    icon: IconName;
  } => {
    switch (status) {
      case 'Acknowledged':
        return {
          foregroundColor: Color.PRIMARY_7,
          icon: 'status-pending'
        };
      case 'Resolved':
        return {
          foregroundColor: Color.GREEN_700,
          icon: 'tick'
        };
      case 'Identified':
        return {
          foregroundColor: Color.PRIMARY_7,
          icon: 'status-pending'
        };
      case 'Investigating':
        return {
          foregroundColor: Color.PRIMARY_7,
          icon: 'status-pending'
        };
      case 'Mitigated':
        return {
          foregroundColor: Color.PRIMARY_7,
          icon: 'status-pending'
        };
      case 'Started':
        return {
          foregroundColor: Color.PRIMARY_7,
          icon: 'status-pending'
        };
      default:
        return {
          foregroundColor: Color.BLACK,
          icon: 'status-pending'
        };
    }
  };

  const { foregroundColor, icon } = getSeverityProps();

  return (
    <Layout.Horizontal flex={{ align: 'center-center' }} className={css.badgeContainer}>
      <Icon name={icon} size={14} color={foregroundColor} />
      <Text font={{ variation: FontVariation.SMALL }} color={foregroundColor}>
        {status || 'N/A'}
      </Text>
    </Layout.Horizontal>
  );
};

export default withErrorBoundary(StatusBadge, { FallbackComponent: Fallback });
