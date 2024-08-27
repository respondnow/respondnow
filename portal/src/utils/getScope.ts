import { Scope } from '@context';
import { useAppStore } from '@hooks';

export function getScope(): Required<Scope> {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const { scope } = useAppStore();
  return {
    accountIdentifier: scope.accountIdentifier || '',
    orgIdentifier: scope.orgIdentifier || '',
    projectIdentifier: scope.projectIdentifier || ''
  };
}

export function scopeExists(): boolean {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const { scope } = useAppStore();
  return !!scope.accountIdentifier || !!scope.orgIdentifier || !!scope.projectIdentifier;
}
