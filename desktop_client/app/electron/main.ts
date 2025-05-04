import * as remote from "@electron/remote/main";
import { BrowserWindow, app, session } from "electron";
import path from "node:path";
import { fileURLToPath } from "node:url";

// const require = createRequire(import.meta.url)
// eslint-disable-next-line @typescript-eslint/naming-convention
const __dirname = path.dirname(fileURLToPath(import.meta.url));

remote.initialize();

// The built directory structure
//
// ├─┬─┬ dist
// │ │ └── index.html
// │ │
// │ ├─┬ dist-electron
// │ │ ├── main.js
// │ │ └── preload.mjs
// │
process.env.APP_ROOT = path.join(__dirname, "..");

// 🚧 Use ['ENV_NAME'] avoid vite:define plugin - Vite@2.x
export const VITE_DEV_SERVER_URL = process.env.VITE_DEV_SERVER_URL;
export const MAIN_DIST = path.join(process.env.APP_ROOT, "dist-electron");
export const RENDERER_DIST = path.join(process.env.APP_ROOT, "dist");

process.env.VITE_PUBLIC = VITE_DEV_SERVER_URL
  ? path.join(process.env.APP_ROOT, "public")
  : RENDERER_DIST;

let win: BrowserWindow | null;

function createWindow() {
  win = new BrowserWindow({
    icon: path.join(process.env.VITE_PUBLIC, "electron-vite.svg"),
    title: "RPG Book",
    webPreferences: {
      contextIsolation: false,
      nodeIntegration: true,
      preload: path.join(__dirname, "preload.mjs"),
      webSecurity: false,
    },
  });

  // Test active push message to Renderer-process.
  win.webContents.on("did-finish-load", () => {
    win.webContents.send("main-process-message", new Date().toLocaleString());
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

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
    win = undefined;
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
  // TODO: use the generated const, but things are broken and I no longer care and I want a pint
  const certificate = process.env.RPG_BOOK_CERTIFICATE ?? "";

  session.defaultSession.setCertificateVerifyProc((request, callback) => {
    try {
      const presentedCert = request.certificate.data.replace(/\s/, "");
      const trustedCert = certificate.replace(/\s/, "");
      const eq = presentedCert === trustedCert;

      if (!eq) {
        console.error("The certificates do not match");
      }

      callback(eq ? 0 : -3);
    } catch {
      callback(-2); // FAILED
    }
  });

  createWindow();
});
