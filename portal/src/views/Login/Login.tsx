import React from 'react';
import { Button, ButtonVariation, Card, Container, FormInput, Layout, Text } from '@harnessio/uicore';
import { Color, FontVariation } from '@harnessio/design-system';
import { Form, Formik } from 'formik';
import * as Yup from 'yup';
import mainLogo from '@images/respondNow.svg';
import { useStrings } from '@strings';
import PasswordInput from '@components/PasswordInput';
import css from './Login.module.scss';

interface LoginFormProps {
  username: string | undefined;
  password: string | undefined;
}

const LoginView: React.FC = () => {
  const { getString } = useStrings();

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
            initialValues={{ username: undefined, password: undefined }}
            validationSchema={Yup.object().shape({
              username: Yup.string().email().required(getString('emailRequired')),
              password: Yup.string().required(getString('passwordRequired'))
            })}
            onSubmit={() => void 0}
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
              <Button type="submit" text={getString('continue')} variation={ButtonVariation.PRIMARY} width="100%" />
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

export default LoginView;
