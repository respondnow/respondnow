import { createLocaleStrings } from '@harnessio/uicore';
import type { StringsMap } from './types';

const { useLocaleStrings: useStrings, LocaleString, StringsContextProvider } = createLocaleStrings<StringsMap>();

export { useStrings, LocaleString, StringsContextProvider };
