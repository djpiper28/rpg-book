import { Settings } from "lucide-react";
import { client, logger } from "./lib/grpcClient/client";
import { Button, MantineProvider, Title } from "@mantine/core";
import { useEffect } from "react";
import { useSettingsStore } from "./stores/settingsStore";
import { BrowserRouter, Routes, Route } from "react-router";
import { IndexPage } from "./pages";
import { SettingsPage } from "./pages/settings/settings";

function App() {
  const { settings, setSettings } = useSettingsStore((s) => s);
  useEffect(() => {
    client
      .getSettings({})
      .then((x) => {
        setSettings(x.response);
      })
      .catch((e) => {
        logger.error(`Cannot get settings ${e}`, {});
      });
  }, [setSettings]);

  return (
    <MantineProvider
      withCSSVariables
      withGlobalStyles
      withNormalizeCSS
      forceColorScheme={settings.darkMode ? "dark" : "light"}
    >
      <div className="flex flex-col gap-3 p-2">
        <div
          id="menu"
          className="flex flex-row gap-3 justify-between items-center border-b border-b-gray-100 border-r-2"
        >
          <Title>
            <a href="/">RPG Book</a>
          </Title>
          <Button
            aria-label="Open settings"
            variant="filled"
            onClick={() => (window.location.href = "/settings")}
          >
            <Settings /> Settings
          </Button>
        </div>

        <div className="flex flex-col gap-3 p-2">
          <BrowserRouter>
            <Routes>
              <Route path="/" element={<IndexPage />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Routes>
          </BrowserRouter>
        </div>
      </div>
    </MantineProvider>
  );
}

export default App;
