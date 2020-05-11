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
  config.module.rules.push({
    test: /\.scss$/,
    use: [
      // creates style nodes from JS strings
      { loader: "style-loader" },
      { loader: "css-loader", options: { sourceMap: true }}, // translates CSS into CommonJS
      { loader: "sass-loader", options: { sourceMap: true }} // compiles Sass to CSS
    ]
  });
  config.module.rules.push({
    test: /\.css$/,
    use: [
      // creates style nodes from JS strings
      { loader: "style-loader" },
      { loader: "css-loader", options: { sourceMap: true }}, // translates CSS into CommonJS
    ]
  });
  config.resolve.extensions.push('.ts', '.tsx');
  return config;
};
