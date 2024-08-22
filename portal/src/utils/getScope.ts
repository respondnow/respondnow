import { Scope, useAppStore } from '@context';

export function getScope(): Scope {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const { scope } = useAppStore();
  return {
    ...scope
  };
}
