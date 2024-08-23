import React from 'react';
import { useHistory } from 'react-router-dom';
import LoginView from '@views/Login';
import { useLoginMutation } from '@services/server';
import { setUserDetails, updateLocalStorage } from '@utils';
import { paths } from '@routes/RouteDefinitions';
import { useAppStore } from '@hooks';

const LoginController: React.FC = () => {
  const history = useHistory();
  const { updateAppStore } = useAppStore();

  const { mutate: loginMutation, isLoading: loginMutationLoading } = useLoginMutation(
    {},
    {
      onSuccess: data => {
        const accessToken = data.data?.token?.split(' ')[1] || '';
        updateLocalStorage('accessToken', accessToken);
        updateLocalStorage('isInitialLogin', data.data?.changeUserPassword ? 'true' : 'false');
        setUserDetails(updateAppStore, accessToken, data.data?.changeUserPassword);
        if (data.data?.changeUserPassword) {
          history.push(paths.toPasswordReset());
        } else {
          history.push(paths.toGetStarted());
        }
      }
    }
  );

  return <LoginView mutation={loginMutation} loading={loginMutationLoading} />;
};

export default LoginController;
