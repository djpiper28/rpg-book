import React, { ReactNode } from "react";
import { render } from "@testing-library/react";
import { BrowserRouter } from "react-router";
import { MantineProvider } from "@mantine/core";

export function wrappedRender(node: ReactNode) {
  return render(
    <React.StrictMode>
      <BrowserRouter>
        <MantineProvider>{node}</MantineProvider>
      </BrowserRouter>
    </React.StrictMode>,
  );
}
