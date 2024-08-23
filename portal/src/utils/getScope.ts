import { useAppStore } from '@hooks';

export interface Scope {
  accountIdentifier: string;
  orgIdentifier: string;
  projectIdentifier: string;
}

export function getScope(): Scope {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const { scope } = useAppStore();
  return {
    accountIdentifier: scope.accountIdentifier || '',
    orgIdentifier: scope.orgIdentifier || '',
    projectIdentifier: scope.projectIdentifier || ''
  };
}
