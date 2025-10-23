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

function filesToOpen(): string[] {
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  return JSON.parse(process.env.FILES_TO_OPEN) as string[];
}

contextBridge.exposeInMainWorld("electron", {
  dialog: electronDialog,
  filesToOpen,
  getBuildType,
  getSystemVersion,
  shell: electronShell,
});

contextBridge.exposeInMainWorld("getEnv", getEnv);
