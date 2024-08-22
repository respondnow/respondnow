import React, { useState } from 'react';
import { Color, FontVariation } from '@harnessio/design-system';
import { Button, ButtonVariation, Layout, Text } from '@harnessio/uicore';
import cx from 'classnames';
import css from './InfoBanner.module.scss';

type BannerType = 'info' | 'levelup';

export interface InfoBannerProps {
  title?: string;
  message: React.ReactNode;
  type: BannerType;
}

export const LevelUpText: React.FC<Pick<InfoBannerProps, 'title' | 'message'>> = ({ title, message }) => {
  return (
    <Layout.Horizontal flex={{ alignItems: 'center' }} padding={{ right: 'small' }}>
      <Text
        icon="flash"
        color={Color.ORANGE_800}
        font={{ variation: FontVariation.FORM_MESSAGE_WARNING, weight: 'bold' }}
        iconProps={{ color: Color.ORANGE_800, size: 25 }}
        padding={{ right: 'medium' }}
      >
        {title}
      </Text>
      <Text color={Color.PRIMARY_10} font={{ variation: FontVariation.SMALL, weight: 'semi-bold' }}>
        {message}
      </Text>
    </Layout.Horizontal>
  );
};

export const InfoText: React.FC<Pick<InfoBannerProps, 'message'>> = ({ message }) => {
  return (
    <Text
      icon="info-message"
      color={Color.PRIMARY_10}
      font={{ variation: FontVariation.SMALL, weight: 'semi-bold' }}
      iconProps={{ padding: { right: 'small' }, size: 20, className: css.infoIcon }}
    >
      {message}
    </Text>
  );
};

const InfoBanner: React.FC<InfoBannerProps> = ({ message, type }) => {
  const [display, setDisplay] = useState(true);

  if (!display) {
    return <></>;
  }

  return (
    <div className={cx(css.infoBanner, { [css.levelUp]: type === 'levelup', [css.info]: type === 'info' })}>
      <Layout.Horizontal spacing="medium" width="95%">
        {type === 'info' ? <InfoText message={message} /> : <LevelUpText message={message} />}
      </Layout.Horizontal>
      <Button
        variation={ButtonVariation.ICON}
        icon="cross"
        data-testid="cross-button"
        iconProps={{ size: 18 }}
        onClick={() => {
          setDisplay(false);
        }}
      />
    </div>
  );
};

export default InfoBanner;
