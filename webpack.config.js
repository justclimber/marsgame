const path = require('path');

module.exports = {
    mode: 'development',
    entry: './js/index.js',
    output: {
        filename: 'app.js',
        path: path.resolve(__dirname, 'static/js'),
    },
    watch: true
};