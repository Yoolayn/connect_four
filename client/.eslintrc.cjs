module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    "eslint:recommended",
    "plugin:react-hooks/recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:react/recommended",
    "plugin:react/jsx-runtime"
  ],
  ignorePatterns: ["dist", ".eslintrc.cjs"],
  parser: "@typescript-eslint/parser",
  plugins: ["react-refresh"],
  rules: {
    "react-refresh/only-export-components": [
      "warn",
      { allowConstantExport: true },
    ],
    "no-console": 0,
    "no-empty": "error",
    "no-empty-function": "error",
    "prefer-const": "error",
    "no-var": "error",
    "semi": ["error", "always"],
    "no-extra-semi": "error",
    "no-async-promise-executor": "error",
    "no-constant-condition": "error",
    "no-debugger": "error",
    "no-dupe-args": "error",
    "no-dupe-class-members": "error",
    "no-dupe-else-if": "error",
    "no-dupe-keys": "error",
    "no-duplicate-case": "error",
    "no-empty-pattern": "error",
    "no-fallthrough": "error",
    "no-unreachable": "error",
    "no-redeclare": "error",
    "no-regex-spaces": "error",
    "no-use-before-define": "error",
    "default-case": "error",
    "eqeqeq": "error",
    "no-alert": "error",
    "no-continue": "error",
    "arrow-spacing": ["error"],
    "block-spacing": ["error"],
    "brace-style": "error",
    "comma-spacing": "error",
    "lines-between-class-members": ["error", "always"],
    "no-multi-spaces": "error",
    "no-multiple-empty-lines": "error",
    "no-trailing-spaces": "error",
    "semi-spacing": "error",
    "semi-style": ["error", "last"],
    "init-declarations": [ "error", "always" ],
    "block-scoped-var": "error",
    "strict": [ "off", "never" ],
    "no-undef": "error",
    "no-unused-vars": "error"
  },
}
