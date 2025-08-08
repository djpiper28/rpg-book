import { Button, MantineProvider, Tabs, Title } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { Plus, Settings, X } from "lucide-react";
import { useEffect } from "react";
import { Route, Routes, useNavigate } from "react-router";
import { ErrorModal } from "./components/modal/errorModal";
import { P } from "./components/typography/P";
import { projectClient, systemClient } from "./lib/grpcClient/client";
import { IndexPage } from "./pages";
import { CreateProjectPage } from "./pages/createProject/createProject";
import { ProjectPage } from "./pages/project/project";
import { SettingsPage } from "./pages/settings/settings";
import { useGlobalErrorStore } from "./stores/globalErrorStore";
import { useSettingsStore } from "./stores/settingsStore";
import { useTabStore } from "./stores/tabStore";

function App() {
  const { setSettings, settings } = useSettingsStore((s) => s);
  const tabs = useTabStore((x) => x);
  const navigate = useNavigate();
  const [opened, { close, open }] = useDisclosure(false);
  const { currentError, setError } = useGlobalErrorStore((x) => x);

  useEffect(() => {
    if (currentError) {
      open();
    }
  }, [currentError, open]);

  useEffect(() => {
    systemClient
      .getSettings({})
      .then((x) => {
        setSettings(x.response);
      })
      .catch((error: unknown) => {
        setError({
          body: String(error),
        });
      });
  }, [setSettings, setError]);

  return (
    <MantineProvider
      forceColorScheme={settings.darkMode ? "dark" : "light"}
      // withCSSVariables
      // withGlobalStyles
      // withNormalizeCSS
    >
      <ErrorModal
        close={close}
        opened={opened}
        title={currentError?.title ?? "An Error Has Ocurred"}
      >
        <P>{currentError?.body}</P>
      </ErrorModal>
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
                      rightSection={
                        <X
                          className="cursor-pointer z-10"
                          color="red"
                          onClick={() => {
                            projectClient
                              .closeProject(tab.handle)
                              .then(async () => {
                                tabs.removeTab(tab.handle);

                                // TODO: figure out why this logic is not working
                                // if (tabs.selectedTab === tab.handle) {
                                await navigate("/");
                                // }
                              })
                              .catch(console.error);
                          }}
                          role="button"
                        />
                      }
                      value={tab.handle.id}
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
            <Route Component={IndexPage} path="/" />
            <Route Component={SettingsPage} path="/settings" />
            <Route Component={ProjectPage} path="/project" />
            <Route Component={CreateProjectPage} path="/create-project" />
          </Routes>
        </div>
      </div>
    </MantineProvider>
  );
}

export default App;
