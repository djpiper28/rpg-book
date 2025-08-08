import { Button, Table } from "@mantine/core";
import dayjs from "dayjs";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { H2 } from "@/components/typography/H2";
import { Link } from "@/components/typography/Link";
import { P } from "@/components/typography/P";
import { bytesToFriendly } from "@/lib/fileSizeHelpers";
import { projectClient } from "@/lib/grpcClient/client";
import { type RecentProjectsResp } from "@/lib/grpcClient/pb/project";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useSettingsStore } from "@/stores/settingsStore";
import { useTabStore } from "@/stores/tabStore";
import { SiGithub } from "@icons-pack/react-simple-icons";

export function IndexPage() {
  const [recentProjects, setRecentProjects] = useState<RecentProjectsResp>({
    projects: [],
  });

  const tabs = useTabStore((x) => x);
  const projects = useProjectStore((x) => x);
  const navigate = useNavigate();
  const { setError } = useGlobalErrorStore((x) => x);
  const settings = useSettingsStore((x) => x.settings);

  useEffect(() => {
    projectClient
      .recentProjects({})
      .then((x) => {
        setRecentProjects(x.response);
      })
      .catch((error: unknown) => {
        setError({
          body: String(error),
        });
      });
  }, [setError]);

  return (
    <>
      <H2>Recent Projects:</H2>
      <P>Lorem ipsum sit amet dolor</P>
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
                projectClient
                  .openProject({ fileName: x.fileName })
                  .then(async (resp) => {
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
                      x.projectName,
                    );

                    projects.newProject(resp.response.handle, resp.response);
                    await navigate("/project");
                  })
                  .catch((error: unknown) => {
                    setError({
                      body: String(error),
                      title: "Cannot open project",
                    });
                  });
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
                setError({
                  body: "Test Error Message.",
                  title: "Dev Mode Test Error",
                });
              }}
            >
              Force Error
            </Button>
          </>
        ) : (
          ""
        )}
      </div>
    </>
  );
}
