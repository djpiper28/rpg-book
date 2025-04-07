import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { ThemeProvider } from "./components/theme-provider.tsx";
import { EnvVarCertificate, EnvVarPort } from "./lib/launcherTypes/index.ts";
import { credentials } from "@grpc/grpc-js";
import { Empty } from "./lib/grpcClient/pb/common_pb";
import { SystemSvcClient } from "./lib/grpcClient/pb/system_grpc_pb";
import { LogLevel, LogRequest } from "./lib/grpcClient/pb/system_pb";

function mustHave(x: string | undefined): string {
  if (x) {
    return x;
  }

  throw new Error("Cannot find env vars");
}

const certificate = mustHave(process.env[EnvVarCertificate]);
const port = mustHave(process.env[EnvVarPort]);

Promise.resolve(async () => {
  const grpcClient = new SystemSvcClient(
    `localhost:${port}`,
    credentials.createSsl(Buffer.from(certificate)),
  );

  grpcClient.getSettings(new Empty(), (settings) => {
    console.log("Settings are:", settings);
  });

  const req = new LogRequest();
  req.setLevel(LogLevel.INFO)
  req.setCaller(__filename)
  req.setMessage("Hello world from JS")
  grpcClient.log(req, () => {})
});

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider defaultTheme="light" storageKey="vite-ui-theme">
      <App />
    </ThemeProvider>
  </React.StrictMode>,
);

// Use contextBridge
window.ipcRenderer.on("main-process-message", (_event, message) => {
  console.log(message);
});
