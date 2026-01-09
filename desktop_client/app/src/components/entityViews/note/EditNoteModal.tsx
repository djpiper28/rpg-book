import { Button } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import {
  type NoteDetails,
  type NoteHandle,
} from "@/lib/grpcClient/pb/project_note";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import { NoteEdit } from "./NoteEdit";

interface Props {
  closeDialog: () => void;
  noteHandle: NoteHandle;
}

export function EditNoteModal(props: Readonly<Props>): ReactNode {
  const [noteDetails, setNoteDetails] = useState<NoteDetails>({
    handle: {
      id: "",
    },
    markdown: "",
    name: "",
  });

  const projectStore = useProjectStore((x) => x);
  const { setError } = useGlobalErrorStore((x) => x);
  const projectHandle = useTabStore((x) => x.selectedTab);

  useEffect(() => {
    if (!projectHandle) {
      return;
    }

    getProjectClient()
      .getNote({
        note: props.noteHandle,
        project: projectHandle,
      })
      .then((resp) => {
        if (!resp.response.details?.details) {
          return;
        }

        setNoteDetails(resp.response.details.details);
      })
      .catch((error: unknown) => {
        setError({
          body: String(error),
        });
      });
  }, [props.noteHandle, projectHandle, setError]);

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <NoteEdit
        markdown={noteDetails.markdown}
        name={noteDetails.name}
        onMarkdownChange={(markdown) => {
          setNoteDetails((prev) => ({ ...prev, markdown }));
        }}
        onNameChange={(name) => {
          setNoteDetails((prev) => ({ ...prev, name }));
        }}
      />
      <Button
        onClick={() => {
          // const f = async (): Promise<void> => {
          //   try {
          //     await getProjectClient().updateNote({
          //       details: characterDetails,
          //       handle: props.noteHandle,
          //       project: projectHandle,
          //     });
          //
          //     projectStore.addCharacter(projectHandle, characterDetails);
          //
          //     props.closeDialog();
          //   } catch (error: unknown) {
          //     getLogger().error("Cannot edit note", {
          //       error: String(error),
          //     });
          //
          //     setError({
          //       body: String(error),
          //     });
          //   }
          // };
          //
          // // eslint-disable-next-line @typescript-eslint/no-floating-promises
          // f();
        }}
      >
        Save
      </Button>
    </div>
  );
}
