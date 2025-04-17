import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "@mantine/core/styles.css";
import "./index.css";

ReactDOM.createRoot(document.querySelector("#root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);

// Use contextBridge
globalThis.ipcRenderer.on("main-process-message", (_event, message) => {
  console.log(message);
});
