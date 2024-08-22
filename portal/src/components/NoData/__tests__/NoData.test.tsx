import '@testing-library/jest-dom/extend-expect';
import React from 'react';
import { render, screen } from '@testing-library/react';
import { TestWrapper } from 'utils/tests';
import NoData from '..';

describe('NoData Card Test ', () => {
  test('should show correct title and image', () => {
    render(
      <TestWrapper>
        <NoData title={'test_title'} />
      </TestWrapper>
    );

    expect(screen.getByTestId('no-data-img')).toBeInTheDocument();
    expect(screen.getByText('test_title')).toBeInTheDocument();
  });
});
