import { Button, MantineProvider, Tabs, Title } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { Plus, Settings, X } from "lucide-react";
import { useEffect } from "react";
import {
  Outlet,
  RouterProvider,
  createBrowserRouter,
  useNavigate,
} from "react-router";
import { ErrorModal } from "./components/modal/errorModal";
import { P } from "./components/typography/P";
import { projectClient, systemClient } from "./lib/grpcClient/client";
import { IndexPage } from "./pages";
import { createProjectPath } from "./pages/createProject/path.ts";
import { indexPath, withLayoutPath } from "./pages/path.ts";
import { projectPath } from "./pages/project/path.ts";
import { settingsPath } from "./pages/settings/path.ts";
import { useGlobalErrorStore } from "./stores/globalErrorStore";
import { useSettingsStore } from "./stores/settingsStore";
import { useTabStore } from "./stores/tabStore";

function IndexRedirect() {
  const navigate = useNavigate();
  // eslint-disable-next-line @typescript-eslint/no-floating-promises
  navigate(indexPath);

  return "Almost ready!";
}

const router = createBrowserRouter([
  {
    children: [
      {
        element: <IndexPage />,
        path: "index",
      },
      {
        lazy: () => import("./pages/settings/settings.tsx"),
        path: "settings",
      },
      {
        lazy: () => import("./pages/project/project.tsx"),
        path: "project",
      },
      {
        lazy: () => import("./pages/createProject/createProject.tsx"),
        path: "create-project",
      },
    ],
    element: <Layout />,
    path: withLayoutPath,
  },
  {
    element: <IndexRedirect />,
    path: "/",
  },
]);

function Layout() {
  const { setSettings } = useSettingsStore((s) => s);
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
    <>
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
                navigate(indexPath);
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
                navigate(projectPath);
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
                                await navigate(indexPath);
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
                navigate(createProjectPath);
              }}
              variant="subtle"
            >
              <Plus />
            </Button>
            <Button
              aria-label="Open settings"
              onClick={() => {
                // eslint-disable-next-line @typescript-eslint/no-floating-promises
                navigate(settingsPath);
              }}
              variant="subtle"
            >
              <Settings />
            </Button>
          </div>
        </div>

        <div className="flex flex-col gap-3 p-2">
          <Outlet />
        </div>
      </div>
    </>
  );
}

function App() {
  const { settings } = useSettingsStore((s) => s);

  return (
    <MantineProvider
      forceColorScheme={settings.darkMode ? "dark" : "light"}
      // withCSSVariables
      // withGlobalStyles
      // withNormalizeCSS
    >
      <RouterProvider router={router} />
    </MantineProvider>
  );
}

export default App;
