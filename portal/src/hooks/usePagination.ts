import React from 'react';
import { useSearchParams } from './useSearchParams';
import { useUpdateSearchParams } from './useUpdateSearchParams';

interface UseSearchParamsBasedPaginationReturn {
  page: number;
  limit: number;
  setPage: (newPage: number) => void;
  setLimit: (newLimit: number) => void;
  resetPage: () => void;
  pageSizeOptions: number[];
}

export function usePagination(
  pageSizes?: number[],
  initial?: { page?: number; limit?: number },
  local?: boolean
): UseSearchParamsBasedPaginationReturn {
  const searchParams = useSearchParams();
  const updateSearchParams = useUpdateSearchParams();
  const [localPage, setLocalPage] = React.useState<number>(initial?.page || 0);
  const [localLimit, setLocalLimit] = React.useState<number>(initial?.limit || 15);

  const page = local ? localPage : parseInt(searchParams.get('page') ?? '0');
  const limit = local ? localLimit : parseInt(searchParams.get('limit') ?? '15');

  const pageSizeOptions = pageSizes ? [...new Set([...pageSizes, limit])].sort() : [...new Set([15, 30, limit])].sort();

  const setPage = (newPage: number): void =>
    local ? setLocalPage(newPage) : updateSearchParams({ page: newPage.toString() });
  const setLimit = (newLimit: number): void =>
    local ? setLocalLimit(newLimit) : updateSearchParams({ limit: newLimit.toString() });
  const resetPage = (): void => {
    page !== 0 && updateSearchParams({ page: '0' });
  };
  return { page, limit, setPage, setLimit, resetPage, pageSizeOptions };
}
