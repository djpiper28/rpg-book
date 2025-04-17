// eslint.config.mjs
import config from "@repodog/eslint-config";
import reactConfig from "@repodog/eslint-config-react";

// eslint convention is to export default
// eslint-disable-next-line import-x/no-default-export
export default [
  ...config,
  ...reactConfig,
  {
    ignores: ["./tsconfig.json", "./tsconfig.node.json", "./vite.config.ts"],
    rules: {
      "import-x/default": "off",
      "import-x/extensions": "off",
      "import-x/no-default-export": "off",
      "prefer-arrow/prefer-arrow-functions": "off",
      "unicorn/filename-case": "off",
    },
  },
];
