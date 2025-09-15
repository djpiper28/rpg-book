import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import path from "node:path";
import { defineConfig } from "vite";
import commonjs from "vite-plugin-commonjs";
import electron from "vite-plugin-electron";

const nodeBuiltins = [
  "assert",
  "buffer",
  "child_process",
  "crypto",
  "events",
  "fs",
  "os",
  "path",
  "stream",
  "url",
  "util",
  "zlib",
];

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
    electron([
      {
        entry: "electron/main.ts",
        vite: {
          build: {
            rollupOptions: {
              external: ["electron", ...nodeBuiltins],
            },
          },
        },
      },
      {
        entry: "electron/preload.ts",
        onstart(options): void {
          options.reload();
        },
        vite: {
          build: {
            rollupOptions: {
              external: ["electron", ...nodeBuiltins],
            },
          },
        },
      },
    ]),
    tailwindcss(),
    commonjs(),
  ],
  resolve: {
    alias: {
      // eslint-disable-next-line unicorn/prefer-module
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    strictPort: false, // Allow vite to pick a random port
  },
});
