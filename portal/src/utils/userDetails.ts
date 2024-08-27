import jwtDecode from 'jwt-decode';
import { DecodedTokenType } from 'models';
import { AppStoreContextProps } from '@context';

export function setUserDetails(
  updateAppStore: AppStoreContextProps['updateAppStore'],
  accessToken: string,
  isInitialLogin?: boolean
): void {
  const email = accessToken ? (jwtDecode(accessToken) as DecodedTokenType).email : '';
  const name = accessToken ? (jwtDecode(accessToken) as DecodedTokenType).name : '';
  const username = accessToken ? (jwtDecode(accessToken) as DecodedTokenType).username : '';

  updateAppStore({
    currentUserInfo: { email, name, username },
    isInitialLogin
  });
}

export function updateLocalStorage(key: string, value: string): void {
  localStorage.setItem(key, value);
}

export function setScope(
  updateAppStore: AppStoreContextProps['updateAppStore'],
  accountIdentifier: string,
  orgIdentifier: string,
  projectIdentifier: string
): void {
  updateAppStore({
    scope: {
      accountIdentifier,
      orgIdentifier,
      projectIdentifier
    }
  });
}
