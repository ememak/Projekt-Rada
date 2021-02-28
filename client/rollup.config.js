const {nodeResolve} = require('@rollup/plugin-node-resolve');
const commonjs = require('@rollup/plugin-commonjs');

/*class NormalizePaths {
  resolveId(importee, importer) {
    if (importee.startsWith('Projekt_Rada')) {
      return `${importee.replace('Projekt_Rada', '')}.js`;
    }
  }
}*/

module.exports = {
  plugins: [
    //new NormalizePaths(),
    nodeResolve({
      mainFields: ['browser', 'es2015', 'module', 'jsnext:main', 'main'],
    }),
    commonjs(),
  ],
};
