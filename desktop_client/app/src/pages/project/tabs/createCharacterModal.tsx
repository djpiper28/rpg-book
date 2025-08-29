import { Button, TextInput } from "@mantine/core";
import { useState } from "react";
import { MarkdownEditor } from "@/components/input/markdownEditor";
import { getProjectClient } from "@/lib/grpcClient/client";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";

interface Props {
  closeDialog: () => void;
}

export default function CreateCharacterModal(props: Readonly<Props>) {
  const [characterName, setCharacterName] = useState("");
  const [characterDescription, setCharacterDescription] = useState("");
  const { setError } = useGlobalErrorStore((x) => x);
  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <TextInput
        label="Character Name"
        onChange={(x) => {
          setCharacterName(x.target.value);
        }}
        placeholder="John Smith"
        required={true}
        value={characterName}
      />
      <MarkdownEditor
        label="Description and notes"
        setValue={(value) => {
          setCharacterDescription(value);
        }}
        value={characterDescription}
      />
      <Button
        onClick={() => {
          getProjectClient()
            .createCharacter({
              name: characterName,
              project: projectHandle,
            })
            .then((resp) => {
              projectStore.addCharacter(projectHandle, {
                description: "",
                handle: resp.response,
                icon: new Uint8Array(),
                name: characterName,
              });

              props.closeDialog();
            })
            .catch((error: unknown) => {
              setError({
                body: String(error),
                title: "Cannot create character",
              });
            });
        }}
      >
        Create
      </Button>
    </div>
  );
}
