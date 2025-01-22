import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Layout, Text, Utils } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import { Icon } from '@harnessio/icons';
import { Fallback } from '@errors';
import { Incident } from '@services/server';
import css from './SeverityBadge.module.scss';

const SeverityBadge: React.FC<{ severity: Incident['severity'] | undefined }> = ({ severity }) => {
  const getSeverityProps = (): {
    foregroundColor: string;
    text: string;
  } => {
    switch (severity) {
      case 'SEV0':
        return {
          foregroundColor: Color.RED_700,
          text: 'SEV0'
        };
      case 'SEV1':
        return {
          foregroundColor: Color.ORANGE_700,
          text: 'SEV1'
        };
      case 'SEV2':
        return {
          foregroundColor: Color.GREEN_700,
          text: 'SEV2'
        };
      default:
        return {
          foregroundColor: Color.GREY_500,
          text: 'N/A'
        };
    }
  };

  const { foregroundColor, text } = getSeverityProps();

  return (
    <Layout.Horizontal
      flex={{ align: 'center-center' }}
      background={Color.WHITE}
      border
      style={{ borderColor: Utils.getRealCSSColor(foregroundColor) }}
      className={css.badgeContainer}
    >
      <Icon name="full-circle" size={6} color={foregroundColor} />
      <Text font={{ variation: FontVariation.TINY_SEMI }} color={foregroundColor}>
        {text}
      </Text>
    </Layout.Horizontal>
  );
};

export default withErrorBoundary(SeverityBadge, { FallbackComponent: Fallback });
