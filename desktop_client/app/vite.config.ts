import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import path from "node:path";
import { defineConfig } from "vite";
import commonjs from "vite-plugin-commonjs";
import electron from "vite-plugin-electron/simple";

export default defineConfig({
  appType: "spa",
  build: {
    rollupOptions: {
      treeshake: "smallest",
    },
  },
  envPrefix: ["VITE_", "RPG_BOOK_"],
  plugins: [
    react(),
    electron({
      main: {
        entry: "electron/main.ts",
      },
      preload: {
        // eslint-disable-next-line unicorn/prefer-module
        input: path.join(__dirname, "electron/preload.ts"),
      },
      renderer:
        process.env.NODE_ENV === "test"
          ? // https://github.com/electron-vite/vite-plugin-electron-renderer/issues/78#issuecomment-2053600808
            undefined
          : {},
    }),
    tailwindcss(),
    commonjs(),
  ],
  resolve: {
    alias: {
      // eslint-disable-next-line unicorn/prefer-module
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
