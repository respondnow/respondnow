import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Button, Card, Layout, Text, useToggleOpen } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import cx from 'classnames';
import { Fallback } from '@errors';
import { useStrings } from '@strings';
import { Incident } from '@services/server';
import IncidentTimeline from '@components/IncidentTimeline';
import css from '../IncidentDetails.module.scss';

interface TimelineSectionProps {
  incidentData: Incident | undefined;
}

const TimelineSection: React.FC<TimelineSectionProps> = props => {
  const { incidentData } = props;
  const { getString } = useStrings();
  const { isOpen: showCommentsOnly, toggle } = useToggleOpen();

  return (
    <Card className={css.timelineCardContainer}>
      <Layout.Vertical
        flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
        width="100%"
        height="100%"
        padding="medium"
        className={css.timelineSectionContainer}
      >
        <Layout.Horizontal
          width="100%"
          flex={{ alignItems: 'center', justifyContent: 'space-between' }}
          spacing="large"
        >
          <Text font={{ variation: FontVariation.H5 }}>{getString('incidentTimeline')}</Text>
          <Button
            noStyling
            onClick={toggle}
            className={cx(css.filterButton, {
              [css.active]: showCommentsOnly
            })}
            width={180}
          >
            <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
              {showCommentsOnly ? getString('showAll') : getString('showCommentsOnly')}
            </Text>
          </Button>
        </Layout.Horizontal>
        <IncidentTimeline incident={incidentData} showCommentsOnly={showCommentsOnly} />
      </Layout.Vertical>
    </Card>
  );
};

export default withErrorBoundary(TimelineSection, { FallbackComponent: Fallback });
