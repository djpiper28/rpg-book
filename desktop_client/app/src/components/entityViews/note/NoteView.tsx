import { Button } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import MarkdownRenderer from "@/components/renderers/markdown";
import { H2 } from "@/components/typography/H2";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type Note } from "@/lib/grpcClient/pb/project_note";
import { useTabStore } from "@/stores/tabStore";

interface Props {
  isEditModalOpen: boolean;
  noteId: string;
  onDelete: () => void;
  onEdit: () => void;
}

export function NoteView(props: Readonly<Props>): ReactNode {
  const [note, setNote] = useState<Note | undefined>();
  const projectHandle = useTabStore((x) => x.selectedTab);

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
    const sync = async () => {
      if (!props.noteId) {
        return;
      }

      if (!projectHandle) {
        return;
      }

      try {
        const resp = await getProjectClient().getNote({
          note: {
            id: props.noteId,
          },
          project: projectHandle,
        });

        setNote(resp.response.details);
      } catch (error: unknown) {
        getLogger().error("Cannot get note", {
          error: JSON.stringify(error),
          note: props.noteId,
          project: JSON.stringify(projectHandle),
        });
      }
    };

    // eslint-disable-next-line @typescript-eslint/no-floating-promises
    sync();
  }, [projectHandle, props.noteId, props.isEditModalOpen]);

  if (!note) {
    return <div className="flex-1" />;
  }

  return (
    <div className="flex-1 gap-3 overflow-y-auto flex flex-col">
      <div className="flex flex-row justify-between items-center">
        <H2>{note.details?.name}</H2>
        <div className="flex flex-row gap-2">
          <Button onClick={props.onEdit}>Edit</Button>
          <Button color="red" onClick={props.onDelete}>
            Delete
          </Button>
        </div>
      </div>
      <MarkdownRenderer markdown={note.details?.markdown ?? ""} />
    </div>
  );
}
