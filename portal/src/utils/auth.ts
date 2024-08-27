import jwtDecode from 'jwt-decode';
interface UserDetails {
  role: string;
  uid: string;
  username: string;
  name?: string;
  email?: string;
  exp: Date;
  iat: Date;
}

// Checks if the user is  authenticated
export function isUserAuthenticated(): boolean {
  const token = localStorage.getItem('accessToken');
  if (token) {
    const userDetails = jwtDecode(token) as unknown as UserDetails;
    const expiry = (new Date(userDetails.exp).getTime() as number) * 1000;
    if (new Date(expiry) > new Date()) return true;
  }
  return false;
}

export function getTokenFromLocalStorage(): {
  token: string;
  isInitialLogin: boolean;
} {
  return {
    token: localStorage.getItem('accessToken') || '',
    isInitialLogin: localStorage.getItem('isInitialLogin') === 'true'
  };
}
