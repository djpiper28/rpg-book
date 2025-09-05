import { SiGithub } from "@icons-pack/react-simple-icons";
import { Button, Table } from "@mantine/core";
import dayjs from "dayjs";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { H2 } from "@/components/typography/H2";
import { Link } from "@/components/typography/Link";
import { P } from "@/components/typography/P";
import { DbExtension } from "@/lib/databaseTypes";
import { electronDialog } from "@/lib/electron";
import { bytesToFriendly } from "@/lib/fileSizeHelpers";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type RecentProjectsResp } from "@/lib/grpcClient/pb/project";
import { mustVoid } from "@/lib/utils/errorHandlers";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useSettingsStore } from "@/stores/settingsStore";
import { useTabStore } from "@/stores/tabStore";
import { projectPath } from "./project/path";

export function Component() {
  const [recentProjects, setRecentProjects] = useState<RecentProjectsResp>({
    projects: [],
  });

  const tabs = useTabStore((x) => x);
  const projects = useProjectStore((x) => x);
  const navigate = useNavigate();
  const { setError } = useGlobalErrorStore((x) => x);
  const settings = useSettingsStore((x) => x.settings);

  const openProject = (filename: string, projectNameOverride?: string) => {
    getProjectClient()
      .openProject({ fileName: filename })
      .then((resp) => {
        if (!resp.response.handle) {
          setError({
            body: "Invalid project returned by server - no handle",
          });

          return;
        }

        tabs.addTab(
          {
            id: resp.response.handle.id,
          },
          projectNameOverride ?? filename,
        );

        projects.newProject(resp.response.handle, resp.response);
        mustVoid(navigate(projectPath));
      })
      .catch((error: unknown) => {
        setError({
          body: String(error),
          title: "Cannot open project",
        });
      });
  };

  useEffect(() => {
    getProjectClient()
      .recentProjects({})
      .then((x) => {
        setRecentProjects(x.response);
      })
      .catch((error: unknown) => {
        setError({
          body: String(error),
        });
      });
  }, [setError, projects.projects]);

  return (
    <>
      <H2>Recent Projects:</H2>
      <P>
        Open an recent project or{" "}
        <Button
          onClick={() => {
            electronDialog
              .showOpenDialog({
                buttonLabel: "Open project",
                filters: [
                  {
                    extensions: [DbExtension.replace(".", "")],
                    name: "Project (*.sqlite)",
                  },
                ],
                properties: ["openFile", "multiSelections"],
                title: "Chose a project to open",
              })
              .then((result: Electron.OpenDialogReturnValue) => {
                if (result.canceled) {
                  return;
                }

                for (const file of result.filePaths) {
                  openProject(file);
                }
              })
              .catch(console.error);
          }}
        >
          Browser Your Files
        </Button>
      </P>
      <Table variant="vertical">
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Project Name</Table.Th>
            <Table.Th>Last Opened</Table.Th>
            <Table.Th>Size</Table.Th>
            <Table.Th>Location</Table.Th>
          </Table.Tr>
        </Table.Thead>
        {recentProjects.projects.map((x) => {
          const time = dayjs(x.lastOpened);
          const size = bytesToFriendly(Number(x.fileSizeBytes));

          return (
            <Table.Tbody
              className="cursor-pointer"
              key={x.fileName}
              onClick={() => {
                openProject(x.fileName, x.projectName);
              }}
              role="button"
            >
              <Table.Tr>
                <Table.Th>{x.projectName}</Table.Th>
                <Table.Th>{time.format()}</Table.Th>
                <Table.Th>{size}</Table.Th>
                <Table.Th>{x.fileName}</Table.Th>
              </Table.Tr>
            </Table.Tbody>
          );
        })}
      </Table>
      <div className="flex flex-col gap-5 py-10">
        <div className="flex flex-row gap-2">
          <SiGithub />
          <P>
            View this project on{" "}
            <Link
              href="https://github.com/djpiper28/rpg-book"
              openInBrowser={true}
              safe={true}
            >
              Github
            </Link>
            . Made by Danny Piper (djpiper28).
          </P>
        </div>
        {settings.devMode ? (
          <>
            <Button
              color="red"
              onClick={() => {
                getLogger().info("Causing a handled error on purpose", {});

                setError({
                  body: "Test Error Message.",
                  title: "Dev Mode Test Error",
                });
              }}
            >
              Force handled error
            </Button>
            <Button
              color="red"
              onClick={() => {
                getLogger().info("Causing an unhandled error on purpose", {});
                throw new Error("Test");
              }}
            >
              Force unhandled error
            </Button>
            <Button
              color="red"
              onClick={() => {
                getLogger().info(
                  "Going to an error boundary page on purpose",
                  {},
                );

                mustVoid(navigate("/dev/not-found"));
              }}
            >
              Error Boundary Page
            </Button>
            <Button
              onClick={() => {
                mustVoid(navigate("/dev/loading"));
              }}
            >
              Loading Page
            </Button>
          </>
        ) : (
          ""
        )}
      </div>
    </>
  );
}
