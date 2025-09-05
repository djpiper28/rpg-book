import { Button, Input, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useEffect, useState } from "react";
import { Modal } from "@/components/modal/modal";
import MarkdownRenderer from "@/components/renderers/markdown";
import { H2 } from "@/components/typography/H2";
import { P } from "@/components/typography/P";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type BasicCharacterDetails } from "@/lib/grpcClient/pb/project_character";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import CreateCharacterModal from "./createCharacterModal";

export function CharacterTab() {
  const [selectedCharacterId, setSelectedCharacterId] = useState("");

  const [selectedCharcter, setSelectedCharacter] = useState<
    BasicCharacterDetails | undefined
  >();

  const [opened, { close, open }] = useDisclosure(false);
  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);

  useEffect(() => {
    const sync = async () => {
      if (!selectedCharacterId) {
        return;
      }

      if (!projectHandle) {
        return;
      }

      try {
        const resp = await getProjectClient().getCharacter({
          character: {
            id: selectedCharacterId,
          },
          project: projectHandle,
        });

        setSelectedCharacter(resp.response);
      } catch (error: unknown) {
        getLogger().error("Cannot get character", {
          character: selectedCharacterId,
          error: JSON.stringify(error),
          project: JSON.stringify(projectHandle),
        });
      }
    };

    // eslint-disable-next-line @typescript-eslint/no-floating-promises
    sync();
  }, [projectHandle, selectedCharacterId]);

  if (!projectHandle) {
    return "No project selected";
  }

  const thisProject = projectStore.getProject(projectHandle);

  if (!thisProject) {
    return "Project not found";
  }

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
                const selected = selectedCharacterId == id;

                return (
                  <Table.Tr
                    className={selected ? "bg-gray-500" : ""}
                    key={id}
                    onClick={() => {
                      setSelectedCharacterId(id);
                    }}
                  >
                    <Table.Th>
                      <div className="flex flex-row gap-2">
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
          {selectedCharcter && (
            <>
              <H2>{selectedCharcter.name}</H2>
              <MarkdownRenderer markdown={selectedCharcter.description} />
            </>
          )}
        </div>
      </div>
    </>
  );
}
