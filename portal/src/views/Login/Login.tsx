import React from 'react';
import { Button, ButtonVariation, Card, Container, FormInput, Layout, Text, useToaster } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import { Form, Formik } from 'formik';
import * as Yup from 'yup';
import { UseMutateFunction } from '@tanstack/react-query';
import { withErrorBoundary } from 'react-error-boundary';
import mainLogo from '@images/respondNow.svg';
import { useStrings } from '@strings';
import PasswordInput from '@components/PasswordInput';
import { UtilsDefaultResponseDto, LoginMutationProps, UserLoginResponseDto } from '@services/server';
import { Fallback } from '@errors';
import css from './Login.module.scss';

interface LoginViewProps {
  mutation: UseMutateFunction<UserLoginResponseDto, UtilsDefaultResponseDto, LoginMutationProps<never>, unknown>;
  loading: boolean;
}

interface LoginFormProps {
  username: string;
  password: string;
}

const LoginView: React.FC<LoginViewProps> = props => {
  const { mutation, loading } = props;
  const { getString } = useStrings();
  const { showError } = useToaster();

  return (
    <Container height="100%" background={Color.PRIMARY_BG} flex={{ align: 'center-center' }}>
      <Card className={css.loginCardContainer}>
        <Layout.Vertical flex={{ alignItems: 'center', justifyContent: 'space-between' }} height="100%">
          <Layout.Vertical width="100%" spacing="small">
            <img src={mainLogo} alt="RespondNow" height={40} />
            <Text font={{ variation: FontVariation.H6, align: 'center' }} color={Color.GREY_400}>
              {getString('loginSubHeading')}
            </Text>
          </Layout.Vertical>
          <Formik<LoginFormProps>
            initialValues={{ username: '', password: '' }}
            validationSchema={Yup.object().shape({
              username: Yup.string().email(getString('emailInvalid')).required(getString('emailRequired')),
              password: Yup.string().required(getString('passwordRequired'))
            })}
            onSubmit={values => {
              mutation(
                {
                  queryParams: {},
                  body: {
                    email: values.username,
                    password: values.password
                  }
                },
                {
                  onError: error => showError(error.message)
                }
              );
            }}
          >
            <Form className={css.formContainer}>
              <FormInput.Text
                name="username"
                inputGroup={{
                  type: 'email'
                }}
                placeholder={`${getString('enter')} ${getString('email')}`}
                label={<Text font={{ variation: FontVariation.FORM_LABEL }}>{getString('email')}</Text>}
              />
              <PasswordInput
                name="password"
                label={<Text font={{ variation: FontVariation.FORM_LABEL }}>{getString('password')}</Text>}
                placeholder={`${getString('enter')} ${getString('password')}`}
                disabled={false}
              />
              <Button
                type="submit"
                text={getString('continue')}
                variation={ButtonVariation.PRIMARY}
                width="100%"
                disabled={loading}
                loading={loading}
              />
            </Form>
          </Formik>
          <Layout.Horizontal flex={{ align: 'center-center' }} spacing="xsmall">
            <Text font={{ variation: FontVariation.SMALL_SEMI }}>{getString('loginFooterText')}</Text>
            <Text font={{ variation: FontVariation.SMALL_SEMI }} color={Color.PRIMARY_7}>
              {getString('termsAndConditions')}
            </Text>
          </Layout.Horizontal>
        </Layout.Vertical>
      </Card>
    </Container>
  );
};

export default withErrorBoundary(LoginView, { FallbackComponent: Fallback });
