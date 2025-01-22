import React from 'react';
import { useHistory } from 'react-router-dom';
import LoginView from '@views/Login';
import { useLoginMutation } from '@services/server';
import { updateLocalStorage, updateUserAndScopeFromAPI } from '@utils';
import { paths } from '@routes/RouteDefinitions';
import { useAppStore } from '@hooks';

const LoginController: React.FC = () => {
  const history = useHistory();
  const { updateAppStore } = useAppStore();

  const { mutate: loginMutation, isLoading: loginMutationLoading } = useLoginMutation(
    {},
    {
      onSuccess: async loginData => {
        const accessToken = loginData.data?.token || '';
        const changePassword = loginData.data?.changeUserPassword;
        const isInitialLogin = !loginData.data?.lastLoginAt;

        updateLocalStorage('accessToken', accessToken);
        updateLocalStorage('isInitialLogin', String(isInitialLogin));

        await updateUserAndScopeFromAPI(updateAppStore, accessToken, isInitialLogin);

        if (changePassword) {
          history.push(paths.toPasswordReset());
        } else {
          history.push(paths.toIncidentDashboard());
        }
      }
    }
  );

  return <LoginView mutation={loginMutation} loading={loginMutationLoading} />;
};

export default LoginController;
