import React from 'react';
import { FontVariation, Color, PaddingProps, Spacing } from '@harnessio/design-system';
import { Container, HarnessDocTooltip, Heading, Layout, Page, PageSpinner } from '@harnessio/uicore';
import cx from 'classnames';
import NoData, { NoDataProps } from '@components/NoData';
import { ErrorWrapper } from '@errors';
import SideNav from '@components/SideNav';
import css from './DefaultLayout.module.scss';

interface DefaultLayoutProps {
  loading?: boolean;
  title?: string | React.ReactNode;
  subtitle?: React.ReactNode;
  tooltipId?: string;
  toolbar?: React.ReactNode;
  subHeader?: React.ReactNode;
  footer?: React.ReactNode;
  padding?: Spacing | PaddingProps;
  noData?: boolean;
  noDataProps?: NoDataProps;
  popovers?: React.ReactNode;
  scale?: 'full-screen' | 'full-parent';
  showSideNav?: boolean;
}

const DefaultLayout: React.FC<DefaultLayoutProps> = ({
  loading,
  title,
  subtitle,
  tooltipId,
  toolbar,
  subHeader,
  footer,
  // infoBannerProps,
  noData,
  noDataProps,
  popovers,
  padding = { top: 'medium', bottom: 'medium', left: 'xlarge', right: 'xlarge' },
  scale = 'full-screen',
  showSideNav = true,
  children
}) => {
  return (
    <Layout.Horizontal
      className={cx({
        [css.fullParentHeight]: scale === 'full-parent',
        [css.fullScreenHeight]: scale === 'full-screen'
      })}
      width="100%"
    >
      {showSideNav && <SideNav />}
      <main className={css.layout}>
        {title && (
          <ErrorWrapper>
            <Page.Header
              className={css.header}
              size={subtitle ? 'large' : 'standard'}
              toolbar={toolbar}
              title={
                typeof title === 'string' || typeof title === 'undefined' ? (
                  <Container>
                    <Heading level={4} font={{ variation: FontVariation.H4 }} color={Color.GREY_700}>
                      {title} <HarnessDocTooltip tooltipId={tooltipId} useStandAlone />
                    </Heading>
                    {subtitle}
                  </Container>
                ) : (
                  <>
                    {title}
                    {subtitle && subtitle}
                  </>
                )
              }
            />
          </ErrorWrapper>
        )}
        {subHeader && <Page.SubHeader className={css.subHeader}>{subHeader}</Page.SubHeader>}
        <div className={cx(css.content, { [css.contentLoading]: loading, [css.contentNoData]: noData })}>
          {loading && <PageSpinner />}
          {!noData ? (
            <Container height="100%" width="100%" padding={padding}>
              {children}
            </Container>
          ) : (
            <Container height="100%" width="100%" flex={{ align: 'center-center' }}>
              <NoData {...noDataProps} />
            </Container>
          )}
        </div>
        <div className={cx(css.footer, { [css.noFooterPadding]: padding === 'none' })}>{footer}</div>
        {popovers}
      </main>
    </Layout.Horizontal>
  );
};

export default DefaultLayout;
