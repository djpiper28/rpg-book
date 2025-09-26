import { TextInput } from "@mantine/core";
import { type ReactNode } from "react";
import { IconSelector } from "@/components/input/iconSelector";
import { MarkdownEditor } from "@/components/input/markdownEditor";
import { type BasicCharacterDetails } from "@/lib/grpcClient/pb/project_character";

interface Props {
  character: BasicCharacterDetails;
  iconPath: string;
  imageDataB64?: string;
  setCharacter: (details: BasicCharacterDetails) => void;
  setIconPath: (icon: string) => void;
}

export function CharacterEdit(props: Readonly<Props>): ReactNode {
  return (
    <div className="flex flex-col gap-3">
      <div className="flex flex-row gap-3 justify-between items-start">
        <TextInput
          label="Character Name"
          onChange={(x) => {
            const details = structuredClone(props.character);
            details.name = x.target.value;
            props.setCharacter(details);
          }}
          placeholder="John Smith"
          required={true}
          value={props.character.name}
        />
        <IconSelector
          description="Select an icon for your character"
          imageDataB64={props.imageDataB64}
          imagePath={props.iconPath}
          setImagePath={(src) => {
            props.setIconPath(src);
          }}
        />
      </div>
      <MarkdownEditor
        label="Description and notes"
        setValue={(value) => {
          const details = structuredClone(props.character);
          details.description = value;
          props.setCharacter(details);
        }}
        value={props.character.description}
      />
    </div>
  );
}
