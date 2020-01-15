const path = require('path');

module.exports = {
    mode: 'development',
    entry: './js/index.js',
    devtool: 'inline-source-map',
    devServer: {
        contentBase: path.resolve(__dirname, 'static'),
        watchContentBase: true,
    },
    output: {
        filename: 'js/app.js',
        path: path.resolve(__dirname, 'static'),
    },
    watch: true,
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: ['style-loader', 'css-loader'],
            },
        ],
    },
};