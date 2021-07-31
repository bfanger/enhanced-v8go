// https://nodejs.org/dist/latest-v16.x/docs/api/modules.html

const stack = [];

function resolve(id) {
  const origin = stack[stack.length - 1];
  const [filepath, err] = go.requireResolve(id, origin);
  if (err) {
    throw err;
  }
  return filepath;
}

/**
 *
 * @param {string} module name or path
 * @returns {any} exported module content
 */
function require(id) {
  if (!require.main) {
    throw new Error("can't use require outside of a module");
  }
  if (!require.cache) {
    require.cache = { [require.main.id]: require.main.exports };
  }
  stack.push(require.main.filename);
  const filepath = resolve(id);
  if (require.cache[filepath]) {
    return require.cache[filepath].exports;
  }

  const cache = { exports: {} };
  require.cache[filepath] = cache;
  stack.push(filepath);
  const [exports, err] = go.requireFile(filepath);
  stack.pop();
  if (err) {
    throw err;
  }
  Object.assign(cache.exports, exports); // cyclic dependencies
  return exports;
}
require.resolve = resolve;
module.exports = require;
