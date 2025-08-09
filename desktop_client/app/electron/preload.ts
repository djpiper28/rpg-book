import { contextBridge } from "electron";
import {
  electronDialog,
  electronShell,
  env,
  getSystemVersion,
} from "../src/lib/electron/main";

contextBridge.exposeInMainWorld("electron", {
  dialog: electronDialog,
  getSystemVersion,
  shell: electronShell,
});

contextBridge.exposeInMainWorld("env", env);
