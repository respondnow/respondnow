import React, { createContext } from 'react';

export interface Scope {
  accountId?: string;
  orgIdentifier?: string;
  projectIdentifier?: string;
}

export interface AppStoreContextProps {
  scope: Scope;
  updateAppStore(data: Partial<AppStoreContextProps>): void;
}

export const AppStoreContext = createContext<AppStoreContextProps>({
  scope: {},
  updateAppStore: () => void 0
});

export function useAppStore(): AppStoreContextProps {
  return React.useContext(AppStoreContext);
}

export const AppStoreProvider: React.FC<AppStoreContextProps> = ({ children }) => {
  const [appStore, setAppStore] = React.useState<AppStoreContextProps>({
    scope: {},
    updateAppStore: () => void 0
  });
  const updateAppStore = React.useCallback(
    (data: Partial<AppStoreContextProps>) => {
      setAppStore(prev => ({ ...prev, ...data }));
    },
    [setAppStore]
  );

  return <AppStoreContext.Provider value={{ ...appStore, updateAppStore }}>{children}</AppStoreContext.Provider>;
};
