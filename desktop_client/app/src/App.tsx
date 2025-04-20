import { Button, MantineProvider, Title } from "@mantine/core";
import { Plus, Settings } from "lucide-react";
import { useEffect } from "react";
import { BrowserRouter, Route, Routes } from "react-router";
import { systemClient } from "./lib/grpcClient/client";
import { IndexPage } from "./pages";
import { CreateProjectPage } from "./pages/createProject/createProject";
import { ProjectPage } from "./pages/project/project";
import { SettingsPage } from "./pages/settings/settings";
import { useSettingsStore } from "./stores/settingsStore";

function App() {
  const { setSettings, settings } = useSettingsStore((s) => s);

  useEffect(() => {
    systemClient
      .getSettings({})
      .then((x) => {
        setSettings(x.response);
      })
      .catch(console.error);
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
          <div className="flex flex-row">
            <Button
              aria-label="Create project"
              onClick={() => (globalThis.location.href = "/create-project")}
              variant="subtle"
            >
              <Plus />
            </Button>
            <Button
              aria-label="Open settings"
              onClick={() => (globalThis.location.href = "/settings")}
              variant="subtle"
            >
              <Settings />
            </Button>
          </div>
        </div>

        <div className="flex flex-col gap-3 p-2">
          <BrowserRouter>
            <Routes>
              <Route element={<IndexPage />} path="/" />
              <Route element={<SettingsPage />} path="/settings" />
              <Route element={<ProjectPage />} path="/project" />
              <Route element={<CreateProjectPage />} path="/create-project" />
            </Routes>
          </BrowserRouter>
        </div>
      </div>
    </MantineProvider>
  );
}

export default App;
