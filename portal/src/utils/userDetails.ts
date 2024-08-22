import { jwtDecode } from 'jwt-decode';
import { DecodedTokenType } from 'models';

export interface UserDetailsProps {
  accessToken: string;
  email: string;
  name: string;
  username: string;
  isInitialLogin: boolean;
}

export function decode<T = unknown>(arg: string): T {
  return JSON.parse(decodeURIComponent(atob(arg)));
}

export function getUserDetails(): UserDetailsProps {
  const accessToken = localStorage.getItem('accessToken') ?? '';
  const email = accessToken ? (jwtDecode(accessToken) as DecodedTokenType).email : '';
  const name = accessToken ? (jwtDecode(accessToken) as DecodedTokenType).name : '';
  const username = accessToken ? (jwtDecode(accessToken) as DecodedTokenType).username : '';
  const isInitialLogin = localStorage.getItem('isInitialLogin') === 'true';
  return { accessToken, email, isInitialLogin, name, username };
}

export function setUserDetails({
  accessToken,
  isInitialLogin
}: Partial<Omit<UserDetailsProps, 'accountID' | 'accountRole'>>): void {
  if (accessToken) localStorage.setItem('accessToken', accessToken);
  if (isInitialLogin !== undefined) localStorage.setItem('isInitialLogin', `${isInitialLogin}`);
}
