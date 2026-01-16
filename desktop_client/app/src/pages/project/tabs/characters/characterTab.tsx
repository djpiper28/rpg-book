import { Button, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { type ReactNode, useEffect, useState } from "react";
import { CharacterView } from "@/components/entityViews/character/CharacterView";
import CreateCharacterModal from "@/components/entityViews/character/CreateCharacterModal";
import EditCharacterModal from "@/components/entityViews/character/EditCharacterModal";
import { ConfirmModal } from "@/components/modal/confirmModal";
import { Modal } from "@/components/modal/modal";
import { Search } from "@/components/search/search";
import { P } from "@/components/typography/P";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type CharacterDetails } from "@/lib/grpcClient/pb/project_character";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";

export function CharacterTab(): ReactNode {
  const [selectedCharacterId, setSelectedCharacterId] = useState("");

  const [createOpened, { close: createClose, open: createOpen }] =
    useDisclosure(false);

  const [editOpened, { close: editClose, open: editOpen }] =
    useDisclosure(false);

  const [deleteOpened, { close: deleteClose, open: deleteOpen }] =
    useDisclosure(false);

  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);
  const thisProject = projectHandle && projectStore.getProject(projectHandle);
  const errorStore = useGlobalErrorStore((x) => x);
  const [queryText, setQueryText] = useState("");
  const [queryResult, setQueryResult] = useState<CharacterDetails[]>([]);
  const [queryError, setQueryError] = useState("");

  const selectedCharacter = thisProject?.project.characters.find(
    (c) => c.handle?.id === selectedCharacterId,
  );

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
    const sync = async () => {
      if (!projectHandle) {
        return;
      }

      if (!thisProject) {
        return;
      }

      if (queryText.trim() === "") {
        setQueryResult(thisProject.project.characters);
        setQueryError("");
        return;
      }

      try {
        const resp = await getProjectClient().searchCharacter({
          project: projectHandle,
          query: queryText,
        });

        setQueryResult(
          thisProject.project.characters.filter((allCharacters) => {
            const inSearchRes = resp.response.details.find(
              (c) => c.id === allCharacters.handle?.id,
            );

            return !!inSearchRes;
          }),
        );

        setQueryError("");
      } catch (error: unknown) {
        getLogger().error("Cannot search characters", {
          error: JSON.stringify(error),
          project: JSON.stringify(projectHandle),
          query: queryText,
        });

        setQueryError(JSON.stringify(error));
      }
    };

    // eslint-disable-next-line @typescript-eslint/no-floating-promises
    sync();
  }, [projectHandle, queryText, thisProject, thisProject?.project.characters]);

  if (!projectHandle) {
    return "No project selected";
  }

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

      <ConfirmModal
        close={deleteClose}
        onConfirm={() => {
          if (!selectedCharacter?.handle) {
            deleteClose();
            return;
          }

          const characterHandle = selectedCharacter.handle;

          getProjectClient()
            .deleteCharacter({
              handle: characterHandle,
              project: thisProject.handle,
            })
            .then(() => {
              projectStore.deleteCharacter(projectHandle, characterHandle);
              setSelectedCharacterId("");
            })
            .catch((error: unknown) => {
              errorStore.setError({
                body: JSON.stringify(error),
                title: "Cannot delete character",
              });
            });
        }}
        opened={deleteOpened}
        title={`Delete Character "${selectedCharacter?.name ?? "no name"}"`}
      />

      <div className="flex flex-row gap-2 pt-2 justify-between">
        <div className="flex flex-col gap-2 flex-1">
          <Search<CharacterDetails>
            elementWrapper={(children: ReactNode[]): ReactNode => {
              return (
                <Table variant="vertical">
                  <Table.Thead>
                    <Table.Tr>
                      <Table.Th>Name</Table.Th>
                      <Table.Th>Factions</Table.Th>
                    </Table.Tr>
                  </Table.Thead>
                  {children}
                </Table>
              );
            }}
            error={queryError}
            onChange={(txt: string) => {
              setQueryText(txt);
            }}
            placeholder="search here or click help"
            render={(character: CharacterDetails) => {
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
            }}
            rightElement={
              <Button
                onClick={() => {
                  createOpen();
                }}
              >
                Add Character
              </Button>
            }
            searchRes={queryResult}
          />
        </div>

        {selectedCharacterId && (
          <CharacterView
            characterId={selectedCharacterId}
            isEditModalOpen={editOpened}
            onDelete={deleteOpen}
            onEdit={editOpen}
          />
        )}
      </div>
    </>
  );
}
