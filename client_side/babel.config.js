module.exports = {
  presets: [
    '@vue/cli-plugin-babel/preset',
  ],
  plugins: [
    ['prismjs', {
      'languages': ['cpp', 'c', 'java', 'python', 'javascript', 'css', 'html', 'bash', 'json', 'markdown', 'php', 'sql', 'typescript', 'yaml'],
      'plugins': ['line-numbers'],
      'theme': 'tomorrow',
      'css': true
    }]
  ]
}
