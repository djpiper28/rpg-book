import { Table } from "@mantine/core";
import dayjs from "dayjs";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { H2 } from "@/components/typography/H2";
import { P } from "@/components/typography/P";
import { bytesToFriendly } from "@/lib/fileSizeHelpers";
import { projectClient } from "@/lib/grpcClient/client";
import { type RecentProjectsResp } from "@/lib/grpcClient/pb/project";
import { useTabStore } from "@/stores/tabStore";

export function IndexPage() {
  const [recentProjects, setRecentProjects] = useState<RecentProjectsResp>({
    projects: [],
  });

  const tabs = useTabStore((x) => x);
  const navigate = useNavigate();

  useEffect(() => {
    projectClient
      .recentProjects({})
      .then((x) => {
        setRecentProjects(x.response);
      })
      .catch(console.error);
  }, []);

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
                    tabs.addTab(
                      {
                        id: resp.response.id,
                      },
                      x.projectName,
                    );

                    await navigate("/project");
                  })
                  .catch(console.error);
              }}
              role="button"
            >
              <Table.Tr>
                <Table.Th>{x.projectName}</Table.Th>
                <Table.Th>{time.format()}</Table.Th>
                <Table.Th>{size.toString()}</Table.Th>
                <Table.Th>{x.fileName}</Table.Th>
              </Table.Tr>
            </Table.Tbody>
          );
        })}
      </Table>
    </>
  );
}
