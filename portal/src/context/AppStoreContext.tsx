import React, { createContext } from 'react';

export interface AppStoreContextProps {
  readonly matchPath?: string;
  readonly renderUrl?: string;
  updateAppStore(data: Partial<Pick<AppStoreContextProps, 'matchPath' | 'renderUrl'>>): void;
}

export const AppStoreContext = createContext<AppStoreContextProps>({
  matchPath: '',
  renderUrl: '',
  updateAppStore: () => void 0
});

export function useAppStore(): AppStoreContextProps {
  return React.useContext(AppStoreContext);
}

export const AppStoreProvider: React.FC<AppStoreContextProps> = ({ children }) => {
  const [appStore, setAppStore] = React.useState<AppStoreContextProps>({
    renderUrl: '',
    matchPath: '',
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
