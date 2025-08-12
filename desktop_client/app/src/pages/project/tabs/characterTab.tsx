import { Button, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { Modal } from "@/components/modal/modal";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import CreateCharacterModal from "./createCharacterModal";
import { useState } from "react";
import { P } from "@/components/typography/P";
import { Pencil } from "lucide-react";

export function CharacterTab() {
  const [selectedCharacter, setSelectedCharacter] = useState("");
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
          {thisProject.project.characters.map((character) => {
            const id = character.handle?.id ?? "";
            return (
              <Table.Tr
                className={selectedCharacter == id ? "bg-gray-600" : ""}
                key={id}
                onClick={() => {
                  setSelectedCharacter(id);
                }}
              >
                <Table.Th>
                  <div className="flex flex-row gap-2">
                    <button
                      className="cursor-pointer"
                      onClick={() => {
                        // TODO: open view / update dialog
                      }}
                    >
                      <Pencil className="p-1" />
                    </button>
                    <P className="text-wrap">{character.name}</P>
                  </div>
                </Table.Th>
                <Table.Th>TODO: Change me</Table.Th>
              </Table.Tr>
            );
          })}
        </Table.Thead>
      </Table>
    </>
  );
}
