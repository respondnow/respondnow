import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Layout, Text } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import { Icon } from '@harnessio/icons';
import { Fallback } from '@errors';
import { IncidentSeverity } from '@services/server';
import css from './SeverityBadge.module.scss';

const SeverityBadge: React.FC<{ severity: IncidentSeverity | undefined }> = ({ severity }) => {
  const getSeverityProps = (): {
    foreGroundColor: string;
    text: string;
  } => {
    switch (severity) {
      case 'SEV0 - Critical, High Impact':
        return {
          foreGroundColor: Color.RED_700,
          text: 'SEV0'
        };
      case 'SEV1 - Major, Significant Impact':
        return {
          foreGroundColor: Color.ORANGE_700,
          text: 'SEV1'
        };
      case 'SEV2 - Minor, Low Impact':
        return {
          foreGroundColor: Color.GREEN_700,
          text: 'SEV2'
        };
      default:
        return {
          foreGroundColor: Color.GREY_500,
          text: 'N/A'
        };
    }
  };

  const { foreGroundColor, text } = getSeverityProps();

  return (
    <Layout.Horizontal
      flex={{ align: 'center-center' }}
      background={Color.WHITE}
      border={{ color: foreGroundColor }}
      className={css.badgeContainer}
    >
      <Icon name="full-circle" size={6} color={foreGroundColor} />
      <Text font={{ variation: FontVariation.TINY_SEMI }} color={foreGroundColor}>
        {text}
      </Text>
    </Layout.Horizontal>
  );
};

export default withErrorBoundary(SeverityBadge, { FallbackComponent: Fallback });
