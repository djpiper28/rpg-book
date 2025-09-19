import { dialog, shell } from "@electron/remote";
import { contextBridge } from "electron";

const electronShell = shell;
const electronDialog = dialog;

function getSystemVersion(): string {
  return process.getSystemVersion();
}

function getEnv(): Record<string, string> {
  return {
    RPG_BOOK_PORT: process.env.RPG_BOOK_PORT,
  };
}

function getBuildType(): string {
  return process.env.VITE_DEV_SERVER_URL ? "Development" : "Production";
}

contextBridge.exposeInMainWorld("electron", {
  dialog: electronDialog,
  getBuildType,
  getSystemVersion,
  shell: electronShell,
});

contextBridge.exposeInMainWorld("getEnv", getEnv);
