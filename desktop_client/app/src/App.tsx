import { Button, MantineProvider, Tabs, Title } from "@mantine/core";
import { X, Plus, Settings } from "lucide-react";
import { useEffect } from "react";
import { Route, Routes, useNavigate } from "react-router";
import { projectClient, systemClient } from "./lib/grpcClient/client";
import { IndexPage } from "./pages";
import { CreateProjectPage } from "./pages/createProject/createProject";
import { ProjectPage } from "./pages/project/project";
import { SettingsPage } from "./pages/settings/settings";
import { useSettingsStore } from "./stores/settingsStore";
import { useTabStore } from "./stores/tabStore";

function App() {
  const { setSettings, settings } = useSettingsStore((s) => s);
  const tabs = useTabStore((x) => x);
  const navigate = useNavigate();

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
          <div className="flex flex-row gap-3">
            <Title
              className="cursor-pointer"
              onClick={() => {
                // eslint-disable-next-line @typescript-eslint/no-floating-promises
                navigate("/");
              }}
              role="button"
            >
              RPG Book
            </Title>
            <Tabs
              onChange={(x) => {
                if (!x) {
                  return;
                }

                tabs.setSelectedTab({ id: x });

                // eslint-disable-next-line @typescript-eslint/no-floating-promises
                navigate("/project");
              }}
              value={tabs.selectedTab?.id}
              variant="outline"
            >
              <Tabs.List grow>
                {Object.values(tabs.tabs).map((tab) => {
                  return (
                    <Tabs.Tab
                      key={tab.handle.id}
                      value={tab.handle.id}
                      rightSection={
                        <X
                          color="red"
                          role="button"
                          className="cursor-pointer z-10"
                          onClick={() => {
                            projectClient
                              .closeProject(tab.handle)
                              .then(async () => {
                                tabs.removeTab(tab.handle);

                                if (tabs.selectedTab === tab.handle) {
                                  await navigate("/");
                                }
                              })
                              .catch(console.error);
                          }}
                        />
                      }
                    >
                      {tab.name}
                    </Tabs.Tab>
                  );
                })}
              </Tabs.List>
            </Tabs>
          </div>
          <div className="flex flex-row">
            <Button
              aria-label="Create project"
              onClick={() => {
                // eslint-disable-next-line @typescript-eslint/no-floating-promises
                navigate("/create-project");
              }}
              variant="subtle"
            >
              <Plus />
            </Button>
            <Button
              aria-label="Open settings"
              onClick={() => {
                // eslint-disable-next-line @typescript-eslint/no-floating-promises
                navigate("/settings");
              }}
              variant="subtle"
            >
              <Settings />
            </Button>
          </div>
        </div>

        <div className="flex flex-col gap-3 p-2">
          <Routes>
            <Route element={<IndexPage />} path="/" />
            <Route element={<SettingsPage />} path="/settings" />
            <Route element={<ProjectPage />} path="/project" />
            <Route element={<CreateProjectPage />} path="/create-project" />
          </Routes>
        </div>
      </div>
    </MantineProvider>
  );
}

export default App;
