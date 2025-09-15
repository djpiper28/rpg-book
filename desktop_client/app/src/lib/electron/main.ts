import { shell, dialog } from "@electron/remote";
import { promises } from "fs";

export const electronShell = shell;
export const electronDialog = dialog;

export function getSystemVersion() {
  return process.getSystemVersion();
}

export function getEnv(): { [key: string]: string } {
  return {
    RPG_BOOK_PORT: process.env.RPG_BOOK_PORT,
  };
}

export function getBuildType() {
  return process.env.VITE_DEV_SERVER_URL ? "Development" : "Production";
}

export async function read(url: string): Promise<Uint8Array> {
  return await promises.readFile(url);
}
