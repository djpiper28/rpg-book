import { Button, TextInput } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import { IconSelector } from "@/components/input/iconSelector";
import { MarkdownEditor } from "@/components/input/markdownEditor";
import { P } from "@/components/typography/P";
import { getProjectClient } from "@/lib/grpcClient/client";
import { type CharacterHandle } from "@/lib/grpcClient/pb/project_character";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useTabStore } from "@/stores/tabStore";

interface Props {
  characterHandle: CharacterHandle;
  closeDialog: () => void;
}

export default function CreateCharacterModal(
  props: Readonly<Props>,
): ReactNode {
  const [characterName, setCharacterName] = useState("");
  const [characterDescription, setCharacterDescription] = useState("");
  const [icon, setIcon] = useState<string>("");
  const { setError } = useGlobalErrorStore((x) => x);
  const projectHandle = useTabStore((x) => x.selectedTab);

  useEffect(() => {
    if (!projectHandle) {
      return;
    }

    getProjectClient()
      .getCharacter({
        character: props.characterHandle,
        project: projectHandle,
      })
      .then((resp) => {
        setCharacterName(resp.response.name);
        setCharacterDescription(resp.response.description);
        setIcon(resp.response.icon.toBase64());
      })
      .catch((error: unknown) => {
        setError({
          body: String(error),
        });
      });
  }, [props.characterHandle, projectHandle, setError]);

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <div className="flex flex-row gap-3 justify-between items-start">
        <TextInput
          label="Character Name"
          onChange={(x) => {
            setCharacterName(x.target.value);
          }}
          placeholder="John Smith"
          required={true}
          value={characterName}
        />
        <IconSelector
          description="Select an icon for your character"
          setSrc={(src) => {
            setIcon(src);
          }}
          src={icon}
        />
      </div>
      <MarkdownEditor
        label="Description and notes"
        setValue={(value) => {
          setCharacterDescription(value);
        }}
        value={characterDescription}
      />
      <Button
        onClick={() => {
          // TODO: this
        }}
      >
        Save
      </Button>
    </div>
  );
}
