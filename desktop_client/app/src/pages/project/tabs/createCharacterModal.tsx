import { Button } from "@mantine/core";
import { type ReactNode, useState } from "react";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type BasicCharacterDetails } from "@/lib/grpcClient/pb/project_character";
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
  const [iconPath, setIconPath] = useState("");

  const [character, setCharacter] = useState<BasicCharacterDetails>({
    description: "",
    handle: {
      id: "",
    },
    icon: new Uint8Array(),
    iconPath: "",
    name: "",
  });

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <CharacterEdit
        character={character}
        iconPath={iconPath}
        setCharacter={setCharacter}
        setIconPath={setIconPath}
      />
      <Button
        onClick={() => {
          getProjectClient()
            .createCharacter({
              details: character,
              iconPath,
              project: projectHandle,
            })
            .then((resp) => {
              getProjectClient()
                .getCharacter({
                  character: resp.response,
                  project: projectHandle,
                })
                .then((charDetails) => {
                  projectStore.addCharacter(
                    projectHandle,
                    charDetails.response.details,
                  );

                  props.closeDialog();
                })
                .catch((error: unknown) => {
                  getLogger().error("Cannot get character after creating", {
                    error: String(error),
                  });

                  setError({
                    body: String(error),
                    title: "Cannot get character after creating",
                  });
                });
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
