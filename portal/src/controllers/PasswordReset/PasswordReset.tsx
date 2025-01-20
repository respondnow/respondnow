import React from 'react';
import { useToaster } from '@harnessio/uicore';
import { useHistory } from 'react-router-dom';
import PasswordResetView from '@views/PasswordReset';
import { useChangePasswordMutation } from '@services/server';
import { updateLocalStorage, updateUserAndScopeFromAPI } from '@utils';
import { useAppStore } from '@hooks';
import { paths } from '@routes/RouteDefinitions';

const PasswordResetController = (): React.ReactElement => {
  const { showSuccess, showError } = useToaster();
  const history = useHistory();
  const { updateAppStore } = useAppStore();

  const { mutate: updatePasswordMutation, isLoading: updatePasswordLoading } = useChangePasswordMutation(
    {},
    {
      onError: err => showError(err.message),
      onSuccess: async data => {
        showSuccess(data.message);
        const accessToken = data.data?.token || '';
        const isInitialLogin = !data.data?.lastLoginAt;

        updateLocalStorage('accessToken', accessToken);
        updateLocalStorage('isInitialLogin', String(isInitialLogin));

        await updateUserAndScopeFromAPI(updateAppStore, accessToken, isInitialLogin);

        if (isInitialLogin) {
          history.push(paths.toGetStarted());
        } else {
          history.push(paths.toIncidentDashboard());
        }
      }
    }
  );

  return (
    <PasswordResetView
      updatePasswordMutation={updatePasswordMutation}
      loading={{
        updatePassword: updatePasswordLoading
      }}
    />
  );
};

export default PasswordResetController;
