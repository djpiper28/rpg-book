import { Button, Input, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { Pencil } from "lucide-react";
import { type ReactNode, useEffect, useState } from "react";
import { Modal } from "@/components/modal/modal";
import MarkdownRenderer from "@/components/renderers/markdown";
import { H2 } from "@/components/typography/H2";
import { P } from "@/components/typography/P";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type BasicCharacterDetails } from "@/lib/grpcClient/pb/project_character";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import CreateCharacterModal from "./createCharacterModal";
import EditCharacterModal from "./editCharacterModal";
import { uint8ArrayToBase64 } from "@/lib/utils/base64";

export function CharacterTab(): ReactNode {
  const [selectedCharacterId, setSelectedCharacterId] = useState("");

  const [selectedCharcter, setSelectedCharacter] = useState<
    BasicCharacterDetails | undefined
  >();

  const [createOpened, { close: createClose, open: createOpen }] =
    useDisclosure(false);

  const [editOpened, { close: editClose, open: editOpen }] =
    useDisclosure(false);

  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);
  const [iconB64, setIconB64] = useState("");

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
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
        setIconB64(uint8ArrayToBase64(resp.response.icon));
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
  }, [projectHandle, selectedCharacterId, editOpened]);

  if (!projectHandle) {
    return "No project selected";
  }

  const thisProject = projectStore.getProject(projectHandle);

  if (!thisProject) {
    return "Project not found";
  }

  return (
    <>
      <Modal close={createClose} opened={createOpened} title="Create Character">
        {createOpened && <CreateCharacterModal closeDialog={createClose} />}
      </Modal>

      <Modal close={editClose} opened={editOpened} title="Edit Character">
        {editOpened && (
          <EditCharacterModal
            characterHandle={{ id: selectedCharacterId }}
            closeDialog={editClose}
          />
        )}
      </Modal>

      <div className="flex flex-row gap-2 pt-2 justify-between">
        <div className="flex flex-col gap-2 flex-1">
          <div className="flex flex-row gap-2 justify-between">
            <Input className="flex-grow" placeholder="TODO Search bar" />
            <Button
              onClick={() => {
                createOpen();
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
              {thisProject.project.characters.map((character): ReactNode => {
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
              <div className="flex flex-row gap-3 justify-between">
                {iconB64 && (
                  <img
                    alt="User selected"
                    className="max-w-1/2 max-h-1/2"
                    src={`data:image/jpg;base64,${iconB64}`}
                  />
                )}
                <H2>{selectedCharcter.name}</H2>
                <button
                  className="cursor-pointer"
                  onClick={() => {
                    editOpen();
                  }}
                >
                  <P className="flex flex-row gap-1">
                    <Pencil />
                    Edit
                  </P>
                </button>
              </div>
              <MarkdownRenderer markdown={selectedCharcter.description} />
            </>
          )}
        </div>
      </div>
    </>
  );
}
