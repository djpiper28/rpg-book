/// <reference types="vite-plugin-electron/electron-env" />

declare namespace NodeJS {
  interface ProcessEnv {
    /**
     * The built directory structure
     *
     * ```tree
     * ├─┬─┬ dist
     * │ │ └── index.html
     * │ │
     * │ ├─┬ dist-electron
     * │ │ ├── main.js
     * │
     * ```
     */
    APP_ROOT: string;
    RPG_BOOK_CERTIFICATE: string;
    RPG_BOOK_PORT: string;
    /** /dist/ or /public/ */
    VITE_PUBLIC: string;
  }
}
