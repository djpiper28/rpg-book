import { Button } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import MarkdownRenderer from "@/components/renderers/markdown";
import { H2 } from "@/components/typography/H2";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type Note, type NoteHandle } from "@/lib/grpcClient/pb/project_note";
import { useTabStore } from "@/stores/tabStore";

interface Props {
  isEditModalOpen: boolean;
  noteHandle: NoteHandle;
  onDelete: () => void;
  onEdit: () => void;
}

export function NoteView(props: Readonly<Props>): ReactNode {
  const [note, setNote] = useState<Note | undefined>();
  const projectHandle = useTabStore((x) => x.selectedTab);

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
    const sync = async () => {
      if (!projectHandle) {
        return;
      }

      try {
        const resp = await getProjectClient().getNote({
          note: props.noteHandle,
          project: projectHandle,
        });

        setNote(resp.response.details);
      } catch (error: unknown) {
        getLogger().error("Cannot get note", {
          error: JSON.stringify(error),
          note: props.noteHandle.id,
          project: JSON.stringify(projectHandle),
        });
      }
    };

    // eslint-disable-next-line @typescript-eslint/no-floating-promises
    sync();
  }, [projectHandle, props.noteHandle, props.isEditModalOpen]);

  if (!note) {
    return <div className="flex-1" />;
  }

  return (
    <div className="flex-1 gap-3 flex flex-col">
      <div className="flex flex-row justify-between items-center">
        <H2>{note.details?.name}</H2>
        <div className="flex flex-row gap-2">
          <Button onClick={props.onEdit}>Edit</Button>
          <Button color="red" onClick={props.onDelete}>
            Delete
          </Button>
        </div>
      </div>
      <MarkdownRenderer markdown={note.details?.markdown ?? ""} className="flex-grow min-h-0" />
    </div>
  );
}
