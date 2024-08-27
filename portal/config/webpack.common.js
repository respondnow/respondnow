const path = require('path');
const prettier = require('prettier');

const { DefinePlugin } = require('webpack');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const TsconfigPathsPlugin = require('tsconfig-paths-webpack-plugin');
const { RetryChunkLoadPlugin } = require('webpack-retry-chunk-load-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const { GenerateStringTypesPlugin } = require('@harnessio/uicore/tools/GenerateStringTypesPlugin');

const CONTEXT = process.cwd();

module.exports = {
  target: 'web',
  context: CONTEXT,
  stats: {
    modules: false,
    children: false
  },
  output: {
    publicPath: '/',
    path: path.resolve(CONTEXT, 'dist/'),
    pathinfo: false
  },
  module: {
    rules: [
      {
        test: /\.m?js$/,
        include: /node_modules/,
        type: 'javascript/auto'
      },
      {
        test: /\.(j|t)sx?$/,
        exclude: /node_modules/,
        use: [
          {
            loader: 'ts-loader',
            options: {
              transpileOnly: true
            }
          }
        ]
      },
      {
        test: /\.module\.scss$/,
        exclude: /node_modules/,
        use: [
          MiniCssExtractPlugin.loader,
          {
            loader: '@teamsupercell/typings-for-css-modules-loader',
            options: {
              prettierConfigFile: path.resolve(__dirname, '../.prettierrc.yml'),
              disableLocalsExport: true
            }
          },
          {
            loader: 'css-loader',
            options: {
              importLoaders: 1,
              modules: {
                mode: 'local',
                localIdentName: 'respond_now_[name]_[local]_[hash:base64:6]',
                exportLocalsConvention: 'camelCaseOnly'
              }
            }
          },
          {
            loader: 'sass-loader',
            options: {
              sassOptions: {
                includePaths: [path.join(CONTEXT, 'src')]
              },
              sourceMap: false,
              implementation: require('sass')
            }
          }
        ]
      },
      {
        test: /(?<!\.module)\.scss$/,
        exclude: /node_modules/,
        use: [
          MiniCssExtractPlugin.loader,
          {
            loader: 'css-loader',
            options: {
              importLoaders: 1,
              modules: false
            }
          },
          {
            loader: 'sass-loader',
            options: {
              sassOptions: {
                includePaths: [path.join(CONTEXT, 'src')]
              },
              implementation: require('sass')
            }
          }
        ]
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
        include: /node_modules/
      },
      {
        test: /\.ya?ml$/,
        type: 'json',
        use: [
          {
            loader: 'yaml-loader'
          }
        ]
      },
      {
        test: /\.(jpg|jpeg|png|svg|gif)$/,
        type: 'asset'
      }
    ]
  },
  resolve: {
    extensions: ['.mjs', '.js', '.ts', '.tsx', '.json'],
    plugins: [new TsconfigPathsPlugin()]
  },
  plugins: [
    new DefinePlugin({
      'process.env': '{}' // required for @blueprintjs/core
    }),
    new RetryChunkLoadPlugin({
      retryDelay: 1000,
      maxRetries: 5
    }),
    new CopyWebpackPlugin({
      patterns: [{ from: 'src/static' }]
    }),
    new GenerateStringTypesPlugin({
      input: 'src/strings/strings.en.yaml',
      output: 'src/strings/types.ts',
      partialType: true,
      preProcess: async content => {
        const prettierConfig = await prettier.resolveConfig(process.cwd());
        return prettier.format(content, {
          ...prettierConfig,
          parser: 'typescript'
        });
      }
    })
  ]
};
