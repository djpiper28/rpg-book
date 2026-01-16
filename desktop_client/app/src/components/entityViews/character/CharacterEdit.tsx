import { TextInput } from "@mantine/core";
import { type ReactNode } from "react";
import { IconSelector } from "@/components/input/iconSelector";
import { MarkdownEditor } from "@/components/input/markdownEditor";

interface Props {
  description: string;
  iconPath: string;
  imageUrl?: string;
  name: string;
  onDescriptionChange: (description: string) => void;
  onNameChange: (name: string) => void;
  setIconPath: (icon: string) => void;
}

export function CharacterEdit(props: Readonly<Props>): ReactNode {
  return (
    <div className="flex flex-col gap-3">
      <div className="flex flex-row gap-3 justify-between items-start">
        <TextInput
          label="Character Name"
          onChange={(x) => {
            props.onNameChange(x.target.value);
          }}
          placeholder="John Smith"
          required={true}
          value={props.name}
        />
        <IconSelector
          description="Select an icon for your character"
          imagePath={props.iconPath}
          imageUrl={props.imageUrl}
          setImagePath={props.setIconPath}
        />
      </div>
      <MarkdownEditor
        label="Description"
        setValue={props.onDescriptionChange}
        value={props.description}
      />
    </div>
  );
}
