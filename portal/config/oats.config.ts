import { defineConfig } from '@harnessio/oats-cli';
import reactQueryPlugin from '@harnessio/oats-plugin-react-query';
import { mapKeys } from 'lodash-es';
import dotenv from 'dotenv';
dotenv.config();

const baseUrl = process.env.BASE_URL || 'http://localhost:8080';

function normalizeAPIPath(url: string): string {
  return url.replace(/\/{2,}/g, '/');
}

export default defineConfig({
  services: {
    api: {
      url: normalizeAPIPath(`${baseUrl}/swagger/doc.json`),
      output: 'src/services/server',
      transformer(spec) {
        return {
          ...spec,
          paths: mapKeys(spec.paths, (_val, key) => normalizeAPIPath(`/api/${key}`))
        };
      },
      genOnlyUsed: true,
      plugins: [
        reactQueryPlugin({
          customFetcher: '@services/fetcher',
          allowedOperationIds: [
            'changePassword',
            'login',
            'createIncident',
            'listIncidents',
            'getIncident',
            'getUserMappings'
          ],
          overrides: {}
        })
      ]
    }
  }
});
