import '@testing-library/jest-dom/extend-expect';
import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { Button } from '@harness/uicore';
import { TestWrapper } from 'utils/tests';
import DefaultLayout from '..';

describe('Default Layout Test ', () => {
  test('should show title, breadcrumbs and body', async () => {
    render(
      <TestWrapper>
        <DefaultLayout
          breadcrumbs={[
            {
              label: 'breadcrumb',
              url: '/'
            }
          ]}
          title="title text"
        >
          page body
        </DefaultLayout>
      </TestWrapper>
    );

    expect(screen.getByText('title text')).toBeInTheDocument();
    expect(screen.getByText('breadcrumb')).toBeInTheDocument();
    expect(screen.getByText('page body')).toBeInTheDocument();
  });

  test('should show toolbar and subheader props', async () => {
    const handleClick = jest.fn();
    render(
      <TestWrapper>
        <DefaultLayout
          breadcrumbs={[
            {
              label: 'breadcrumb',
              url: '/'
            }
          ]}
          title="title text"
          toolbar={<Button text="toolbar-button" onClick={handleClick} />}
          subHeader={<>subheader</>}
        >
          page body
        </DefaultLayout>
      </TestWrapper>
    );

    expect(screen.getByText('subheader')).toBeInTheDocument();
    fireEvent.click(screen.getByText('toolbar-button'));
    await waitFor(() => {
      expect(handleClick).toHaveBeenCalled();
    });
  });

  test('should show info banner', async () => {
    const handleClick = jest.fn();
    render(
      <TestWrapper>
        <DefaultLayout
          breadcrumbs={[
            {
              label: 'breadcrumb',
              url: '/'
            }
          ]}
          title="title text"
          toolbar={<Button text="toolbar-button" onClick={handleClick} />}
          subHeader={<>subheader</>}
          infoBannerProps={{
            message: 'info banner text',
            type: 'info'
          }}
        >
          page body
        </DefaultLayout>
      </TestWrapper>
    );

    expect(screen.getByText('info banner text')).toBeInTheDocument();
  });
});
