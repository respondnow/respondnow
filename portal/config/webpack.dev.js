const path = require('path');
const util = require('util');
const fs = require('fs');

require('dotenv').config();

const { createCertificate } = require('pem');
const { merge } = require('webpack-merge');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const HTMLWebpackPlugin = require('html-webpack-plugin');
const { DefinePlugin, WatchIgnorePlugin } = require('webpack');
const commonConfig = require('./webpack.common');

const createCertificateAsync = util.promisify(createCertificate);

const baseUrl = process.env.BASE_URL;
const targetLocalHost = (process.env.TARGET_LOCALHOST && JSON.parse(process.env.TARGET_LOCALHOST)) || true;

const devConfig = {
  mode: 'development',
  entry: './src/index.ts',
  devtool: 'cheap-module-source-map',
  cache: { type: 'filesystem' },
  output: {
    filename: '[name].js',
    chunkFilename: '[name].[id].js'
  },
  devServer: {
    historyApiFallback: true,
    port: 8191,
    proxy: {
      '/api': {
        pathRewrite: targetLocalHost ? { '^/api': '' } : {},
        target: targetLocalHost ? process.env.API_SERVER_URI || 'http://localhost:8080' : baseUrl,
        changeOrigin: true
      }
    }
  },
  plugins: [
    new MiniCssExtractPlugin({
      filename: '[name].css',
      chunkFilename: '[name].[id].css'
    }),
    new HTMLWebpackPlugin({
      template: 'src/index.html',
      filename: 'index.html',
      minify: false
    }),
    new DefinePlugin({
      'process.env': '{}', // required for @blueprintjs/core
      __DEV__: true
    }),
    new WatchIgnorePlugin({
      paths: [/node_modules(?!\/@harnessio)/, /\.(d|test)\.tsx?$/, /stringTypes\.ts/, /\.snap$/]
    })
  ]
};

module.exports = async () => {
  const pemConfig = await fs.promises.readFile(path.resolve(process.cwd(), './config/pem.cfg'));
  const keys = await createCertificateAsync({
    days: 365,
    selfSigned: true,
    commonName: 'localhost',
    country: 'US',
    state: 'California',
    locality: 'San Francisco',
    organization: 'Respond Now',
    altNames: ['localhost'],
    config: pemConfig
  });

  devConfig.devServer.server = {
    type: 'https',
    options: {
      key: keys.serviceKey,
      cert: keys.certificate
    }
  };

  return merge(commonConfig, devConfig);
};
