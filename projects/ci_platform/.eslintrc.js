module.exports = {
  root: true,
  env: {
    node: true,
  },
  plugins: ["prettier"],
  extends: [
    "plugin:vue/vue3-essential",
    "eslint:recommended",
    "prettier",
    "eslint-config-prettier",
  ],
  parserOptions: {
    parser: "@babel/eslint-parser",
  },
  rules: {
    "no-console": process.env.NODE_ENV === "production" ? "warn" : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "warn" : "off",
    //"vue/multi-word-component-names": "off",
  },
};
