const webpack = require('webpack');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HappyPack = require( 'happypack' )

var lessLoader = ExtractTextPlugin.extract(
  "style!css?sourceMap"
);

const config = {
  context: __dirname,
  entry: './index.js',
  output: {
    path: __dirname,
    filename: 'bundle.js'
  },
  module: {
    loaders: [{
      exclude: /node_modules/,
      test: /\.(js|jsx)$/,
      loader: 'happypack/loader?id=babelJs'
    },
    {
      test: /\.css$/,
      loader: ExtractTextPlugin.extract('style-loader', 'css!less?indentedSyntax=true&sourceMap=true')
    },
    {
      test: /\.less$/,
      loader: ExtractTextPlugin.extract('style-loader', 'css!less?indentedSyntax=true&sourceMap=true')
    },
    { test: /\.(eot|woff|woff2|ttf|svg|png|jpe?g|gif)(\?\S*)?$/
    , loader: 'url?limit=100000&name=[name].[ext]'
    },
    {
      test: /\.scss$/,
      loader: ExtractTextPlugin.extract('style-loader', 'css!less?indentedSyntax=true&sourceMap=true')
    }]
  },
  devServer: {
    historyApiFallback: true,
    port: 8000,
    contentBase: './'
  },
  resolve: {
    extensions: ['', '.js', '.jsx'],
  },
  plugins: [
    new webpack.DefinePlugin({ 'process.env':{ 'NODE_ENV': JSON.stringify('production') } }),
    new webpack.optimize.DedupePlugin(),
    new webpack.optimize.OccurenceOrderPlugin(false),
    new HappyPack({
      id: 'babelJs',
      loaders: [ 'babel' ]
    }),
    new webpack.optimize.UglifyJsPlugin({
      compress: { warnings: false },
      output: {comments: false },
      mangle: false,
      sourcemap: false,
      minimize: true,
      mangle: { except: ['$super', '$', 'exports', 'require', '$q', '$ocLazyLoad'] }
    }),
    new ExtractTextPlugin('public/stylesheets/app.css', {
      allChunks: true
    })
  ]
};

module.exports = config;