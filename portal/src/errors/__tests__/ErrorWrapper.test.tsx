import '@testing-library/jest-dom/extend-expect';
import React from 'react';
import { render, screen } from '@testing-library/react';
import { withErrorBoundary } from 'react-error-boundary';
import { TestWrapper } from 'utils/tests';
import { EmptyFallback, ErrorWrapper, Fallback } from '..';

describe('ErrorWrapper Test ', () => {
  beforeEach(() => {
    jest.spyOn(console, 'error').mockImplementation(jest.fn());
  });

  test('should show wrapped children', () => {
    render(
      <TestWrapper>
        <ErrorWrapper>
          <>error</>
        </ErrorWrapper>
      </TestWrapper>
    );

    expect(screen.getByText('error')).toBeInTheDocument();
  });

  test('should show empty DOM when no children', () => {
    const { container } = render(
      <TestWrapper>
        <ErrorWrapper>{undefined}</ErrorWrapper>
      </TestWrapper>
    );

    expect(container).toBeEmptyDOMElement();
  });

  test('test Fallback component', () => {
    const ErroredElement = (): React.ReactElement => {
      throw new Error('error occurred');
    };
    const Element = withErrorBoundary(ErroredElement, { FallbackComponent: Fallback });
    render(
      <TestWrapper>
        <ErrorWrapper>
          <Element />
        </ErrorWrapper>
      </TestWrapper>
    );

    expect(screen.getByText('Something went wrong:')).toBeInTheDocument();
    expect(screen.getByText('error occurred')).toBeInTheDocument();
  });

  test('test Empty Fallback component', () => {
    const ErroredElement = (): React.ReactElement => {
      throw new Error();
    };
    const Element = withErrorBoundary(ErroredElement, { FallbackComponent: EmptyFallback });
    const { container } = render(
      <TestWrapper>
        <ErrorWrapper>
          <Element />
        </ErrorWrapper>
      </TestWrapper>
    );

    expect(container).toBeEmptyDOMElement();
  });
});
