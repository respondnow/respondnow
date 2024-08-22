import React from 'react';
import { withErrorBoundary } from 'react-error-boundary';
import { Avatar, Button, ButtonSize, Container, Layout, Text, TextProps } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import { NavLink, NavLinkProps } from 'react-router-dom';
import { IconName } from '@harnessio/icons';
import cx from 'classnames';
import { PopoverInteractionKind, Position } from '@blueprintjs/core';
import { Fallback } from '@errors';
import respondNowLogo from '@images/respondNow.svg';
import { useStrings } from '@strings';
import { useRouteWithBaseUrl } from '@hooks';
import { useLogout } from '@utils';
import css from './SideNav.module.scss';

interface SidebarLinkProps extends NavLinkProps {
  label: string;
  icon?: IconName;
  className?: string;
  textProps?: TextProps;
}

export const SidebarLink: React.FC<SidebarLinkProps> = ({ label, icon, className, textProps, ...others }) => (
  <NavLink className={cx(css.link, className)} activeClassName={css.selected} {...others}>
    <Text icon={icon} className={css.text} font={{ variation: FontVariation.SMALL_SEMI }} {...textProps}>
      {label}
    </Text>
  </NavLink>
);

const SideNav: React.FC = () => {
  const { getString } = useStrings();
  const paths = useRouteWithBaseUrl();
  const { forceLogout } = useLogout();

  return (
    <Layout.Vertical height="100%" width={250} background={Color.PRIMARY_BG} className={css.sideNavMainContainer}>
      <Container padding="medium">
        <img src={respondNowLogo} alt="RespondNow" height={25} />
      </Container>
      <Layout.Vertical padding="medium" className={css.sideNavLinkContainer}>
        <SidebarLink label={getString('incidents')} to={paths.toDummy()} icon="home" />
      </Layout.Vertical>
      <Layout.Vertical>
        <Container padding="medium" border={{ top: true }}>
          <SidebarLink
            to="#"
            icon="link"
            label={getString('documentation')}
            target="_blank"
            rel="noopener noreferrer"
            activeClassName={undefined}
          />
        </Container>
        <Container padding="medium" border={{ top: true }} className={css.userButtonContainer}>
          <Button
            noStyling
            className={css.userButton}
            tooltip={
              <Layout.Vertical
                height={150}
                width={300}
                padding="medium"
                flex={{ alignItems: 'flex-start', justifyContent: 'space-between' }}
              >
                <Layout.Horizontal
                  flex={{ alignItems: 'center', justifyContent: 'flex-start' }}
                  style={{ gap: '1rem', width: '100%' }}
                >
                  <Avatar
                    size="normal"
                    hoverCard={false}
                    autoFocus={false}
                    name="Aandu"
                    email="aandu@email.com"
                    style={{ margin: 0 }}
                  />
                  <Layout.Vertical style={{ flexGrow: 1, width: 'calc(100% - 48px)' }}>
                    <Text width="100%" lineClamp={1} font={{ variation: FontVariation.H6 }}>
                      Aandu
                    </Text>
                    <Text width="100%" lineClamp={1} font={{ variation: FontVariation.SMALL }}>
                      aandu@email.com
                    </Text>
                  </Layout.Vertical>
                </Layout.Horizontal>
                <Button
                  text={
                    <Text font={{ variation: FontVariation.SMALL_BOLD }} color={Color.WHITE}>
                      {getString('logout')}
                    </Text>
                  }
                  icon="log-out"
                  width="100%"
                  intent="danger"
                  size={ButtonSize.SMALL}
                  onClick={forceLogout}
                />
              </Layout.Vertical>
            }
            tooltipProps={{
              usePortal: true,
              position: Position.RIGHT,
              interactionKind: PopoverInteractionKind.CLICK
            }}
          >
            <Layout.Horizontal flex={{ alignItems: 'center', justifyContent: 'flex-start' }}>
              <Avatar name="Aandu" email="aandu@email.com" size="normal" />
              <Text font={{ variation: FontVariation.SMALL_BOLD }}>Aandu</Text>
            </Layout.Horizontal>
          </Button>
        </Container>
      </Layout.Vertical>
    </Layout.Vertical>
  );
};

export default withErrorBoundary(SideNav, { FallbackComponent: Fallback });
