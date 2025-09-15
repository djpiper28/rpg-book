import React, { type ReactNode } from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "@mantine/core/styles.css";
import "./index.css";

function Root(): ReactNode {
  return (
    <React.StrictMode>
      <App />
    </React.StrictMode>
  );
}

ReactDOM.createRoot(document.querySelector("#root")!).render(<Root />);
