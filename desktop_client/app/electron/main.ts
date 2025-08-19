import * as remote from "@electron/remote/main";
// eslint-disable-next-line import-x/no-unresolved
import { BrowserWindow, app, session } from "electron";
import path from "node:path";
import { fileURLToPath } from "node:url";

// eslint-disable-next-line @typescript-eslint/naming-convention
const __dirname = path.dirname(fileURLToPath(import.meta.url));
remote.initialize();
process.env.APP_ROOT = path.join(__dirname, "..");

export const VITE_DEV_SERVER_URL = process.env.VITE_DEV_SERVER_URL;
export const MAIN_DIST = path.join(process.env.APP_ROOT, "dist-electron");
export const RENDERER_DIST = path.join(process.env.APP_ROOT, "dist");

process.env.VITE_PUBLIC = VITE_DEV_SERVER_URL
  ? path.join(process.env.APP_ROOT, "public")
  : RENDERER_DIST;

const isDevServer = !!VITE_DEV_SERVER_URL;
let win: BrowserWindow | null;

function createWindow() {
  win = new BrowserWindow({
    autoHideMenuBar: true,
    icon: path.join(process.env.VITE_PUBLIC, "electron-vite.svg"),
    minHeight: 400,
    minWidth: 600,
    title: "RPG Book",
    webPreferences: {
      contextIsolation: true,
      devTools: isDevServer,
      nodeIntegration: false,
      preload: path.join(__dirname, "preload.mjs"),
      webSecurity: false,
    },
  });

  // Test active push message to Renderer-process.
  win.webContents.on("did-finish-load", () => {
    win?.webContents.send("main-process-message", new Date().toLocaleString());
  });

  remote.enable(win.webContents);

  if (VITE_DEV_SERVER_URL) {
    win.loadURL(VITE_DEV_SERVER_URL).then().catch(console.error);
  } else {
    win
      .loadFile(path.join(RENDERER_DIST, "index.html"))
      .then()
      .catch(console.error);
  }
}

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
    // eslint-disable-next-line unicorn/no-null
    win = null;
  }
});

app.on("activate", () => {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

// eslint-disable-next-line @typescript-eslint/no-floating-promises,unicorn/prefer-top-level-await
app.whenReady().then(() => {
  const certificate = process.env.RPG_BOOK_CERTIFICATE;

  session.defaultSession.setCertificateVerifyProc((request, callback) => {
    try {
      const presentedCert = request.certificate.data.replace(/\s/, "");
      const trustedCert = certificate.replace(/\s/, "");
      const eq = presentedCert === trustedCert;

      if (!eq) {
        console.error("The certificates do not match");
      }

      callback(eq ? 0 : -2);
    } catch {
      callback(-2);
    }
  });

  session.defaultSession.setUserAgent("RPG-Book");

  createWindow();
});
