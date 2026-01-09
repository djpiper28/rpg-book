import { Button, Table } from "@mantine/core";
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
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";

export function NoteTab(): ReactNode {
  const [createOpened, { close: createClose, open: createOpen }] =
    useDisclosure(false);

  const [editOpened, { close: editClose, open: editOpen }] =
    useDisclosure(false);

  const [deleteOpened, { close: deleteClose, open: deleteOpen }] =
    useDisclosure(false);

  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);
  const thisProject = projectHandle && projectStore.getProject(projectHandle);
  const [queryError, setQueryError] = useState("");
  const [queryText, setQueryText] = useState("");
  const [selectedNoteId, setSelectedNoteId] = useState<string | undefined>();
  const [queryResult, setQueryResult] = useState<Note[]>([]);

  // const selectedNote = thisProject?.project.notes.find(
  //   (n) => n.handle?.id === selectedNoteId,
  // );

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

  return (
    <>
      <Modal close={createClose} opened={createOpened} title="Create Note">
        {createOpened && <CreateNoteModal closeDialog={createClose} />}
      </Modal>

      <Modal close={editClose} opened={editOpened} title="Edit Note">
        {editOpened && (
          <EditNoteModal
            closeDialog={editClose}
            noteHandle={{ id: selectedNoteId }}
          />
        )}
      </Modal>

      <ConfirmModal
        close={deleteClose}
        onConfirm={() => {
          // if (!selectedNote?.handle) {
          //   deleteClose();
          //   return;
          // }
          //
          // const noteHandle = selectedNote.handle;
          //
          // getProjectClient()
          //   .deleteNote({
          //     handle: noteHandle,
          //     project: thisProject.handle,
          //   })
          //   .then(() => {
          //     projectStore.deleteNote(projectHandle, noteHandle);
          //     setSelectedNoteId("");
          //   })
          //   .catch((error: unknown) => {
          //     errorStore.setError({
          //       body: JSON.stringify(error),
          //       title: "Cannot delete note",
          //     });
          //   });
        }}
        opened={deleteOpened}
        title={`Delete Note TODO: name`}
      />

      <div className="flex flex-row gap-2 pt-2 justify-between">
        <div className="flex flex-col gap-2 flex-1">
          <Search<Note>
            elementWrapper={(children: ReactNode[]): ReactNode => {
              return (
                <Table variant="vertical">
                  <Table.Thead>
                    <Table.Tr>
                      <Table.Th>Name</Table.Th>
                      <Table.Th>Related Entities</Table.Th>
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
            render={(note: Note) => {
              const id = note.handle?.id ?? "";
              const selected = selectedNoteId == id;

              return (
                <Table.Tr
                  className={selected ? "bg-gray-500" : ""}
                  key={id}
                  onClick={() => {
                    setSelectedNoteId(id);
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
        </div>

        {selectedNoteId && (
          <NoteView
            isEditModalOpen={editOpened}
            noteId={selectedNoteId}
            onDelete={deleteOpen}
            onEdit={editOpen}
          />
        )}
      </div>
    </>
  );
}
