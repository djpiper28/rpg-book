import { Button, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { Modal } from "@/components/modal/modal";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import CreateCharacterModal from "./createCharacterModal";

export function CharacterTab() {
  const [opened, { close, open }] = useDisclosure(false);
  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);

  if (!projectHandle) {
    return "No project selected";
  }

  const thisProject = projectStore.getProject(projectHandle);

  return (
    <>
      <Modal close={close} opened={opened} title="Create Character">
        {opened && <CreateCharacterModal />}
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
          {thisProject.project.characters.map((character) => (
            <Table.Tr key={character.handle?.id}>
              <Table.Th>{character.name}</Table.Th>
              <Table.Th>TODO: Change me</Table.Th>
            </Table.Tr>
          ))}
        </Table.Thead>
      </Table>
    </>
  );
}
