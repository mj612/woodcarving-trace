module.exports = {
  chainWebpack: config => {
    // 禁用CSP检查
    config.plugin('html').tap(args => {
      args[0].meta = {
        'Content-Security-Policy': "script-src 'self' 'unsafe-eval'"
      };
      return args;
    });
  },
  
  devServer: {
    headers: {
      'Content-Security-Policy': "script-src 'self' 'unsafe-eval'"
    }
  }
}