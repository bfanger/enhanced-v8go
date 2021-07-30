const registry = {};

module.exports = function require(filepath) {
  if (registry[filepath]) {
    return registry[filepath];
  }
  registry[filepath] = {}; //@todo cyclic dependencies
  registry[filepath] = go.require(filepath);
  return registry[filepath];
};
