import { Table } from "@mantine/core";
import { H2 } from "@/components/typography/H2";
import { P } from "@/components/typography/P";

export function IndexPage() {
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
        <Table.Tbody>
          <Table.Tr>
            <Table.Th>Test Project</Table.Th>
            <Table.Th>3 days ago</Table.Th>
            <Table.Th>5 MiB</Table.Th>
            <Table.Th>/home/example/test.rpg</Table.Th>
          </Table.Tr>
        </Table.Tbody>
      </Table>
    </>
  );
}
