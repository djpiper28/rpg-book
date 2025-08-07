import { contextBridge, ipcRenderer } from "electron";

// --------- Expose some API to the Renderer process ---------
contextBridge.exposeInMainWorld("ipcRenderer", {
  invoke(...args: Parameters<typeof ipcRenderer.invoke>) {
    const [channel, ...omit] = args;
    // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
    return ipcRenderer.invoke(channel, ...omit);
  },
  off(...args: Parameters<typeof ipcRenderer.off>) {
    const [channel, ...omit] = args;
    return ipcRenderer.off(channel, ...omit);
  },
  on(...args: Parameters<typeof ipcRenderer.on>) {
    const [channel, listener] = args;
    // eslint-disable-next-line @typescript-eslint/no-shadow
    return ipcRenderer.on(channel, (event, ...args) => {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
      listener(event, ...args);
    });
  },
  send(...args: Parameters<typeof ipcRenderer.send>) {
    const [channel, ...omit] = args;
    // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
    ipcRenderer.send(channel, ...omit);
  },
});

contextBridge.exposeInMainWorld("process", {
  env: {
    RPG_BOOK_PORT: process.env.RPG_BOOK_PORT,
    RPG_BOOK_CERTIFICATE: process.env.RPG_BOOK_CERTIFICATE,
  },
});
