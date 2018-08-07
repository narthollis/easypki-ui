import process from 'process';
import path from 'path';

import webpack from 'webpack';
import HtmlWebpackPlugin from 'html-webpack-plugin';
import TsconfigPathsPlugin from 'tsconfig-paths-webpack-plugin';
import ForkTsCheckerWebpackPlugin from 'fork-ts-checker-webpack-plugin';

// Unset the TS_NODE_PROJECT as we use this to load webpack with ts-node but we do would like to use the regular
// tsconfig.json for everything else
process.env['TS_NODE_PROJECT'] = "";

const config: webpack.Configuration = {
    mode: 'production',
    entry: {
        'main': './src/'
    },
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'easypki-ui.bundle.js'
    },
    // Add the loader for .ts files.
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: {
                    loader: 'ts-loader',
                    options: {
                        transpileOnly: true
                    }
                }
            },
        ]
    },
    // Currently we need to add '.ts' to the resolve.extensions array.
    resolve: {
        extensions: ['.ts', '.tsx', '.js', '.jsx'],
        plugins: [
            new TsconfigPathsPlugin({
                logLevel: 'INFO',
                configFile: 'tsconfig.json',
                mainFields: ['module', 'main']
            })
        ]
    },
    // Source maps support ('inline-source-map' also works)
    devtool: 'source-map',
    plugins: [
        new ForkTsCheckerWebpackPlugin({ tsconfig: './tsconfig.json' }),
        // Build our index.html so we don't have to keep updating it with new hashed chunks
        new HtmlWebpackPlugin({
            template: 'assets/templates/index.ejs',
            title: 'Test',
            chunksSortMode: 'dependency'
        }),
    ]
};

export default config;
