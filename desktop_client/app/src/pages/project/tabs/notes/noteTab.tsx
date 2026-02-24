import { Button, ScrollArea, Table } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { type ReactNode, useEffect, useState } from "react";
import { CreateNoteModal } from "@/components/entityViews/note/CreateNoteModal";
import { EditNoteModal } from "@/components/entityViews/note/EditNoteModal";
import { NoteView } from "@/components/entityViews/note/NoteView";
import { ConfirmModal } from "@/components/modal/confirmModal";
import { Modal } from "@/components/modal/modal";
import { Search } from "@/components/search/search";
import { P } from "@/components/typography/P";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type Note } from "@/lib/grpcClient/pb/project_note";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";

export function NoteTab(): ReactNode {
  const [createOpened, { close: createClose, open: createOpen }] =
    useDisclosure(false);

  const [editOpened, { close: editClose, open: editOpen }] =
    useDisclosure(false);

  const [deleteOpened, { close: deleteClose, open: deleteOpen }] =
    useDisclosure(false);

  const setSelectedNote = useTabStore((x) => x.selected.setSelectedNote);
  const projectHandle = useTabStore((x) => x.selectedTab);
  const tabs = useTabStore((x) => x.tabs);
  const projectStore = useProjectStore((x) => x);
  const errorStore = useGlobalErrorStore((x) => x);
  const thisProject = projectHandle && projectStore.getProject(projectHandle);
  const [queryError, setQueryError] = useState("");
  const [queryText, setQueryText] = useState("");
  const [queryResult, setQueryResult] = useState<Note[]>([]);

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
        setQueryResult(thisProject.project.notes);
        setQueryError("");
        return;
      }

      try {
        const resp = await getProjectClient().searchNote({
          project: projectHandle,
          query: queryText,
        });

        setQueryResult(
          thisProject.project.notes.filter((allNotes) => {
            const inSearchRes = resp.response.details.find(
              (n) => n.id === allNotes.handle?.id,
            );

            return !!inSearchRes;
          }),
        );

        setQueryError("");
      } catch (error: unknown) {
        getLogger().error("Cannot search notes", {
          error: JSON.stringify(error),
          project: JSON.stringify(projectHandle),
          query: queryText,
        });

        setQueryError(JSON.stringify(error));
      }
    };

    // eslint-disable-next-line @typescript-eslint/no-floating-promises
    sync();
  }, [projectHandle, queryText, thisProject, thisProject?.project.notes]);

  if (!projectHandle) {
    return "No project selected";
  }

  if (!thisProject) {
    return "Project not found";
  }

  const selectedNoteHandle = tabs[projectHandle.id].selectedNote;

  const selectedNote = thisProject.project.notes.find(
    (x) => x.handle === selectedNoteHandle,
  );

  return (
    <>
      <Modal close={createClose} opened={createOpened} title="Create Note">
        {createOpened && <CreateNoteModal closeDialog={createClose} />}
      </Modal>

      <Modal close={editClose} opened={editOpened} title="Edit Note">
        {editOpened && selectedNoteHandle && (
          <EditNoteModal
            closeDialog={editClose}
            noteHandle={selectedNoteHandle}
          />
        )}
      </Modal>

      <ConfirmModal
        close={deleteClose}
        onConfirm={() => {
          if (!selectedNoteHandle) {
            deleteClose();
            return;
          }

          getProjectClient()
            .deleteNote({
              handle: selectedNoteHandle,
              project: thisProject.handle,
            })
            .then(() => {
              projectStore.deleteNote(projectHandle, selectedNoteHandle);
              setSelectedNote(projectHandle.id, { id: "" });
              deleteClose();
            })
            .catch((error: unknown) => {
              errorStore.setError({
                body: JSON.stringify(error),
              });
            });
        }}
        opened={deleteOpened}
        title={`Delete Note ${selectedNote?.details?.name ?? ""}`}
      />

      <div className="flex flex-row gap-2 pt-2 justify-between">
        <div className="flex flex-col gap-2 flex-1 h-full">
          <ScrollArea h="100%" type="always">
            <Search<Note>
              elementWrapper={(children: ReactNode): ReactNode => {
                return (
                  <Table className="w-full" variant="vertical">
                    <Table.Thead className="sticky top-0 bg-gray-800 z-10">
                      <Table.Tr>
                        <Table.Th>Name</Table.Th>
                        <Table.Th>Related Entities</Table.Th>
                      </Table.Tr>
                    </Table.Thead>
                    <Table.Tbody>{children}</Table.Tbody>
                  </Table>
                );
              }}
              error={queryError}
              onChange={(txt: string) => {
                setQueryText(txt);
              }}
              placeholder="search here or click help"
              render={(note: Note) => {
                const id = note.handle?.id ?? "";
                const selected = selectedNoteHandle?.id === id;

                return (
                  <Table.Tr
                    className={selected ? "bg-gray-500" : ""}
                    key={id}
                    onClick={() => {
                      setSelectedNote(projectHandle.id, note.handle);
                    }}
                  >
                    <Table.Th>
                      <div className="flex flex-row gap-2">
                        <P className="text-wrap">{note.details?.name ?? ""}</P>
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
                  Add Note
                </Button>
              }
              searchRes={queryResult}
            />
          </ScrollArea>
        </div>

        {selectedNoteHandle && (
          <NoteView
            isEditModalOpen={editOpened}
            noteHandle={selectedNoteHandle}
            onDelete={deleteOpen}
            onEdit={editOpen}
          />
        )}
      </div>
    </>
  );
}
