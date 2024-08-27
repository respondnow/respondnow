const { defaults: tsjPreset } = require('ts-jest/presets');
const { pathsToModuleNameMapper } = require('ts-jest');
const { compilerOptions } = require('./tsconfig');

module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  transform: {
    ...tsjPreset.transform,
    '^.+\\.svg$': '<rootDir>/scripts/jest/svgTransform.js',
    '^.+\\.jsx?$': [
      'ts-jest',
      {
        tsconfig: '<rootDir>/tsconfig.json',
        isolatedModules: true
      }
    ]
  },
  setupFiles: ['fake-indexeddb/auto'],
  testEnvironmentOptions: {
    beforeParse(window) {
      window.document.childNodes.length === 0;
      window.alert = msg => {
        console.log(msg);
      };
      window.matchMedia = () => ({});
      window.scrollTo = () => {};
    }
  },

  transformIgnorePatterns: [
    '<rootDir>/node_modules/(?!(date-fns|lodash-es|p-debounce|@harnessio/react-monitoring-service-client)/)'
  ],
  testPathIgnorePatterns: ['<rootDir>/dist'],
  moduleNameMapper: {
    '^lodash-es$': 'lodash',
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
      '<rootDir>/test/jest/__mocks__/fileMock.js',
    '\\.yaml$': 'yaml-jest',
    '\\.s?css$': 'identity-obj-proxy',
    '@harnessio/monaco-yaml.*': '<rootDir>/scripts/jest/file-mock.js',
    ...pathsToModuleNameMapper(compilerOptions.paths)
  },
  testMatch: ['**/?(*.)+(test).[jt]s?(x)'],
  modulePaths: ['<rootDir>'],
  moduleDirectories: ['node_modules', 'src'],
  coveragePathIgnorePatterns: ['/node_modules/', 'dist/', '/config/', 'nginx/'],
  coverageReporters: ['clover', 'json', 'lcov', 'text'],
  collectCoverageFrom: [
    '<rootDir>/src/**/*.{tsx,jsx,ts}',
    '!<rootDir>/src/_mocks/**',
    '!<rootDir>/src/**/*mock*.{ts,tsx}',
    '!<rootDir>/src/**/*Mock*.{ts,tsx}',
    '!<rootDir>/src/**/*.{d.ts,module.scss.d.ts}',
    '!<rootDir>/src/**/index.ts',
    '!<rootDir>/src/App/App.tsx',
    '!<rootDir>/src/context/**',
    '!<rootDir>/src/controllers/**',
    '!<rootDir>/src/interfaces/**',
    '!<rootDir>/src/routes/**',
    '!<rootDir>/src/services/**',
    '!<rootDir>/src/utils/tests.tsx'
  ]
};
