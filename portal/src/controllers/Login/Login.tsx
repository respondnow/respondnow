import React from 'react';
import { useHistory } from 'react-router-dom';
import LoginView from '@views/Login';
import { useLoginMutation } from '@services/server';
import { setUserDetails } from '@utils';
import { paths } from '@routes/RouteDefinitions';

const LoginController: React.FC = () => {
  const history = useHistory();

  const { mutate: loginMutation, isLoading: loginMutationLoading } = useLoginMutation(
    {},
    {
      onSuccess: data => {
        setUserDetails({
          accessToken: data.data?.token,
          isInitialLogin: data.data?.changeUserPassword
        });
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
