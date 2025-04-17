import { Button, MantineProvider, Title } from "@mantine/core";
import { Settings } from "lucide-react";
import { useEffect } from "react";
import { BrowserRouter, Route, Routes } from "react-router";
import { client, logger } from "./lib/grpcClient/client";
import { IndexPage } from "./pages";
import { SettingsPage } from "./pages/settings/settings";
import { useSettingsStore } from "./stores/settingsStore";

function App() {
  const { setSettings, settings } = useSettingsStore((s) => s);

  useEffect(() => {
    client
      .getSettings({})
      .then((x) => {
        setSettings(x.response);
      })
      .catch((error: unknown) => {
        logger.error(`Cannot get settings ${JSON.stringify(error)}`, {});
      });
  }, [setSettings]);

  return (
    <MantineProvider
      forceColorScheme={settings.darkMode ? "dark" : "light"}
      withCSSVariables
      withGlobalStyles
      withNormalizeCSS
    >
      <div className="flex flex-col gap-3 p-2">
        <div
          className="flex flex-row gap-3 justify-between items-center border-b border-b-gray-100 border-r-2"
          id="menu"
        >
          <Title>
            <a href="/">RPG Book</a>
          </Title>
          <Button
            aria-label="Open settings"
            onClick={() => (globalThis.location.href = "/settings")}
            variant="filled"
          >
            <Settings /> Settings
          </Button>
        </div>

        <div className="flex flex-col gap-3 p-2">
          <BrowserRouter>
            <Routes>
              <Route element={<IndexPage />} path="/" />
              <Route element={<SettingsPage />} path="/settings" />
            </Routes>
          </BrowserRouter>
        </div>
      </div>
    </MantineProvider>
  );
}

export default App;
