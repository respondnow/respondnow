import jwtDecode from 'jwt-decode';
import { pick } from 'lodash-es';
import { DecodedTokenType } from 'models';
import { AppStoreContextProps, Scope } from '@context';
import { getUserMapping } from '@services/server';

export function setUserDetails(
  updateAppStore: AppStoreContextProps['updateAppStore'],
  accessToken: string,
  isInitialLogin?: boolean,
  scope?: Scope
): void {
  const email = accessToken ? (jwtDecode(accessToken) as DecodedTokenType).email : '';
  const name = accessToken ? (jwtDecode(accessToken) as DecodedTokenType).name : '';
  const username = accessToken ? (jwtDecode(accessToken) as DecodedTokenType).username : '';

  updateAppStore({
    currentUserInfo: { email, name, username },
    scope,
    isInitialLogin
  });
}

export function getUsername(accessToken: string): string {
  return accessToken ? (jwtDecode(accessToken) as DecodedTokenType).username : '';
}

export function updateLocalStorage(key: string, value: string): void {
  localStorage.setItem(key, value);
}

export async function updateUserAndScopeFromAPI(
  updateAppStore: AppStoreContextProps['updateAppStore'],
  accessToken: string,
  isInitialLogin?: boolean
): Promise<void> {
  await getUserMapping({
    queryParams: {
      userId: getUsername(accessToken)
    }
  }).then(mappingData => {
    const scope: Scope = pick(mappingData?.data?.defaultMapping, [
      'accountIdentifier',
      'orgIdentifier',
      'projectIdentifier'
    ]);

    setUserDetails(updateAppStore, accessToken, isInitialLogin, scope);
  });
}
