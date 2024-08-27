import { Color, FontVariation } from '@harnessio/design-system';
import { Button, ButtonVariation, Card, Container, Layout, Text } from '@harnessio/uicore';
import React from 'react';
import { Form, Formik } from 'formik';
import * as Yup from 'yup';
import { UseMutateFunction } from '@tanstack/react-query';
import { withErrorBoundary } from 'react-error-boundary';
import { Icon } from '@harnessio/icons';
import mainLogo from '@images/respondNow.svg';
import { useStrings } from '@strings';
import PasswordInput from '@components/PasswordInput';
import { ChangePasswordMutationProps, UserChangePasswordResponseDto, UtilsDefaultResponseDto } from '@services/server';
import { useAppStore } from '@hooks';
import { Fallback } from '@errors';
import css from './PasswordReset.module.scss';

interface PasswordResetViewProps {
  updatePasswordMutation: UseMutateFunction<
    UserChangePasswordResponseDto,
    UtilsDefaultResponseDto,
    ChangePasswordMutationProps<never>,
    unknown
  >;
  loading: {
    updatePassword: boolean;
  };
}

interface AccountPasswordChangeFormProps {
  oldPassword: string;
  newPassword: string;
  reEnterNewPassword: string;
}

const PasswordResetView = (props: PasswordResetViewProps): React.ReactElement => {
  const { updatePasswordMutation, loading } = props;
  const { getString } = useStrings();
  const { currentUserInfo } = useAppStore();

  function handleSubmit(values: AccountPasswordChangeFormProps): void {
    updatePasswordMutation({
      queryParams: {},
      body: {
        email: currentUserInfo.email ?? '',
        password: values.oldPassword,
        newPassword: values.newPassword
      }
    });
  }

  return (
    <Container height="100%" background={Color.PRIMARY_BG} flex={{ align: 'center-center' }}>
      <Card className={css.passwordResetContainer}>
        <Layout.Vertical flex={{ alignItems: 'center', justifyContent: 'space-between' }} height="100%">
          <Layout.Vertical width="100%" spacing="small">
            <img src={mainLogo} alt="RespondNow" height={40} />
          </Layout.Vertical>
          <Formik<AccountPasswordChangeFormProps>
            initialValues={{
              oldPassword: '',
              newPassword: '',
              reEnterNewPassword: ''
            }}
            onSubmit={handleSubmit}
            validationSchema={Yup.object().shape({
              oldPassword: Yup.string().required(getString('enterOldPassword')),
              newPassword: Yup.string()
                .required(getString('enterNewPassword'))
                .min(8, getString('fieldMinLength', { length: 8 }))
                .max(16, getString('fieldMaxLength', { length: 16 })),
              reEnterNewPassword: Yup.string()
                .required(getString('reEnterNewPassword'))
                .oneOf([Yup.ref('newPassword'), null], getString('passwordsDoNotMatch'))
            })}
          >
            {formikProps => {
              return (
                <Form className={css.formContainer}>
                  <Layout.Vertical width="100%">
                    <PasswordInput
                      name="oldPassword"
                      placeholder={getString('oldPassword')}
                      label={<Text font={{ variation: FontVariation.FORM_LABEL }}>{getString('oldPassword')}</Text>}
                    />
                    <PasswordInput
                      name="newPassword"
                      placeholder={getString('newPassword')}
                      label={<Text font={{ variation: FontVariation.FORM_LABEL }}>{getString('newPassword')}</Text>}
                    />
                    <PasswordInput
                      name="reEnterNewPassword"
                      placeholder={getString('reEnterNewPassword')}
                      label={
                        <Text font={{ variation: FontVariation.FORM_LABEL }}>{getString('reEnterNewPassword')}</Text>
                      }
                    />
                  </Layout.Vertical>
                  <Layout.Horizontal width="100%" flex={{ alignItems: 'center', justifyContent: 'center' }}>
                    <Button
                      type="submit"
                      variation={ButtonVariation.PRIMARY}
                      text={loading.updatePassword ? <Icon name="loading" size={16} /> : getString('confirm')}
                      disabled={loading.updatePassword || Object.keys(formikProps.errors).length > 0}
                      style={{ minWidth: '90px' }}
                      width="100%"
                    />
                  </Layout.Horizontal>
                </Form>
              );
            }}
          </Formik>
        </Layout.Vertical>
      </Card>
    </Container>
  );
};

export default withErrorBoundary(PasswordResetView, { FallbackComponent: Fallback });
