import { ReactNode } from "react";
import { render } from "@testing-library/react";
import { BrowserRouter } from "react-router";
import { MantineProvider } from "@mantine/core";

export function wrappedRender(node: ReactNode) {
  return render(
    <BrowserRouter>
      <MantineProvider
        forceColorScheme="dark"
        // withCSSVariables
        // withGlobalStyles
        // withNormalizeCSS
      >
        {node}
      </MantineProvider>
    </BrowserRouter>,
  );
}
