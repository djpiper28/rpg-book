// eslint-disable-next-line import-x/no-unresolved
import { dialog, shell } from "@electron/remote";
import { contextBridge } from "electron";
import { promises } from "node:fs";

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

async function read(url: string): Promise<Uint8Array> {
  return await promises.readFile(url);
}

contextBridge.exposeInMainWorld("electron", {
  dialog: electronDialog,
  getBuildType,
  getSystemVersion,
  read,
  shell: electronShell,
});

contextBridge.exposeInMainWorld("getEnv", getEnv);
