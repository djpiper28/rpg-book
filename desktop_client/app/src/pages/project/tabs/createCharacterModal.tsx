import { Button, TextInput } from "@mantine/core";
import { useState } from "react";
import { MarkdownEditor } from "@/components/input/markdownEditor";
import { projectClient } from "@/lib/grpcClient/client";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";

export default function CreateCharacterModal() {
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
          projectClient
            .createCharacter({
              name: characterName,
              project: projectHandle,
            })
            .then((resp) => {
              console.log("Created a player", resp.response);

              projectStore.addCharacter(projectHandle, {
                handle: resp.response,
                name: characterName,
              });

              close();
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
