import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router";
import App from "./App.tsx";
import "@mantine/core/styles.css";
import "./index.css";

ReactDOM.createRoot(document.querySelector("#root")).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
);

// Use contextBridge
// eslint-disable-next-line
globalThis.ipcRenderer.on("main-process-message", (_event, message) => {
  console.log(message);
});
