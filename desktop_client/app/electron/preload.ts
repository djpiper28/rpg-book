// eslint-disable-next-line import-x/no-unresolved
import { contextBridge } from "electron";
import {
  electronDialog,
  electronShell,
  getBuildType,
  getEnv,
  getSystemVersion,
  read,
} from "../src/lib/electron/main";

contextBridge.exposeInMainWorld("electron", {
  dialog: electronDialog,
  getBuildType,
  getSystemVersion,
  read,
  shell: electronShell,
});

contextBridge.exposeInMainWorld("getEnv", getEnv);
