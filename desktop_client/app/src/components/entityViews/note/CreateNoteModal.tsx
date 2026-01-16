import { Button } from "@mantine/core";
import { type ReactNode, useState } from "react";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type NoteDetails } from "@/lib/grpcClient/pb/project_note";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import { NoteEdit } from "./NoteEdit";

interface Props {
  closeDialog: () => void;
}

export function CreateNoteModal(props: Readonly<Props>): ReactNode {
  const { setError } = useGlobalErrorStore((x) => x);
  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);

  const [note, setNote] = useState<NoteDetails>({
    markdown: "",
    name: "",
  });

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <NoteEdit
        markdown={note.markdown}
        name={note.name}
        onMarkdownChange={(markdown) => {
          setNote((prev) => ({ ...prev, markdown }));
        }}
        onNameChange={(name) => {
          setNote((prev) => ({ ...prev, name }));
        }}
      />
      <Button
        onClick={() => {
          // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
          const f = async () => {
            try {
              const resp = await getProjectClient().createNote({
                characters: [],
                details: note,
                project: projectHandle,
              });

              const noteDetails = await getProjectClient().getNote({
                note: resp.response,
                project: projectHandle,
              });

              if (!noteDetails.response.details) {
                return;
              }

              projectStore.addNote(projectHandle, noteDetails.response.details);

              props.closeDialog();
            } catch (error: unknown) {
              getLogger().error("Cannot create note", {
                error: String(error),
              });

              setError({
                body: String(error),
                title: "Cannot create note",
              });
            }
          };

          // eslint-disable-next-line @typescript-eslint/no-floating-promises
          f();
        }}
      >
        Create Note
      </Button>
    </div>
  );
}
