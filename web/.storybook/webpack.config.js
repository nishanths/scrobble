const path = require("path");

const srcPath = path.join(__dirname, '../');

module.exports = ({ config }) => {
  config.module.rules.push({
    test: /\.(ts|tsx)$/,
    include: [srcPath],
      use: [
        {
          loader: require.resolve('awesome-typescript-loader'),
          options: {
            configFileName: './.storybook/tsconfig.json',
          }
        },
      ]
  });
  config.resolve.extensions.push('.ts', '.tsx');
  return config;
};