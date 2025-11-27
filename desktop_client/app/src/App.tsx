import { Button, MantineProvider, Tabs, Title } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { LoaderCircle, Plus, Settings, X } from "lucide-react";
import { type ReactNode, useEffect } from "react";
import {
  Outlet,
  RouterProvider,
  createHashRouter,
  useNavigate,
} from "react-router";
import { ErrorPage } from "./ErrorPage.tsx";
import { ErrorModal } from "./components/modal/errorModal";
import { H2 } from "./components/typography/H2.tsx";
import { P } from "./components/typography/P";
import { filesToOpen } from "./lib/electron/index.ts";
import {
  getLogger,
  getProjectClient,
  getSystemClient,
  initializeClients,
} from "./lib/grpcClient/client";
import { mustVoid } from "./lib/utils/errorHandlers.tsx";
import { createProjectPath } from "./pages/createProject/path.ts";
import { helpPath } from "./pages/help/path.ts";
import { indexPath, withLayoutPath } from "./pages/path.ts";
import { projectPath } from "./pages/project/path.ts";
import { settingsPath } from "./pages/settings/path.ts";
import { useGlobalErrorStore } from "./stores/globalErrorStore";
import { useProjectStore } from "./stores/projectStore.ts";
import { useSettingsStore } from "./stores/settingsStore";
import { useTabStore } from "./stores/tabStore";

function IndexRedirect(): ReactNode {
  const navigate = useNavigate();

  useEffect(() => {
    mustVoid(navigate(indexPath));
  }, [navigate]);

  return (
    <div className="flex justify-center items-center h-screen">
      <H2>Almost ready!</H2>
    </div>
  );
}

const routesCommon = {
  ErrorBoundary: ErrorPage,
  loader: Loader,
};

const router = createHashRouter([
  {
    children: [
      {
        // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
        lazy: () => import("./pages/index.tsx"),
        path: "index",
      },
      {
        // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
        lazy: () => import("./pages/settings/settings.tsx"),
        path: "settings",
      },
      {
        // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
        lazy: () => import("./pages/project/project.tsx"),
        path: "project",
      },
      {
        // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
        lazy: () => import("./pages/createProject/createProject.tsx"),
        path: "create-project",
      },
    ],
    element: <Layout />,
    path: withLayoutPath,
    ...routesCommon,
  },
  {
    children: [
      {
        // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
        lazy: () => import("./pages/help/search/search.tsx"),
        path: "search",
      },
    ],
    path: helpPath,
    ...routesCommon,
  },
  {
    children: [
      {
        element: <Loader />,
        path: "loading",
      },
    ],
    path: "/dev",
    ...routesCommon,
  },
  {
    element: <IndexRedirect />,
    path: "/",
    ...routesCommon,
  },
]);

function Loader(): ReactNode {
  return (
    <div className="flex flex-row justify-center gap-5 items-center h-screen">
      <span className="animate-spin repeat-infinite">
        <LoaderCircle />
      </span>
      <H2>Loading...</H2>
    </div>
  );
}

function Layout(): ReactNode {
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
    getSystemClient()
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
      <div className="flex flex-col gap-3 p-2 h-screen ">
        <div
          className="flex flex-row gap-3 justify-between items-center border-b border-b-gray-100 border-r-2"
          id="menu"
        >
          <div className="flex flex-row gap-3">
            <Title
              className="cursor-pointer"
              onClick={() => {
                mustVoid(navigate(indexPath));
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
                mustVoid(navigate(projectPath));
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
                            getProjectClient()
                              .closeProject(tab.handle)
                              .then(() => {
                                tabs.removeTab(tab.handle);

                                mustVoid(navigate(indexPath));
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
                mustVoid(navigate(createProjectPath));
              }}
              variant="subtle"
            >
              <Plus />
            </Button>
            <Button
              aria-label="Open settings"
              onClick={() => {
                mustVoid(navigate(settingsPath));
              }}
              variant="subtle"
            >
              <Settings />
            </Button>
          </div>
        </div>

        <div className="flex flex-col gap-3 p-2 overflow-y-auto">
          <Outlet />
        </div>
      </div>
    </>
  );
}

export default function App(): ReactNode {
  const { settings } = useSettingsStore((s) => s);
  const projects = useProjectStore((p) => p);
  const tabs = useTabStore((t) => t);

  useEffect(() => {
    initializeClients();

    // Load files on startup
    const files = filesToOpen();

    if (files.length === 0) {
      return;
    }

    getLogger().info("Loading projects on startup", {
      projects: JSON.stringify(files),
    });

    Promise.all(
      files.map(async (file) => {
        return getProjectClient()
          .openProject({
            fileName: file,
          })
          .then((resp) => {
            if (!resp.response.handle) {
              throw new Error("Invalid project returned");
            }

            projects.newProject(resp.response.handle, resp.response);
            tabs.addTab(resp.response.handle, file);
          })
          .catch((error: unknown) => {
            getLogger().error("Cannot load project", {
              error: JSON.stringify(error),
              project: file,
            });

            throw new Error(`Cannot load project ${file}`);
          });
      }),
    )
      .then(() => {
        getLogger().info("Loaded all projects successfully", {});
      })
      .catch((error: unknown) => {
        getLogger().error("Could not load all projects successfully", {
          error: JSON.stringify(error),
        });
      });
  }, [projects, tabs]);

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
