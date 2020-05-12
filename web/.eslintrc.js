// https://github.com/typescript-eslint/typescript-eslint/blob/master/docs/getting-started/linting/README.md

module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  plugins: [
    '@typescript-eslint',
    'import',
  ],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/eslint-recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react/recommended',
  ],

  // several rules copied from https://github.com/typescript-eslint/typescript-eslint/blob/master/.eslintrc.js
  rules: {
    // disables
    "@typescript-eslint/no-use-before-define": "off",
    "@typescript-eslint/explicit-function-return-type": "off",
    "@typescript-eslint/no-non-null-assertion": "off",
    "@typescript-eslint/no-explicit-any": "off",
    "react/no-unescaped-entities": "off", // e.g. `'` can be escaped with `&apos;`, `&lsquo;`, `&#39;`, `&rsquo;`
    "react/prop-types": "off", // e.g. 'username' is missing in props validation

    // modifications
    "@typescript-eslint/member-delimiter-style": [
      "error",
      {
        "multiline": {
            "delimiter": "none",
            "requireLast": false
        },
        "singleline": {
            "delimiter": "comma",
            "requireLast": false
        },
      },
    ],

    // Forbid the use of extraneous packages
    "import/no-extraneous-dependencies": [
      "error",
      {
        devDependencies: true,
        peerDependencies: true,
        optionalDependencies: false,
      },
    ],
  },

  parserOptions: {
    sourceType: 'module',
    ecmaFeatures: {
      jsx: true,
    },
    project: ['./tsconfig.json'],
    tsconfigRootDir: __dirname,
  },

  settings: {
    // https://github.com/yannickcr/eslint-plugin-react#configuration
    react: {
      version: "detect",
    },
  },
};
