import { useAppStore } from '@context';

interface Scope {
  accountId: string;
  orgIdentifier: string;
  projectIdentifier: string;
}

export function getScope(): Scope {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const { scope } = useAppStore();
  return {
    accountId: scope.accountId || '',
    orgIdentifier: scope.orgIdentifier || '',
    projectIdentifier: scope.projectIdentifier || ''
  };
}
