module.exports = {
    entry: "./static/app.js",
    output: {
        path: __dirname,
        filename: "./static/bundle.js"
    },
    module: {
        loaders: [
            { 
              test: /\.css$/, 
              loader: "style!css" 
            },
            { 
              test: /\.vue$/, 
              loader: 'vue-loader',
              options: {
              // vue-loader options
              } 
            }
        ]
    },
    resolve: {
      alias: {
        'vue$': 'vue/dist/vue.common.js'
      }
    },
};