import { TextInput } from "@mantine/core";
import { type ReactNode } from "react";
import { MarkdownEditor } from "@/components/input/markdownEditor";

interface Props {
  markdown: string;
  name: string;
  onMarkdownChange: (markdown: string) => void;
  onNameChange: (name: string) => void;
}

export function NoteEdit(props: Readonly<Props>): ReactNode {
  return (
    <div className="flex flex-col gap-3">
      <TextInput label="Name" placeholder="Note #23" value={props.name} />
      <MarkdownEditor
        label="Note"
        setValue={props.onMarkdownChange}
        value={props.markdown}
      />
    </div>
  );
}
