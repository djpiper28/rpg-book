import { shell, dialog } from "@electron/remote";

export const electronShell = shell;
export const electronDialog = dialog;

export function getSystemVersion() {
  return process.getSystemVersion();
}

export const env: {[key: string]: string} = {
  RPG_BOOK_PORT: process.env.RPG_BOOK_PORT,
};

export function getBuildType() {
  return process.env.VITE_DEV_SERVER_URL ? "Development" : "Production";
}
