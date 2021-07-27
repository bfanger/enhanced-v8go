function log(...args) {
  // @todo format arguments
  println(...args);
}

export default new Proxy(
  { log },
  {
    get: function (target, prop) {
      if (target[prop]) {
        return target[prop];
      }

      throw new Error("console." + prop + " is not implemented");
    },
  }
);
