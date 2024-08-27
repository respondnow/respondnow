import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Avatar, Layout, Tag, Text } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import { Icon } from '@harnessio/icons';
import { Fallback } from '@errors';
import { useStrings } from '@strings';
import SeverityBadge from '@components/SeverityBadge';
import StatusBadge from '@components/StatusBadge';
import { IncidentIncident } from '@services/server';
import { getDurationBasedOnStatus } from '@utils';
import css from '../IncidentDetails.module.scss';

interface DetailsSectionProps {
  incidentData: IncidentIncident | undefined;
}

const DetailsSection: React.FC<DetailsSectionProps> = props => {
  const { incidentData } = props;
  const { getString } = useStrings();

  return (
    <Layout.Vertical width="100%" height="100%" padding="medium" className={css.detailsSectionContainer}>
      <Layout.Vertical
        flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
        className={css.internalContainers}
        width="100%"
      >
        <Text font={{ variation: FontVariation.H6 }}>{getString('severity')}</Text>
        <SeverityBadge severity={incidentData?.severity} />
      </Layout.Vertical>
      <Layout.Horizontal width="100%">
        <Layout.Vertical
          flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
          className={css.internalContainers}
          width="50%"
        >
          <Text font={{ variation: FontVariation.H6 }}>{getString('status')}</Text>
          <StatusBadge status={incidentData?.status} />
        </Layout.Vertical>
        <Layout.Vertical
          flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
          className={css.internalContainers}
          width="50%"
        >
          <Text font={{ variation: FontVariation.H6 }}>{getString('duration')}</Text>
          <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} width="100%">
            <Icon name="time" size={12} color={Color.GREY_800} margin={{ right: 'small' }} />
            {incidentData?.createdAt && incidentData?.updatedAt && incidentData.status && (
              <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
                {getDurationBasedOnStatus(
                  incidentData.createdAt * 1000,
                  incidentData.updatedAt * 1000,
                  incidentData.status
                )}
              </Text>
            )}
          </Layout.Horizontal>
        </Layout.Vertical>
      </Layout.Horizontal>
      <Layout.Vertical
        flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
        className={css.internalContainers}
        width="100%"
      >
        <Text font={{ variation: FontVariation.H6 }}>{getString('summary')}</Text>
        <Text font={{ variation: FontVariation.BODY }} color={Color.GREY_800}>
          {incidentData?.summary || getString('abbv.na')}
        </Text>
      </Layout.Vertical>
      <Layout.Vertical
        flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
        className={css.internalContainers}
        width="100%"
      >
        <Text font={{ variation: FontVariation.H6 }}>{getString('keyMembers')}</Text>
        <Layout.Vertical width="100%">
          {incidentData?.roles?.map(member => (
            <Layout.Horizontal
              key={member.userDetails?.userName}
              flex={{ alignItems: 'center', justifyContent: 'flex-start' }}
            >
              <Avatar hoverCard={false} name={member.userDetails?.name || member.userDetails?.userName} size="small" />
              <Layout.Vertical>
                <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.GREY_800}>
                  {member.userDetails?.name || member.userDetails?.userName}
                </Text>
                <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_600}>
                  {member.roleType}
                </Text>
              </Layout.Vertical>
            </Layout.Horizontal>
          ))}
        </Layout.Vertical>
      </Layout.Vertical>
      <Layout.Vertical
        flex={{ alignItems: 'flex-start', justifyContent: 'flex-start' }}
        className={css.internalContainers}
        width="100%"
      >
        <Text font={{ variation: FontVariation.H6 }}>{getString('tags')}</Text>
        <Layout.Horizontal width="100%" style={{ flexWrap: 'wrap' }}>
          {incidentData?.tags && incidentData.tags?.length > 0 ? (
            incidentData?.tags?.map(tag => (
              <Tag key={tag} style={{ marginRight: '0.25rem' }}>
                {tag}
              </Tag>
            ))
          ) : (
            <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
              {getString('abbv.na')}
            </Text>
          )}
        </Layout.Horizontal>
      </Layout.Vertical>
    </Layout.Vertical>
  );
};

export default withErrorBoundary(DetailsSection, { FallbackComponent: Fallback });
