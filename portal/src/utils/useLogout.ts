import { useHistory } from 'react-router-dom';
import { paths } from '@routes/RouteDefinitions';

interface UseLogoutReturn {
  forceLogout: () => void;
}

export const useLogout = (): UseLogoutReturn => {
  const history = useHistory();

  const forceLogout = (): void => {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('projectRole');
    localStorage.removeItem('projectID');
    localStorage.removeItem('isInitialLogin');
    history.push(paths.toLogin());
  };
  return { forceLogout };
};
