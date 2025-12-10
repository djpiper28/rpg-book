import { Button } from "@mantine/core";
import { type ReactNode, useState } from "react";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import { type BasicCharacterDetails } from "@/lib/grpcClient/pb/project_character";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";
import { CharacterEdit } from "./CharacterEdit";

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
    name: "",
  });

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <CharacterEdit
        description={character.description}
        iconPath={iconPath}
        name={character.name}
        onDescriptionChange={(description) => {
          setCharacter((prev) => ({ ...prev, description }));
        }}
        onNameChange={(name) => {
          setCharacter((prev) => ({ ...prev, name }));
        }}
        setIconPath={setIconPath}
      />
      <Button
        onClick={() => {
          // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
          const f = async () => {
            try {
              const resp = await getProjectClient().createCharacter({
                details: character,
                iconPath,
                project: projectHandle,
              });

              const charDetails = await getProjectClient().getCharacter({
                character: resp.response,
                project: projectHandle,
              });

              if (!charDetails.response.details) {
                return;
              }

              projectStore.addCharacter(
                projectHandle,
                charDetails.response.details,
              );

              props.closeDialog();
            } catch (error: unknown) {
              getLogger().error("Cannot create character", {
                error: String(error),
              });

              setError({
                body: String(error),
                title: "Cannot create character",
              });
            }
          };

          // eslint-disable-next-line @typescript-eslint/no-floating-promises
          f();
        }}
      >
        Create Character
      </Button>
    </div>
  );
}
