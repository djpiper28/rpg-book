import { Button } from "@mantine/core";
import { type ReactNode, useState } from "react";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type BasicCharacterDetails } from "@/lib/grpcClient/pb/project_character";
import { base64ToUint8Array } from "@/lib/utils/base64";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import { CharacterEdit } from "./characterEdit";

interface Props {
  closeDialog: () => void;
}

export default function CreateCharacterModal(
  props: Readonly<Props>,
): ReactNode {
  const { setError } = useGlobalErrorStore((x) => x);
  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);
  const [iconB64, setIconB64] = useState("");

  const [character, setCharacter] = useState<BasicCharacterDetails>({
    description: "",
    handle: {
      id: "",
    },
    icon: new Uint8Array(),
    name: "",
  });

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <CharacterEdit
        character={character}
        iconB64={iconB64}
        setCharacter={setCharacter}
        setIconB64={setIconB64}
      />
      <Button
        onClick={() => {
          getProjectClient()
            .createCharacter({
              details: {
                ...character,
                icon: base64ToUint8Array(iconB64),
              },
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
              getLogger().error("Cannot create character", {
                error: String(error),
              });

              setError({
                body: String(error),
                title: "Cannot create character",
              });
            });
        }}
      >
        Create Character
      </Button>
    </div>
  );
}
