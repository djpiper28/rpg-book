import { defineConfig } from "vite";
import viteConfig from "./vite.config.ts";

export default defineConfig({
  ...viteConfig,
  test: {
    environment: "jsdom",
    globals: true,
    include: ["**/*.test.tsx", "**/*.test.ts"],
    setupFiles: "./setup-vitest.ts",
  },
});
