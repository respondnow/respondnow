import { defineConfig } from '@harnessio/oats-cli';
import reactQueryPlugin from '@harnessio/oats-plugin-react-query';
import { mapKeys, omit } from 'lodash-es';
import dotenv from 'dotenv';
dotenv.config();

const baseUrl = process.env.BASE_URL || 'https://api.respondnow.io';

function normalizeAPIPath(url: string): string {
  return url.replace(/\/{2,}/g, '/');
}

export default defineConfig({
  services: {
    auth: {
      url: normalizeAPIPath(`${baseUrl}/auth/swagger/doc.json`),
      output: 'src/services/server',
      transformer(spec) {
        return {
          ...spec,
          paths: mapKeys(spec.paths, (_val, key) => normalizeAPIPath(`/auth/${key}`))
        };
      },
      genOnlyUsed: true,
      plugins: [
        reactQueryPlugin({
          customFetcher: '@services/fetcher',
          allowedOperationIds: ['ChangePassword', 'Login', 'SignUp', 'CreateIncident', 'ListIncidents'],
          overrides: {}
        })
      ]
    }
  }
});
