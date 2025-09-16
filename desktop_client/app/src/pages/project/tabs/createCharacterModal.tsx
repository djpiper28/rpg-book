import { Button } from "@mantine/core";
import { type ReactNode, useState } from "react";
import { getProjectClient } from "@/lib/grpcClient/client";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import { CharacterEdit } from "./characterEdit";
import { BasicCharacterDetails } from "@/lib/grpcClient/pb/project_character";

interface Props {
  closeDialog: () => void;
}

export default function CreateCharacterModal(
  props: Readonly<Props>,
): ReactNode {
  const { setError } = useGlobalErrorStore((x) => x);
  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);
  const [character, setCharacter] = useState<BasicCharacterDetails>({});
  const [icon, setIcon] = useState("");

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <CharacterEdit
        character={character}
        icon={icon}
        setCharacter={setCharacter}
        setIcon={setIcon}
      />
      <Button
        onClick={() => {
          getProjectClient()
            .createCharacter({
              details: character,
              project: projectHandle,
            })
            .then((resp) => {
              projectStore.addCharacter(projectHandle, {
                description: character.description,
                handle: resp.response,
                icon: character.icon,
                name: character.name,
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
