// https://github.com/webpack-contrib/sass-loader
// https://www.typescriptlang.org/docs/handbook/react-&-webpack.html

const MiniCssExtractPlugin = require("mini-css-extract-plugin");

module.exports = {
  entry: {
    home: "./entrypoints/home.main.tsx",
    root: "./entrypoints/root.main.tsx",
    u: "./entrypoints/u.main.tsx",
  },
  output: {
    filename: "[name].js",
    path: __dirname + "/dist"
  },
  // Enable sourcemaps for debugging webpack's output.
  devtool: "source-map",
  resolve: {
    // Add '.ts' and '.tsx' as resolvable extensions.
    extensions: [".ts", ".tsx", ".js", ".json", ".scss"]
  },
  module: {
    rules: [
      // All files with a '.ts' or '.tsx' extension will be handled by 'awesome-typescript-loader'.
      { test: /\.tsx?$/, loader: "awesome-typescript-loader" },
      // All output '.js' files will have any sourcemaps re-processed by 'source-map-loader'.
      { enforce: "pre", test: /\.js$/, loader: "source-map-loader" },
      {
        test: /\.scss$/,
        use: [
          // creates style nodes from JS strings
          { loader: process.env.NODE_ENV !== "production" ? "style-loader" : MiniCssExtractPlugin.loader },
          { loader: "css-loader", options: { sourceMap: true }}, // translates CSS into CommonJS
          { loader: "sass-loader", options: { sourceMap: true }} // compiles Sass to CSS
        ]
      },
      {
        test: /\.css$/,
        use: [
          // creates style nodes from JS strings
          { loader: process.env.NODE_ENV !== "production" ? "style-loader" : MiniCssExtractPlugin.loader },
          { loader: "css-loader", options: { sourceMap: true }}, // translates CSS into CommonJS
        ]
      }
    ]
  },
  plugins: [
    new MiniCssExtractPlugin({
      // Options similar to the same options in webpackOptions.output
      // both options are optional
      filename: "[name].css",
      chunkFilename: "[id].css"
    })
  ]
};
