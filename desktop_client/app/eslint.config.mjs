import config from "@repodog/eslint-config";
import reactConfig from "@repodog/eslint-config-react";

export default [
  ...config,
  ...reactConfig,
  {
    languageOptions: {
      parserOptions: {
        project: "./tsconfig.json",
      },
    },
    rules: {
      "@typescript-eslint/no-non-null-assertion": "off",
      "eslint-comments/no-unlimited-disable": "warn",
      "import-x/default": "off",
      "import-x/extensions": "off",
      "import-x/no-default-export": "off",
      "prefer-arrow/prefer-arrow-functions": "off",
      "unicorn/filename-case": "off",
      "unicorn/no-abusive-eslint-disable": "warn",
    },
  },
  {
    languageOptions: {
      parserOptions: {
        project: "./tsconfig.node.json",
      },
    },
  },
];
