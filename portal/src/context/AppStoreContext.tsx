import React, { createContext } from 'react';

export interface Scope {
  accountIdentifier?: string;
  orgIdentifier?: string;
  projectIdentifier?: string;
}

export interface AppStoreContextProps {
  scope: Scope;
  currentUserInfo: {
    name?: string;
    username?: string;
    email: string;
  };
  isInitialLogin?: boolean;
  updateAppStore: (data: Partial<AppStoreContextProps>) => void;
}

export const initialAppContext: AppStoreContextProps = {
  scope: {},
  currentUserInfo: {
    email: ''
  },
  updateAppStore: () => void 0
};

export const AppStoreContext = createContext<AppStoreContextProps>(initialAppContext);

export const AppStoreProvider: React.FC = ({ children }) => {
  const [appStore, setAppStore] = React.useState<AppStoreContextProps>(initialAppContext);

  const updateAppStore = React.useCallback(
    (data: Partial<AppStoreContextProps>) => {
      setAppStore(prev => ({ ...prev, ...data }));
    },
    [setAppStore]
  );

  return <AppStoreContext.Provider value={{ ...appStore, updateAppStore }}>{children}</AppStoreContext.Provider>;
};
