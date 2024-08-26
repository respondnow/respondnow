import React from 'react';
import { CellProps, Renderer } from 'react-table';
import { Avatar, Button, ButtonVariation, Layout, Tag, Text } from '@harnessio/uicore';
import { Link } from 'react-router-dom';
import { Color, FontVariation } from '@harnessio/design-system';
import { Icon } from '@harnessio/icons';
import { IncidentIncident } from '@services/server';
import { useStrings } from '@strings';
import SeverityBadge from '@components/SeverityBadge';
import { paths } from '@routes/RouteDefinitions';
import { getDetailedTime, getDurationBasedOnStatus } from '@utils';
import StatusBadge from '@components/StatusBadge';
import css from '../CommonTableStyles.module.scss';

type CellRendererType = Renderer<CellProps<IncidentIncident>>;

export const IncidentsName: CellRendererType = ({ row }) => {
  const { getString } = useStrings();
  const { name, severity, description, tags, identifier } = row.original;
  return (
    <Layout.Horizontal
      flex={{ alignItems: 'center', justifyContent: 'flex-start' }}
      width="100%"
      style={{ gap: '0.5rem' }}
    >
      <SeverityBadge severity={severity} />
      <Layout.Vertical className={css.incidentsNameContainer}>
        <Link
          to={paths.toIncidentDetails({
            incidentId: identifier
          })}
          className={css.textLink}
        >
          <Text color={Color.PRIMARY_7} font={{ variation: FontVariation.BODY, weight: 'bold' }} lineClamp={1}>
            {name || getString('abbv.na')}
          </Text>
        </Link>
        {description && (
          <Text font={{ variation: FontVariation.SMALL }} lineClamp={1}>
            {description}
          </Text>
        )}
        {tags &&
          (tags.length < 3 ? (
            <Layout.Horizontal style={{ gap: '0.25rem' }}>
              {tags.map(tag => (
                <Tag key={tag}>{tag}</Tag>
              ))}
            </Layout.Horizontal>
          ) : (
            <Layout.Horizontal style={{ gap: '0.25rem' }}>
              {tags.slice(0, 3).map(tag => (
                <Tag key={tag}>{tag}</Tag>
              ))}
              <Tag>+{tags.length - 1}</Tag>
            </Layout.Horizontal>
          ))}
      </Layout.Vertical>
    </Layout.Horizontal>
  );
};

export const IncidentReportedBy: CellRendererType = ({ row }) => {
  const { getString } = useStrings();
  const { createdBy, createdAt } = row.original;
  return (
    <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} width="100%">
      <Avatar name={createdBy?.name} hoverCard={false} size="small" />
      <Layout.Vertical className={css.incidentsNameContainer}>
        <Text font={{ variation: FontVariation.SMALL }} lineClamp={1} color={Color.GREY_800}>
          {createdBy?.name || getString('abbv.na')}
        </Text>
        {createdAt && (
          <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_500}>
            {getDetailedTime(createdAt * 1000, true)}
          </Text>
        )}
      </Layout.Vertical>
    </Layout.Horizontal>
  );
};

export const IncidentStatus: CellRendererType = ({ row }) => {
  const { status } = row.original;
  return (
    <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }}>
      <StatusBadge status={status} />
    </Layout.Horizontal>
  );
};

export const IncidentDuration: CellRendererType = ({ row }) => {
  const { createdAt, updatedAt, status } = row.original;
  return (
    <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }} width="100%" spacing="small">
      <Icon name="time" size={12} color={Color.GREY_800} />
      {createdAt && updatedAt && (
        <Text font={{ variation: FontVariation.SMALL }} color={Color.GREY_800}>
          {getDurationBasedOnStatus(createdAt * 1000, updatedAt * 1000, status)}
        </Text>
      )}
    </Layout.Horizontal>
  );
};

export const IncidentCTA: CellRendererType = () => {
  const { getString } = useStrings();
  return (
    <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-end' }}>
      <Button text={getString('viewChannel')} variation={ButtonVariation.SECONDARY} />
    </Layout.Horizontal>
  );
};
