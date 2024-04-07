import { defineConfig } from 'vite';
import path from 'path';
import inject from '@rollup/plugin-inject';

export default defineConfig({
  base: '',
  root: path.resolve(__dirname, 'src'),

  resolve: {
    alias: {
      '~coreui': path.resolve(__dirname, 'node_modules/@coreui/coreui-pro'),
    },
  },
  build: {
    minify: true,
    manifest: true,
    rollupOptions: {
      input: {
        scans: './src/app/scans.js',
        main: './src/app/main.js',
        app: './src/scss/app.scss',
      },
    },
    outDir: '../static',
  },
  plugins: [
    inject({
      include: '**/*.js', // Only include JavaScript files
      exclude: '**/*.scss', // Exclude SCSS files
      $: 'jquery',
      jQuery: 'jquery',
    }),
  ],
});
