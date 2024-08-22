import React from 'react';
import { Container, Text } from '@harnessio/uicore';
import { FontVariation } from '@harnessio/design-system';
import noFilteredData from '@images/no-data.svg';
import css from './NoData.module.scss';

export interface NoDataProps {
  height?: string | number;
  width?: string | number;
  title?: string;
  image?: string;
  subtitle?: string;
  ctaButton?: React.ReactElement;
  documentationLink?: React.ReactElement;
}

export default function NoData({
  height = 200,
  title,
  image,
  subtitle,
  ctaButton,
  documentationLink
}: NoDataProps): React.ReactElement {
  return (
    <Container height={height} className={css.container}>
      <div className={css.image}>
        <img data-testid="no-data-img" src={image ?? noFilteredData} />
      </div>
      <div className={css.content}>
        <Text margin={{ top: 'medium' }} font={{ variation: FontVariation.H5 }}>
          {title}
        </Text>
        <Text margin={{ top: 'small', bottom: 'medium' }} font={{ variation: FontVariation.BODY }}>
          {subtitle}
        </Text>
        {documentationLink}
        {ctaButton}
      </div>
    </Container>
  );
}
