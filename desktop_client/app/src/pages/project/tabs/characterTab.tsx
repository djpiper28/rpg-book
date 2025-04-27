import { Button, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { Modal } from "@/components/modal/modal";

export function CharacterTab() {
  const [opened, { close, open }] = useDisclosure(false);

  return (
    <>
      <Modal close={close} opened={opened} title="Create Character">
        <p>UwU</p>
      </Modal>
      <Button
        onClick={() => {
          open();
        }}
      >
        Create Character
      </Button>

      <Table variant="vertical">
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Name</Table.Th>
            <Table.Th>Factions</Table.Th>
          </Table.Tr>
        </Table.Thead>
      </Table>
    </>
  );
}
