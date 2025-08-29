import { Button, Input, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { Pencil } from "lucide-react";
import { useState } from "react";
import { Modal } from "@/components/modal/modal";
import { P } from "@/components/typography/P";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import CreateCharacterModal from "./createCharacterModal";

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
        {opened && <CreateCharacterModal closeDialog={close} />}
      </Modal>

      <div className="flex flex-row gap-2 pt-2 justify-between">
        <div className="flex flex-col gap-2 flex-2">
          <div className="flex flex-row gap-2 justify-between">
            <Input className="flex-grow" placeholder="TODO Search bar" />
            <Button
              onClick={() => {
                open();
              }}
            >
              Create Character
            </Button>
          </div>
          <Table variant="vertical">
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Name</Table.Th>
                <Table.Th>Factions</Table.Th>
              </Table.Tr>
              {thisProject.project.characters.map((character) => {
                const id = character.handle?.id ?? "";
                const selected = selectedCharacter == id;

                return (
                  <Table.Tr
                    className={selected ? "bg-gray-500" : ""}
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
        </div>

        <div className="flex-1 overflow-x-auto">
          TODO: {selectedCharacter} has been selected.
        </div>
      </div>
    </>
  );
}
