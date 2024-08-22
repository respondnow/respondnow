import '@testing-library/jest-dom/extend-expect';
import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { TestWrapper } from 'utils/tests';
import InfoBanner from '..';
import type { InfoBannerProps } from '../InfoBanner';

const defaultInfoProps: InfoBannerProps = {
  title: 'info_dummy',
  message: 'info_message',
  type: 'info'
};

const defaultLevelUpProps: InfoBannerProps = {
  title: 'levelup_dummy',
  message: 'levelup_message',
  type: 'levelup'
};

describe('Info Banner Tests ', () => {
  test('should show correct details for info', () => {
    render(
      <TestWrapper>
        <InfoBanner {...defaultInfoProps} />
      </TestWrapper>
    );
    expect(screen.getByText('info_message')).toBeInTheDocument();
  });

  test('should show correct details for level up', () => {
    render(
      <TestWrapper>
        <InfoBanner {...defaultLevelUpProps} />
      </TestWrapper>
    );
    expect(screen.getByText('levelup_message')).toBeInTheDocument();
  });

  test('should hide when close button is clicked', async () => {
    const { container } = render(
      <TestWrapper>
        <InfoBanner {...defaultInfoProps} />
      </TestWrapper>
    );
    expect(screen.getByText('info_message')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('cross-button'));
    await waitFor(() => {
      expect(container).toBeEmptyDOMElement();
    });
  });
});
